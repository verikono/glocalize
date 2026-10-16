package internal

import (
	"fmt"

	"github.com/verikono/glocalize/packages/go/data/prisma/graphql"
	"github.com/verikono/glocalize/packages/go/puck"
	T "github.com/verikono/glocalize/types/go"
)

type RouteMeta struct {
	command string
	target  string
	fn      T.SchemaParser
}

var routes = []RouteMeta{
	{command: "parse", target: "prisma", fn: graphql.ParseGQL},
}

func GetParser(source string, targetFormat string) error {
	s := puck.From(source)
	fmt.Printf("GetParser: source=%s, targetFormat=%s\n", source, targetFormat)

	// look up routes
	for _, route := range routes {
		if route.command == "parse" && route.target == targetFormat {
			result, err := route.fn(source)
			if err != nil {
				fmt.Errorf("failed parsing %s", s)
			}
			fmt.Printf("\n\n---\n%s\n\n", result)
			return nil
		}
	}
	return fmt.Errorf("unknown source/target combination [%s, %s]", source, targetFormat)
}

func KeepMe(source string, route []string) error {

	s := puck.From(source)
	path := route[1:]

	fmt.Printf("/... %s %s", s, path)
	return nil
}
