package main

import (
	"github.com/febster16/go-auth/config"
	"github.com/febster16/go-auth/database"
	_ "github.com/febster16/go-auth/docs"
	"github.com/febster16/go-auth/internal/api"
	"github.com/febster16/go-auth/migration"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	database.ConnectToDB()
	migration.MigrateDB()
}

// @title			GO API Course
// @version		1.0
// @description	GO API using Gin framework.
// @host			localhost:8000
// @BasePath		/
func main() {
	r := gin.Default()

	api.SetupRoutes(r)

	r.Run()
}
