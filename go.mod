module github.com/MinterTeam/minter-explorer-api/v2

go 1.13

require (
	github.com/MinterTeam/minter-explorer-extender/v2 v2.8.1-0.20210122125525-c4124b1bc91c
	github.com/MinterTeam/minter-go-node v1.0.5
	github.com/MinterTeam/minter-go-sdk/v2 v2.1.0-rc2.0.20210201213932-3e24f1f8d283
	github.com/MinterTeam/node-grpc-gateway v1.2.2-0.20210201213545-40b0690753c9
	github.com/centrifugal/centrifuge-go v0.3.0
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-pg/pg/v9 v9.2.0
	github.com/go-pg/urlstruct v0.2.8
	github.com/joho/godotenv v1.3.0
	github.com/prometheus/client_golang v1.4.1
	github.com/sirupsen/logrus v1.4.2
	github.com/zsais/go-gin-prometheus v0.1.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/genproto v0.0.0-20210202153253-cf70463f6119 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
