package main

import (
	"drive"

	u "github.com/lemoras/goutils/api"
)

func Main(in drive.Request) (*u.Response, error) {
	return drive.Invoke(in)
}
