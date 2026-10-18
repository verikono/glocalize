package web

type SheetName = string
type RelativeFilePath = string
type Utf8String = string
type CssContent = Utf8String

type GeneralizedKind string

type WebScriptKind GeneralizedKind

const (
	CachedWebScript       WebScriptKind = "cached"
	RemoteCommonWebScript WebScriptKind = "remote-common"
	RemoteModuleWebScript WebScriptKind = "remote-module"
)

type RoutedStaticContent interface {
	GetKind() GeneralizedKind
	GetUrlPath() RelativeFilePath
	GetContent() Utf8String
}

type ScriptLanguage = string

const (
	CommonJavascript  ScriptLanguage = "cjs"
	EcmascriptModule  ScriptLanguage = "esm"
	TypeescriptModule ScriptLanguage = "typescript"
)
