# Catalog-go

## Getting started
It logs cat. To get started, go with me
```shell
cd catalog-go
go mod tidy
go run . migrate
go run .
```

## Deployment

### Build

To officially build the project, simply run `go build -ldflags "-s -w " .` to build the target binary.

And if you want a systemd configuration example, here it is:

```ini
[Unit]
Description=CatALog Golang backend
After=network.target
Wants=network.target

[Service]
Type=simple
ExecStart=/opt/catalog-go/catalog-go
Environment="GIN_MODE=release"
WorkingDirectory=/opt/catalog-go

[Install]
WantedBy=default.target
```

## Configurations
- Prerequisite: Getting started

The configuration file is located at `config.toml` as a structure of the following

- Server
  - Listen: string address to listen on, corresponding to `VITE_BACKEND_ADDR`
  - AllowOrigins: string array of CORS allowed-origins
- Database
  - Type: enumerate of `sqlite3` and nothing else
  - Path: the string filename to the database
  - AllowedReactions: string array of unicodes for possible reactions
- Store
  - StorePath: mock CDN, corresponding to `VITE_CND_ADDR`