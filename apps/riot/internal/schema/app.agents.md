### app identity

An `App` is keyed by `@@id([language, name])`, deliberately **not** by name alone:
two different projects both legitimately contain a package called `utils`, and a
single project may hold a TypeScript `utils` package and a Go `utils` module at
the same time. The project relation is therefore *not* part of the identity — an
app belongs to exactly one project (`project_id`), but its uniqueness is scoped
by language rather than by project.

### `name` derivation

The name is not free-form; it is read from the language's own manifest so that a
row can always be traced back to real files:

| language     | source                                  | example                    |
|--------------|-----------------------------------------|----------------------------|
| `typescript` | `package.json` `name`                   | `@kaltra/app-kaltra-fi`     |
| `go`         | `go.mod` `module`                       | `github.com/verikono/mcpikon` |
| `python`     | `pyproject.toml` `project.name`         | `riot-kit`                 |

### `location`

`location` is a path **relative to the owning project's root** (`Project.location`),
because an app can only ever belong to a single project. Absolute app paths would
duplicate the project root on every row and go stale the moment a checkout moves.
