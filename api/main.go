package api

import (
	"github.com/MinterTeam/minter-explorer-api/api/v1"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
)

// Run API
func Run(db *pg.DB) {
	router := SetupRouter(db)
	err := router.Run()
	helpers.CheckErr(err)
}

// Setup router
func SetupRouter(db *pg.DB) *gin.Engine {
	router := gin.Default()
	router.Use(gin.ErrorLogger())                  // print all errors
	router.Use(gin.Recovery())                     // returns 500 on any code panics
	router.Use(apiMiddleware(db))                  // init global context

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

	// Create Swagger UI
	router.Static("/help", "./help/dist")

	return router
}

// Add necessary services to global context
func apiMiddleware(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}