package modules

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"

)

var ErrorMqConnectionFail = errors.New("ErrorMqConnectionFail")
var ErrorMqChannelExist  =  errors.New("ErrorMqChannelExist")

type RabbitMq struct {
	conn    *amqp.Connection

	channelMap map[string]*amqp.Channel
}

func NewRabbitMq() *RabbitMq {
	mq :=&RabbitMq{
		conn :nil,
		channelMap:make(map[string]*amqp.Channel),
	}

	return mq
}


func (mq *RabbitMq)Connect(amqpUrl string) error{

	var err error
	mq.conn , err = amqp.Dial(amqpUrl)

	if err == nil {
		log.Printf("mq connect : %s\n",amqpUrl)
		go func() {
			errChan := <-mq.conn.NotifyClose(make(chan *amqp.Error))

			if errChan != nil {
				log.Println(errChan)
			}
		}()
	}
	return err

}

func (mq *RabbitMq)ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args map[string]interface{} ) error {
	c, err := mq.getChannel()

	if err == nil {
		defer func(){err = c.Close()}()

		err = c.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	}

	return err
}

func (mq *RabbitMq)QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args map[string]interface{} ) error  {
	c, err := mq.conn.Channel()

	if err == nil {
		defer func(){err = c.Close()}()

		_ , err = c.QueueDeclare(name , durable , autoDelete , exclusive , noWait , args)
	}
	return err
}

func (mq *RabbitMq)QueueBind( queue string , key string ,  exchange string  , noWait bool , args map[string]interface{} )  error {

	c, err := mq.getChannel()
	if err == nil {
		defer func(){err = c.Close()}()

		err = c.QueueBind(queue , key , exchange , noWait , args)
	}
	return err
}


func (mq *RabbitMq)Consume(id string , queue string) (<-chan amqp.Delivery ,error){

	var err error
	if mq.conn != nil {

		_,ok := mq.channelMap[id]
		if ok {
			return nil , ErrorMqChannelExist
		}

		chnl , err := mq.conn.Channel()
		if err == nil {
			chnl.Qos(1 , 0 , false)

			deliveries , err := chnl.Consume(
				queue,
				"",
				false,
				false,
				false,
				false,
				nil,
			)

			if err == nil {
				mq.channelMap[id] = chnl
				return deliveries , nil
			}
		}

	}
	return nil , err
}

func (mq *RabbitMq)Publish(exchange string , exchangeType string)(chan<- amqp.Publishing , error){
	chnl , err := mq.conn.Channel()
	if err == nil {
		err = chnl.ExchangeDeclare(
			exchange,
			exchangeType,
			true,
			false,
			false,
			false,
			nil,
		)
		if err == nil {
			err = chnl.Confirm( false )
			if err == nil {


				confirms := chnl.NotifyPublish(make(chan amqp.Confirmation, 1))

				pubChan := make(chan amqp.Publishing)
				go func(){

					for {
						select {
						case pub :=<-pubChan:
							err := chnl.Publish(
								exchange,
								"",
								false,
								false,
								pub)
							if err != nil {
								log.Println(err)
							}
							confirmOne(confirms)
						}
					}
				}()

				return pubChan, nil
			}
		}
	}
	return nil , err
}

func (mq *RabbitMq)QueueUnbind(name, key, exchange string) error  {
	c, err := mq.getChannel()

	if err == nil{
		defer func(){err = c.Close()}()

		err = c.QueueUnbind(name , key , exchange , nil)
	}
	return err
}

func (mq *RabbitMq)ExchangeDelete(exchange string ) error {
	c, err := mq.getChannel()

	if err == nil {
		defer func(){err = c.Close()}()

		err = c.ExchangeDelete(exchange , true , false)
	}
	return err
}


func (mq *RabbitMq)QueueDelete( queue string ) error  {

	c, err := mq.getChannel()

	if err == nil {
		defer func(){err = c.Close()}()

		var cnt int
		cnt , err = c.QueueDelete(queue , true , true , false )
		if err == nil {
			if cnt != 0 {
				return errors.New("IsNotEmpty")
			}
		}
	}
	return err
}

func (mq *RabbitMq)CloseConsumeChannel( id string) error  {
	chnl , ok := mq.channelMap[id]
	if ok {
		delete(mq.channelMap , id)
		return chnl.Close()
	}
	return  errors.New("cannot_find_channel")
}


func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func (mq *RabbitMq)Shutdown() error{
	if mq.conn == nil {
		log.Println("mq.conn nil")
	}
	for k, v := range mq.channelMap {
		if err :=v.Cancel("", true); err != nil {

			return fmt.Errorf("Consumer_cancel_failed: %s %v\n", k,err)
		}
		if err :=v.Close(); err != nil {

			return fmt.Errorf("Close(): %s %v\n", k,err)
		}

	}
	if mq.conn != nil {
		if err := mq.conn.Close(); err != nil {
			return fmt.Errorf("AMQP connection close error: %s", err)
		}
	}
	return nil
}

func (mq *RabbitMq)StopConsume(id string) error {

	if mq.conn == nil {
		log.Println("StopConsume mq.conn nil")
	}

	chnl ,ok := mq.channelMap[id]
	if ok {
		if err :=chnl.Cancel("", true); err != nil {
			return fmt.Errorf("Consumer_cancel_failed: %s %v\n", id,err)
		}

		if err :=chnl.Close(); err != nil {

			return fmt.Errorf("Close(): %s %v\n", id,err)
		}
		mq.channelMap[id] = nil
		delete(mq.channelMap,id)
	}
	return nil
}

func (mq *RabbitMq)getChannel() (*amqp.Channel , error){
	if mq.conn == nil {
		return nil, ErrorMqConnectionFail
	}
	return mq.conn.Channel()
}