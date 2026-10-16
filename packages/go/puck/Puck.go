package puck

import (
	"fmt"
	"strings"

	GT "github.com/verikono/glocalize/types/go"
)

type DataFormat string
type UrlProtocol string
type Url = string

const (
	ProtocolFile UrlProtocol = "file"
	ProtocolHttp UrlProtocol = "http"
)

const (
	FormatGraphQL DataFormat = "graphql"
	FormatPrisma  DataFormat = "prisma"
	FormatUnknown DataFormat = "unknown"
)

type Puck struct {
	url      Url
	protocol UrlProtocol
	format   DataFormat
	IPuck
}

func NewPuck(options ...GT.OptionerFn) *Puck {
	puck := Puck{
		url:      "file:///dev/null",
		format:   FormatUnknown,
		protocol: ProtocolFile,
	}

	for _, option := range options {
		err := option(&puck)
		if err != nil {
			panic(err)
		}
	}
	return &puck
}

func (p *Puck) Path() string {
	switch p.protocol {
	case ProtocolFile:
		return p.url[7:]
	default:
		return ""
	}
}

func (p *Puck) Protocol() UrlProtocol {
	if len(p.url) > 0 {
		parts := strings.Split(p.url, "://")
		if IsUrlProtocol(parts[0]) {
			return UrlProtocol(parts[0])
		}
	}
	panic(fmt.Errorf("invalid url: %s", p.url))
}

func (p *Puck) Format() DataFormat {
	if len(p.url) > 0 {
		parts := strings.Split(p.url, "://")
		basepath := parts[len(parts)-1]
		segs := strings.Split(basepath, ".")
		if len(segs) < 2 {
			p.format = FormatUnknown
			return p.format
		}
		ext := segs[len(segs)-1]
		if IsKnownFormat(ext) {
			p.format = DataFormat(ext)
		} else {
			p.format = FormatUnknown
		}
		return p.format
	}
	panic(fmt.Errorf("invalid url: %s", p.url))
}

func IsUrlProtocol(proto string) bool {
	if proto == "file" || proto == "http" {
		return true
	}
	return false
}

func IsKnownFormat(format string) bool {
	if format == "graphql" || format == "prisma" {
		return true
	}
	return false
}
