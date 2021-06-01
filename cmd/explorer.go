package main

import (
	"context"
	"github.com/MinterTeam/minter-explorer-api/v2/api"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/ws"
	"github.com/MinterTeam/minter-explorer-api/v2/database"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/metrics"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"os"
)

func main() {
	// Init Logger
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}

	// init environment
	env := core.NewEnvironment()

	startTracing(context.Background())

	// connect to database
	db := database.Connect(env)
	defer database.Close(db)

	// create explorer
	explorer := core.NewExplorer(db, env)

	// run market price update
	go explorer.MarketService.Run()
	go explorer.PoolService.RunPoolUpdater()

	// create ws extender
	extender := ws.NewExtenderWsClient(explorer)
	defer extender.Close()

	// create ws channel handler
	blocksChannelHandler := ws.NewBlocksChannelHandler()
	blocksChannelHandler.AddSubscriber(explorer.Cache)
	blocksChannelHandler.AddSubscriber(metrics.NewLastBlockMetric())
	blocksChannelHandler.AddSubscriber(explorer.PoolService)
	blocksChannelHandler.AddSubscriber(&explorer.CoinRepository)

	// subscribe to channel and add cache handler
	sub := extender.CreateSubscription(explorer.Environment.WsBlocksChannel)
	sub.OnPublish(blocksChannelHandler)
	extender.Subscribe(sub)

	// run api
	api.Run(db, explorer)
}

func startTracing(ctx context.Context) {
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithEndpoint("otlp.uptrace.dev:4317"),
		otlpgrpc.WithHeaders(map[string]string{"uptrace-dsn": os.Getenv("UPTRACE_DSN")}),
		otlpgrpc.WithCompressor("gzip"),
		otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)

	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithMaxQueueSize(1000), sdktrace.WithMaxExportBatchSize(1000))

	tracerProvider := sdktrace.NewTracerProvider()
	tracerProvider.RegisterSpanProcessor(bsp)

	// Install our tracer provider and we are done.
	otel.SetTracerProvider(tracerProvider)
}