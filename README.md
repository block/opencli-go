# opencli-go

Go types for the [OpenCLI spec](https://opencli.org).

This project implements a subset of a variation of the OpenCLI specification. The current version is **`0.1-block.1`**, which differs from the [upstream spec](https://opencli.org) in two ways:

1. The Document object extends the Command object — see [spectreconsole/open-cli#70](https://github.com/spectreconsole/open-cli/discussions/70)
2. Command supports the field `defaultCommand` — see [spectreconsole/open-cli#77](https://github.com/spectreconsole/open-cli/discussions/77)

## Installation

```sh
go get github.com/block/opencli-go
```

## Usage

```go
import "github.com/block/opencli-go"
```

### Marshaling

`Document.MarshalJSON` always writes the current version constant:

```go
doc := opencli.Document{
    Info: opencli.CliInfo{Version: "1.0.0"},
    Command: opencli.Command{
        Name: "myapp",
        Commands: []opencli.Command{
            {Name: "serve"},
        },
    },
}
data, _ := json.Marshal(doc)
// {"opencli":"0.1-block.1","info":{"version":"1.0.0"},"name":"myapp",...}
```

### Unmarshaling

`Document.UnmarshalJSON` validates the `"opencli"` field and rejects missing or unrecognized versions:

```go
var doc opencli.Document
err := json.Unmarshal(data, &doc)
// Returns an error if "opencli" is missing or unsupported
```

## Version support

This module can unmarshal any OpenCLI version it supports and always marshals to the latest version (`opencli.Version`).

When a future version of this module introduces breaking type changes, it will use [Go's semantic import versioning](https://research.swtch.com/vgo-import) (e.g. `github.com/block/opencli-go/v2`). A program can import multiple major versions simultaneously:

```go
import (
    opencli   "github.com/block/opencli-go"
    opencliv2 "github.com/block/opencli-go/v2"
)
```

## Adding a new schema version

Types are code-generated from JSON Schema files in `schema/`. To add support for a new version:

1. Add the schema file (e.g. `schema/0.2-block.1.json`).
2. Run `just generate-types` to create `internal/v0_2_block_1/types_gen.go`.
3. Add `internal/v0_2_block_1/unmarshal.go` with any custom unmarshal logic (e.g. Arity defaults).
4. Add a `case "0.2-block.1":` branch to `Document.UnmarshalJSON` in `opencli.go`.
5. If this becomes the new current version, update the `Version` constant and re-point the public type aliases to the new internal package.

## Development

This project uses [Hermit](https://cashapp.github.io/hermit/) for tool management. Activate the environment:

```sh
. ./bin/activate-hermit
```

Generate types from schemas:

```sh
just generate-types
```

Run tests:

```sh
go test ./...
```

## Project Resources
 
| Resource                                   | Description                                                                    |
| ------------------------------------------ | ------------------------------------------------------------------------------ |
| [CODEOWNERS](./CODEOWNERS)                 | Outlines the project lead(s)                                                   |
| [GOVERNANCE.md](./GOVERNANCE.md)           | Project governance                                                             |
| [LICENSE](./LICENSE)                       | Apache License, Version 2.0                                                    |
