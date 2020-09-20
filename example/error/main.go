package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {
	i, err := a()
	log.Printf("i=%d err=%v", i, err)
}

//
//
//

func a() (int, error) {
	i, err := b()
	if errors.Is(err, ErrFoo) {
		return 0, fmt.Errorf("tragedy: %w", err)
	}

	var bar BarError
	if errors.As(err, &bar) {
		return 1, fmt.Errorf("comedy: %w", err)
	}

	var baz BazError
	if errors.As(err, &baz) {
		return 2, fmt.Errorf("farce: %w", err)
	}

	return i, nil
}

func b() (int, error) {
	if err := c(); err != nil {
		return 0, fmt.Errorf("error executing c: %w", err)
	}
	return 1, nil
}

func c() error {
	// return ErrFoo
	// return BarError{Reason: "😫"}
	// return BazError{Reason: "☹️"}
	return BazError{Reason: "😟", Inner: ErrFoo}
}

//
//
//

var ErrFoo = errors.New("foo error")

//
//
//

type BarError struct {
	Reason string
}

func (e BarError) Error() string {
	return fmt.Sprintf("bar error: %s", e.Reason)
}

//
//
//

type BazError struct {
	Reason string
	Inner  error
}

func (e BazError) Unwrap() error {
	fmt.Println("fuck")
	return e.Inner
}

func (e BazError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("baz error: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("baz error: %s", e.Reason)
}
