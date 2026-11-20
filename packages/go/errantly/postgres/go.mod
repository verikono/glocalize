module github.com/verikono/glocalize/packages/go/errantly/postgres

go 1.27.0

require (
	github.com/verikono/glocalize/packages/go/errantly v0.1.0
	github.com/verikono/glocalize/types/go v0.1.0
)

require github.com/lib/pq v1.12.3

replace (
	github.com/verikono/glocalize/packages/go/errantly => ../
	github.com/verikono/glocalize/types/go => ../../../../types/go
)
