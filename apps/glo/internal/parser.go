package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/verikono/glocalize/packages/go/data/prisma/graphql"
	"github.com/verikono/glocalize/packages/go/puck"
	T "github.com/verikono/glocalize/types/go"
)

type RouteMeta struct {
	command string
	target  T.DataFormat
	fn      T.SchemaParser
}

var routes = []RouteMeta{
	{command: "parse", target: "prisma", fn: graphql.ParseGQL},
	{command: "parse", target: "graphql", fn: graphql.ParsePrisma},
}

type ParseConfig struct {
	source       FileResourceIdentifier
	target       T.AbsoluteFilePath
	targetFormat T.DataFormat
	sourceFormat T.DataFormat
}

func WithFileOut(path string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if config, ok := arg.(*ParseConfig); ok {
				config.target = path
			}
		}
		return nil
	}
}

func WithTargetFormat(format string) OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(*ParseConfig); ok {
				targetFormat, _, err := T.IsDataFormat(format)
				if err != nil {
					return err
				}
				inst.targetFormat = targetFormat
			}
		}
		return nil
	}
}

func getParserFunction(targetFormat T.DataFormat) (T.SchemaParser, error) {
	for _, route := range routes {
		if route.command == "parse" && route.target == targetFormat {
			return route.fn, nil
		}
	}
	return nil, fmt.Errorf("failed identifying parser function for format %s", targetFormat)
}

func ParseSchema(source FileResourceIdentifier, target T.AbsoluteFilePath, options ...OptionerFn) (*T.Schema, error) {
	return ParseSchemas([]FileResourceIdentifier{source}, target, options...)
}

// ParseSchemas merges every source document into one schema before parsing, so
// a shell-expanded glob (a.prisma b.prisma ...) is treated as a single schema.
func ParseSchemas(sources []FileResourceIdentifier, target T.AbsoluteFilePath, options ...OptionerFn) (*T.Schema, error) {

	if len(sources) == 0 {
		return nil, fmt.Errorf("ParseSchemas requires at least one source")
	}

	config := ParseConfig{
		source: sources[0],
		target: target,
	}

	for _, option := range options {
		if err := option(&config); err != nil {
			return nil, err
		}
	}

	if config.targetFormat == "" {
		targetFormat, _, err := T.IsDataFormat(puck.From(config.target).Extension())
		if err != nil {
			return nil, fmt.Errorf("ParseSchema failed resolving target format: %w", err)
		}
		config.targetFormat = targetFormat
	}

	parser, err := getParserFunction(config.targetFormat)
	if err != nil {
		return nil, err
	}

	var merged strings.Builder
	for _, source := range sources {
		sourceContent, err := puck.From(source).Read()
		if err != nil {
			return nil, err
		}
		merged.WriteString(string(sourceContent))
		merged.WriteString("\n")
	}

	result, err := parser(merged.String())
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(config.target, []byte(*result), 0644); err != nil {
		return nil, fmt.Errorf("ParseSchema failed writing %s: %w", config.target, err)
	}

	return result, nil
}

func KeepMe(source string, route []string) error {
	s := puck.From(source)
	path := route[1:]
	fmt.Printf("/... %s %s", s, path)
	return nil
}
