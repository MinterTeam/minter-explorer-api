# Minter Explorer api

The official repository of Minter Explorer API service.

## API Docs
https://app.swaggerhub.com/apis-docs/GrKamil/minter-explorer_api

## Related services:
- [explorer-extender](https://github.com/MinterTeam/minter-explorer-extender)
- [explorer-gate](https://github.com/MinterTeam/explorer-gate)
- [explorer-validators](https://github.com/MinterTeam/minter-explorer-validators) - API for validators meta
- [explorer-tools](https://github.com/MinterTeam/minter-explorer-tools) - common packages
- [explorer-genesis-uploader](https://github.com/MinterTeam/explorer-genesis-uploader)

## BUILD

- `go mod download`

- `go build ./cmd/explorer.go`

## USE

### Requirement

- PostgresSQL

- Centrifugo (WebSocket server) [GitHub](https://github.com/centrifugal/centrifugo)

- Explorer Extender [GitHub](https://github.com/MinterTeam/minter-explorer-extender)

### Setup

- be sure [explorer-extender](https://github.com/MinterTeam/minter-explorer-extender) has been installed and work correctly

- build and move the compiled file to the directory e.g. `/opt/minter/explorer`

- copy .env file in extender's directory and fill with own values

#### Run

./explorer

### Env file

```
EXPLORER_DEBUG=1
APP_BASE_COIN=BIP
CENTRIFUGO_LINK=wss://exporer-rtm.minter.network/connection/websocket
CENTRIFUGO_BLOCK_CHANNEL=blocks
DB_HOST=localhost
DB_PORT=5432
DB_POOL_SIZE=20
DB_NAME=explorer
DB_USER=minter
DB_PASSWORD=password
EXPLORER_PORT=8080
MARKET_HOST=https://api.coingecko.com
```