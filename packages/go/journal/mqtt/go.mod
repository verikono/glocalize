module github.com/verikono/glocalize/packages/go/journal/mqtt

go 1.27.4


require (
	github.com/verikono/glocalize/packages/go/journal v0.1.0
	github.com/verikono/glocalize/types/go v0.1.0
)

replace (
	github.com/verikono/glocalize/packages/go/journal => ../
	github.com/verikono/glocalize/types/go => ../../../../types/go
)
