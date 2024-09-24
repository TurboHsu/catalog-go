# Catalog Go

## Getting started
It logs cat. To get started, go with me
```shell
cd catalog-go
go mod tidy
go run . migrate
go run .
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