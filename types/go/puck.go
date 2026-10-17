package types

import (
	"fmt"
	"regexp"
)

// -- Interfaces -- //
type (
	IPuck interface {
		Protocol() UrlProtocol
		Path() string
		Format() DataFormat
		Read() (string, error)
		ReadGlob(pattern string) (string, error)
		ResourceType() ResourceType
		Sync() error
		Extension() FileExtension
		Exists() bool
		Size() int // size in bytes
	}
)

// -- Enumerators -- //

// Enum (string): Named structured data format
type DataFormat string

const (
	FormatGraphQL DataFormat = "graphql"
	FormatPrisma  DataFormat = "prisma"
	FormatUnknown DataFormat = "unknown"
)

var DataFormatRx = regexp.MustCompile(`^(?i:graphql|prisma)$`)

func IsDataFormat(format string) (DataFormat, bool, error) {
	if DataFormatRx.MatchString(format) {
		return DataFormat(format), true, nil
	}
	return "", false, fmt.Errorf("DataFormat TypeGuard has not passed for value '%s'", format)
}

// Enum (string): protocol declaration embedded in url eg: "file", "http", "https"
type UrlProtocol string

var UrlProtocolRx = regexp.MustCompile(`^(?i:file|http)$`)

func IsUrlProtocol(proto string) bool {
	return UrlProtocolRx.MatchString(proto)
}

const (
	ProtocolFile UrlProtocol = "file"
	ProtocolHttp UrlProtocol = "http"
)

// Enum (string): Named structured data format
type ResourceType string

const (
	DirectoryResource ResourceType = "directory"
	FileResource      ResourceType = "file"
	GlobFileResource  ResourceType = "file"
	UnknownResource   ResourceType = "unknown"
)

// -- TypeAliases -- //
type Url = string
