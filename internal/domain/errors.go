package domain

import (
	"errors"
	"fmt"
)

// Ошибки
const (
	CounterNoAdded   = "counter have not been added"
	CounterNotFound  = "no counter found"
	URLSNoAdded      = "urls have not been added"
	URLSNotFound     = "no urls found"
	URLAlreadyExists = "the URL already exists"
	EmptyData        = "the data is empty"
	EmptyParams      = "the params is empty"
)

var (
	ErrCounterNoAdded   = errors.New(CounterNoAdded)
	ErrCounterNotFound  = errors.New(CounterNotFound)
	ErrURLSNoAdded      = errors.New(URLSNoAdded)
	ErrURLSNotFound     = errors.New(URLSNotFound)
	ErrURLAlreadyExists = errors.New(URLAlreadyExists)
	ErrEmptyData        = errors.New(EmptyData)
	ErrEmptyParams      = errors.New(EmptyParams)
)

func ErrCounterNoAddedMsg(code, msg string) error {
	return fmt.Errorf("%s: code=%s, message=%s", ErrCounterNoAdded.Error(), code, msg)
}
