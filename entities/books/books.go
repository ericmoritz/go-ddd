package books

type Entity interface {
	Get(name Key) BookUnion
}

type BookUnion func(func(Available) error, func(CheckedOut) error) error

type Available interface {
	Book() Book
	Checkout() (CheckedOut, error)
}

type CheckedOut interface {
	Book() Book
	Return() (Available, error)
}

type Key string
type Book struct {
	Key Key
}
