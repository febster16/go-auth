package main

import (
	"github.com/febster16/go-auth/config"
	"github.com/febster16/go-auth/database"
	"github.com/febster16/go-auth/internal/api"
	"github.com/febster16/go-auth/migration"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	database.ConnectToDB()
	migration.MigrateDB()
}

func main() {
	r := gin.Default()

	api.SetupRoutes(r)

	r.Run()
}
