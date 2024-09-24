# Catalog-go

Cats are good, so its a cat providing api server.

## Deployment

### Build

To build the project, simply run `go build -ldflags "-s -w " .` to build the target binary.

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

### Setting up

Run `./catalog-go migrate` to create/update the database schema.

Simply run `./catalog-go`. If no `config.toml` is found, the program will generate a config template for you, just fill it in!

After configuring everything, just crank up the systemd service, and you're good to serve cats. _You may also need a gateway server like `caddy` to serve in HTTPS._
