package modules


type Reporter interface {
	Report(string) error
}

