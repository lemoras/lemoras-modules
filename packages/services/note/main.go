package main

import (
	"note"

	u "github.com/lemoras/goutils/api"
)

func Main(in note.Request) (*u.Response, error) {
	return note.Invoke(in)
}
