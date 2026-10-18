package web

import (
	_ "embed"

	T "github.com/verikono/glocalize/types/go"
)

//go:embed style/main.css
var styleMainCss string

// StyleSheets is keyed by urlPath, matching how ScriptFiles is keyed, so a
// request path can be looked up directly.
var StyleSheets = map[T.UrlPath]StyleSheet{
	"/styles/main.css": {name: "main", route: "*", urlPath: "/styles/main.css", path: "./style/main.css"},
}

type StyleSheet struct {
	name    SheetName        // eg "main"
	urlPath T.UrlPath        // eg "/chart/main.style.css"
	route   T.UrlPath        // eg "/chart/" -- when set to "*", loads into every route
	path    RelativeFilePath // eg "./style/main.css"
	content CssContent
	RoutedStaticContent
}

func (s *StyleSheet) GetUrlPath() RelativeFilePath {
	return s.urlPath
}

func (s *StyleSheet) GetContent() Utf8String {
	return s.content
}

func init() {
	LoadSheets()
}

// LoadSheets loads each sheet's css from its embedded source into content, so
// it can be served without a filesystem read at runtime.
func LoadSheets() {
	for key, sheet := range StyleSheets {
		if content, ok := embeddedStyles[sheet.path]; ok {
			sheet.content = content
			StyleSheets[key] = sheet
		}
	}
}

var embeddedStyles = map[RelativeFilePath]CssContent{
	"./style/main.css": CssContent(styleMainCss),
}
