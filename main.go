package main

import (
	"fmt"
	"github.com/ericmoritz/go-ddd/entities/books"
	"github.com/ericmoritz/go-ddd/impl/books/json"
	"github.com/spf13/cobra"
)

func main() {
	var cmdRoot = &cobra.Command{Use: "go-ddd"}
	var cmdBook = &cobra.Command{
		Use:   "book",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
	}
	var cmdBookCheckout = &cobra.Command{
		Use:   "checkout [book-key]",
		Short: "checkout a book",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := Checkout(books.Key(args[0]))
			if err != nil {
				fmt.Printf("ERROR: %s", err)
			}

		},
	}
	var cmdBookReturn = &cobra.Command{
		Use:   "return [book-key]",
		Short: "return a book",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := Return(books.Key(args[0]))
			if err != nil {
				fmt.Printf("ERROR: %s", err)
			}

		},
	}

	cmdBook.AddCommand(cmdBookCheckout, cmdBookReturn)
	cmdRoot.AddCommand(cmdBook)
	cmdRoot.Execute()
}

func Checkout(key books.Key) error {
	jsonEntity := json.Init()
	return jsonEntity.Get(key)(
		func(book books.Available) error {
			checkedOut, err := book.Checkout()
			if err != nil {
				return err
			}
			fmt.Printf("Book %v has been checkout out", checkedOut.Book().Key)
			return nil
		},
		func(book books.CheckedOut) error {
			return fmt.Errorf("Sorry, %s has already been checked out", book.Book().Key)
		},
	)
}

func Return(key books.Key) error {
	jsonEntity := json.Init()
	return jsonEntity.Get(key)(
		func(book books.Available) error {
			return fmt.Errorf("Sorry, %s is not checked out", book.Book().Key)
		},
		func(book books.CheckedOut) error {
			returned, err := book.Return()
			if err != nil {
				return err
			}
			fmt.Printf("Book %v has been returned", returned.Book().Key)
			return nil
		},
	)
}
