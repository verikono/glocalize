# clone the mcp service from /home/bren/mcpikon
# provide a graphql server on /home/bren/glocalize/apps/riot/internal/gql (RiotGraph)

```.env
RIOT_URL=http://localhost:6661/ (default port is 6660) # serves the riot app at ()
RIOT_DATAURL=sqlite:///home/bren/.config/riot/data.sqlite (triple slash is correct - for a relative URL it would be sqlite://../../some/dir/db.sqlite)
RIOT_DIRECT_DATAURL=/home/bren/.config/riot/data.sqlite
```

note the datastore/database type is determined by RIOT_DATAURL - which is currently sqlite, a postgres://user:pass@host:5432/db shoudl be also supported (we support sqlite and postgres)

## Arch
- urfav CLI app
- commands: 
  - `riot serve` (aliases: `riot up`)
  instances a RiotServer puts it in the background
  - `riot status` a status on the background service:
    - mcp url
    - gql url
    - database url in use
    - uptime
  - `riot down` kill all processes of riot
### graphql
- graphql server: gqlgen
- client+playground: https://altairgraphql.dev/

### mcp server
- replicate : /home/bren/mcpikon/internal/mcp, yes you can use the client directly, or you can use the graphql service for transactions.
