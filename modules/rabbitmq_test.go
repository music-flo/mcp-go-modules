package modules

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"strings"
	"testing"
	"time"
)

type TestMessage struct {
	Test string `json:"test"`
}

func TestRabbitMq(t *testing.T) {

	mq := NewRabbitMq()
	now := time.Now().Unix()

	ExchangeName := fmt.Sprintf("E.TEST_%X" ,now)
	QueueName := fmt.Sprintf("Q.TEST_%X" ,now)
	Id := "TEST"

	err := mq.Connect("amqp://admin:admin@localhost:5672//mcp")

	if err != nil {
		t.Error(err)
	}
	err = mq.ExchangeDeclare(ExchangeName , "fanout" , true , false , false ,false, nil)
	if err != nil {
		t.Error(err)
	}

	err = mq.QueueDeclare(QueueName, true , false , false  , false , nil)
	if err != nil {
		t.Error(err)
	}

	err = mq.QueueBind(QueueName , "" , ExchangeName , false , nil)
	if err != nil {
		t.Error(err)
	}

	var pub chan<- amqp.Publishing
	var consume <-chan amqp.Delivery

	pub , err = mq.Publish(ExchangeName, "fanout")

	testMsg := TestMessage{Test:"OK"}
	var body []byte
	body , err = json.Marshal(testMsg)

	if err != nil {
		t.Error(err)
	}

	pubUp := amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            body,
		DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
		Priority:        5,
	}
	pub <- pubUp

	consume , err = mq.Consume(Id ,QueueName)
	msg := <-consume


	var consumeMsg TestMessage
	err = json.Unmarshal(msg.Body , &consumeMsg)
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(consumeMsg.Test , testMsg.Test) {
		t.Errorf("message missmatch")
	}

	err = msg.Ack(false)

	if err != nil {
		t.Error(err)
	}

	err = mq.CloseConsumeChannel(Id)
	if err != nil {
		t.Error(err)
	}

	err = mq.QueueUnbind(QueueName , "" , ExchangeName)
	if err != nil {
		t.Error(err)
	}

	err = mq.QueueDelete(QueueName)
	if err != nil {
		t.Error(err)
	}

	err = mq.ExchangeDelete(ExchangeName)
	if err != nil {
		t.Error(err)
	}

	err = mq.Shutdown()
	if err != nil {
		t.Error(err)
	}

}
