package api

import (
	"github.com/MinterTeam/minter-explorer-api/api/v1"
	"github.com/MinterTeam/minter-explorer-api/api/validators"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-pg/pg"
	"golang.org/x/time/rate"
	"gopkg.in/go-playground/validator.v8"
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
	router.Use(cors.Default())              // CORS
	router.Use(gin.ErrorLogger())           // print all errors
	//router.Use(apiRecovery)                 // returns 500 on any code panics
	router.Use(apiMiddleware(db, explorer)) // init global context

	// create ip map
	ipMap := sync.Map{}
	// rate limit
	router.Use(throttle(ipMap))

	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Resource not found.", c)
	})

	// Create base api prefix
	api := router.Group("/api")
	{
		// apply routes of version 1.0
		apiV1.ApplyRoutes(api)
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
