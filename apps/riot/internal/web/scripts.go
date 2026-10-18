package web

import (
	_ "embed"

	T "github.com/verikono/glocalize/types/go"
)

//go:embed scripts/react.development.js
var reactDevJs string

//go:embed scripts/graphql.js
var graphQlJs string

//go:embed scripts/riot.js
var riotJs string

var ScriptFiles = map[T.UrlPath]WebScript{
	"/scripts/react.development.js": {name: "react", route: "*", urlPath: "/scripts/react.development.js", path: "./scripts/react.development.js"},
	"/scripts/graphql.js":           {name: "graphql", route: "*", urlPath: "/scripts/graphql.js", path: "./scripts/graphql.js"},
	"/scripts/riot.js":              {name: "riot", route: "*", urlPath: "/scripts/riot.js", path: "./scripts/riot.js"},
}

type WebScript struct {
	name     string
	kind     WebScriptKind
	url      T.Url
	urlPath  T.UrlPath        // path when using browser
	path     RelativeFilePath // path relative to the script file
	route    T.UrlPath        // route matching pattern "*" is a wildcard for all
	content  Utf8String
	language ScriptLanguage
	RoutedStaticContent
}

func init() {
	LoadScripts()
}

func (s *WebScript) GetKind() WebScriptKind {
	return s.kind
}

func (s *WebScript) GetUrlPath() RelativeFilePath {
	return s.urlPath
}

func (s *WebScript) GetContent() Utf8String {
	return s.content
}

// LoadScripts loads each script's source from its embedded file into content,
// so it can be inlined into the document without a filesystem read at runtime.
func LoadScripts() {
	for key, script := range ScriptFiles {
		if content, ok := embeddedScripts[script.path]; ok {
			script.kind = CachedWebScript
			script.content = content
			ScriptFiles[key] = script
		}
	}
}

var embeddedScripts = map[RelativeFilePath]Utf8String{
	"./scripts/react.development.js": Utf8String(reactDevJs),
	"./scripts/graphql.js":           Utf8String(graphQlJs),
	"./scripts/riot.js":              Utf8String(riotJs),
}

var remoteCommonScripts = []T.Url{
	"https://cdn.jsdelivr.net/npm/graphql@17.0.2/execution/execute.min.js",
}
