name=riot
version=0.1.0

# run phi to regenerate the orm client
generate:
	phi generate

generate-graph:
	gloc parse /home/bren/glocalize/apps/riot/internal/schema/*.prisma --out /home/bren/glocalize/apps/riot/internal/gql/riot.graphql

compile:
	go build -o ./dist/riot

build:
	riot down
	make compile
	make install
	riot up

install:
	cp ./dist/riot ~/.local/bin/riot
