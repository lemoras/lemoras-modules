package main

import (
	"initialize"

	u "github.com/lemoras/goutils/api"
)

func Main(in initialize.Request) (*u.Response, error) {
	return initialize.Invoke(in)
}
