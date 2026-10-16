package puck

import (
	GT "github.com/verikono/glocalize/types/go"
)

type Puck struct {
	url string
}

func NewPuck(options ...GT.OptionerFn) *Puck {
	puck := Puck{
		url: "/dev/null",
	}
	for _, option := range options {
		err := option(puck)
		if err != nil {
			panic(err)
		}
	}
	return &puck
}
