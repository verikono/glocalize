package web

type SheetName = string
type RelativeFilePath = string
type Utf8String = string
type CssContent = Utf8String

type RoutedStaticContent interface {
	GetUrlPath() RelativeFilePath
	GetContent() Utf8String
}

type ScriptLanguage = string

const (
	CommonJavascript  ScriptLanguage = "cjs"
	EcmascriptModule  ScriptLanguage = "esm"
	TypeescriptModule ScriptLanguage = "typescript"
)
