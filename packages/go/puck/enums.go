package puck

import (
	GT "github.com/verikono/glocalize/types/go"
)

// - Enumerators - //
type ResourceType = GT.ResourceType

const (
	DirectoryResource = GT.DirectoryResource
	FileResource      = GT.FileResource
	GlobFileResource  = GT.GlobFileResource
	UnknownResource   = GT.UnknownResource
)

type UrlProtocol = GT.UrlProtocol

const (
	ProtocolFile UrlProtocol = "file"
	ProtocolHttp UrlProtocol = "http"
)

type (
	DataFormat = GT.DataFormat
)

const (
	FormatGraphQL = GT.FormatGraphQL
	FormatPrisma  = GT.FormatPrisma
	FormatUnknown = GT.FormatUnknown
)
