package types

import (
	"strings"
)

// resolves one or multiple files in a filesystem (an absolute path, relative path and glob are examples of "file resource identifiers")
type FileResourceIdentifier = string

type AbsoluteFsPath = string

type AbsoluteFilePath = AbsoluteFsPath

type AbsoluteDirectoryPath = AbsoluteFsPath

type FileExtension = string

type GlobPattern = string

func IsGlobPattern(subject string) bool {
	return strings.Contains(subject, "*")
}
