package types

// A software package version label
type SoftwareVersion = string

type SoftwareLanguageName string

type PackageManagerName string

const (
	PythonLanguage     SoftwareLanguageName = "python"
	TypescriptLanguage SoftwareLanguageName = "typescript"
	ReactLanguage      SoftwareLanguageName = "react"
	LynxLanguage       SoftwareLanguageName = "lynx"
	GoLanguage         SoftwareLanguageName = "go"
	RustLanguage       SoftwareLanguageName = "rust"
	KotlinLanguage     SoftwareLanguageName = "kotlin"
	SwiftLanguage      SoftwareLanguageName = "swift"
)

const (
	Uv   PackageManagerName = "uv"
	Yarn PackageManagerName = "yarn"
)
