module github.com/MinterTeam/minter-explorer-api/v2

go 1.15

require (
	github.com/MinterTeam/explorer-sdk v0.1.1-0.20210521135222-1500182959d7
	github.com/MinterTeam/minter-explorer-extender/v2 v2.12.12-0.20210513095622-810621654289
	github.com/MinterTeam/minter-go-node v1.0.5
	github.com/MinterTeam/minter-go-sdk/v2 v2.2.0-alpha1.0.20210327021641-068a106a883b
	github.com/MinterTeam/node-grpc-gateway v1.2.2-0.20210327013510-012fe2754dfc
	github.com/centrifugal/centrifuge-go v0.3.0
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.7.1
	github.com/go-pg/pg/extra/pgotel/v10 v10.9.3
	github.com/go-pg/pg/v10 v10.9.3
	github.com/go-pg/urlstruct v0.2.8
	github.com/go-playground/validator/v10 v10.4.1
	github.com/joho/godotenv v1.3.0
	github.com/prometheus/client_golang v1.4.1
	github.com/sirupsen/logrus v1.8.1
	github.com/starwander/goraph v0.0.0-20200325033650-cb8f0beb44cc
	github.com/uptrace/uptrace-go v0.20.0
	github.com/uptrace/uptrace-go/extra/otellogrus v0.20.0
	github.com/zsais/go-gin-prometheus v0.1.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)
