module github.com/verikono/glocalize/apps/riot

go 1.27.0

require (
	github.com/coder/websocket v1.8.15
	github.com/fsnotify/fsnotify v1.10.1
	github.com/graphql-go/graphql v0.8.1
	github.com/graphql-go/handler v0.2.4
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.3
	github.com/mark3labs/mcp-go v0.57.0
	github.com/pressly/goose/v3 v3.29.0
	github.com/urfave/cli/v3 v3.11.0
	github.com/verikono/glocalize/packages/go/data/prisma v0.0.0-00010101000000-000000000000
	github.com/verikono/glocalize/packages/go/puck v0.1.0
	github.com/verikono/glocalize/types/go v0.1.0
	modernc.org/sqlite v1.60.1
)

require (
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/google/jsonschema-go v0.4.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	github.com/sethvargo/go-retry v0.5.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	modernc.org/libc v1.78.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)

replace github.com/verikono/glocalize/packages/go/data/prisma => ../../packages/go/data/prisma

replace github.com/verikono/glocalize/types/go => ../../types/go

replace github.com/verikono/glocalize/packages/go/puck => ../../packages/go/puck
