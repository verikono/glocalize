package puck

import (
	GT "github.com/verikono/glocalize/types/go"
)

// - TypeAliases - //
type (
	OptionerFn            = GT.OptionerFn
	AbsoluteFsPath        = string
	AbsoluteFilePath      = AbsoluteFsPath
	AbsoluteDirectoryPath = AbsoluteFsPath
	GlobPattern           = string
	UTF8FileContent       = string
	Url                   = GT.Url
)

// - Interfaces - //
type (
	IPuck = GT.IPuck
)
