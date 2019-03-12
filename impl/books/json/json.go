package json

import (
	"encoding/json"
	"fmt"
	"github.com/ericmoritz/go-ddd/entities/books"
	"io/ioutil"
	"os"
)

func Init() books.Entity {
	return &entity{}
}

type entity struct{}

func (s *entity) Get(key books.Key) books.BookUnion {
	return func(
		onAvailable func(books.Available) error,
		onCheckedout func(books.CheckedOut) error,
	) error {
		book, err := load(key)
		if err != nil {
			return err
		}
		switch book.State {
		case available:
			return onAvailable(&availableImpl{key: key, book: books.Book{Key: key}})
		case checkedout:
			return onCheckedout(&checkedOutImpl{key: key, book: books.Book{Key: key}})
		default:
			return fmt.Errorf("Book is in unknown state: %v", book.State)
		}
	}
}

type state string

const (
	available  = "available"
	checkedout = "checkedout"
)

type book struct {
	State state
}

func load(key books.Key) (book, error) {
	b := book{State: available}

	fn := filename(key)
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return b, nil
	} else if err != nil {
		return b, err
	}
	file, err := os.Open(fn)
	defer file.Close()
	if err != nil {
		return b, err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&b)
	if err != nil {
		return b, err
	}
	return b, nil
}

func save(key books.Key, b book) error {
	output, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("error writing JSON: %#v", err)
	}
	err = ioutil.WriteFile(filename(key), output, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON: %#v", err)
	}
	return nil
}

func filename(key books.Key) string {
	return fmt.Sprintf("./%s.json", string(key))
}

type availableImpl struct {
	key  books.Key
	book books.Book
}

func (b *availableImpl) Book() books.Book {
	return b.book
}
func (b *availableImpl) Checkout() (books.CheckedOut, error) {
	err := save(b.key, book{State: checkedout})
	if err != nil {
		return nil, err
	}
	return &checkedOutImpl{key: b.key, book: b.book}, nil
}

type checkedOutImpl struct {
	key  books.Key
	book books.Book
}

func (b *checkedOutImpl) Book() books.Book {
	return b.book
}
func (b *checkedOutImpl) Return() (books.Available, error) {
	err := save(b.key, book{State: available})
	if err != nil {
		return nil, err
	}
	return &availableImpl{key: b.key, book: b.book}, nil
}
