package api

import (
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2"
	"github.com/MinterTeam/minter-explorer-api/v2/api/validators"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"github.com/zsais/go-gin-prometheus"
	"golang.org/x/time/rate"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"sync"
)

// Run API
func Run(db *pg.DB, explorer *core.Explorer) {
	router := SetupRouter(db, explorer)
	appAddress := ":" + explorer.Environment.ServerPort
	err := router.Run(appAddress)
	helpers.CheckErr(err)
}

// Setup router
func SetupRouter(db *pg.DB, explorer *core.Explorer) *gin.Engine {
	// Set release mode
	if !explorer.Environment.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// metrics
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	router.Use(cors.Default())              // CORS
	router.Use(gin.ErrorLogger())           // print all errors
	router.Use(apiRecovery)                 // returns 500 on any code panics
	router.Use(apiMiddleware(db, explorer)) // init global context

	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Resource not found.", c)
	})

	// Create base api prefix
	api := router.Group("/api")
	{
		// apply routes of version 2.0
		apiV2.ApplyRoutes(api)
	}

	// Register validator for api requests
	registerApiValidators()

	return router
}

// Add necessary services to global context
func apiMiddleware(db *pg.DB, explorer *core.Explorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Set("explorer", explorer)
		c.Next()
	}
}

// Register request validators
func registerApiValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("minterAddress", validators.MinterAddress)
		helpers.CheckErr(err)

		err = v.RegisterValidation("minterTxHash", validators.MinterTxHash)
		helpers.CheckErr(err)

		err = v.RegisterValidation("minterPubKey", validators.MinterPublicKey)
		helpers.CheckErr(err)

		err = v.RegisterValidation("timestamp", validators.Timestamp)
		helpers.CheckErr(err)
	}
}

// Send 500 status and JSON response
func apiRecovery(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			logrus.WithField("err", rec).Error("API error")
			errors.SetErrorResponse(http.StatusInternalServerError, -1, "Internal server error", c)
		}
	}(c)

	c.Next()
}

func throttle(ipMap sync.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter, ok := ipMap.Load(c.ClientIP())
		if !ok {
			limiter = rate.NewLimiter(5, 5)
			ipMap.Store(c.ClientIP(), limiter)
		}

		if !limiter.(*rate.Limiter).Allow() {
			errors.SetErrorResponse(http.StatusTooManyRequests, -1, "Too many requests", c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
