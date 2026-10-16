module github.com/verikono/glocalize/apps/glo

go 1.27.0

require (
	github.com/urfave/cli/v3 v3.11.0
	github.com/verikono/glocalize/packages/go/data/prisma v0.0.0-00010101000000-000000000000
	github.com/verikono/glocalize/packages/go/puck v0.1.0
	github.com/verikono/glocalize/types/go v0.1.0
)

replace github.com/verikono/glocalize/packages/go/data/prisma => ../../packages/go/data/prisma

replace github.com/verikono/glocalize/types/go => ../../types/go

replace github.com/verikono/glocalize/packages/go/puck => ../../packages/go/puck
