module github.com/verikono/glocalize/packages/go/puck

go 1.27.0

require (
	github.com/fsnotify/fsnotify v1.10.1
	github.com/verikono/glocalize/types/go v0.1.0
)

require golang.org/x/sys v0.13.0 // indirect

replace github.com/verikono/glocalize/types/go => ../../../types/go
