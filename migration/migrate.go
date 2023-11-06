package migration

import (
	"github.com/febster16/go-auth/database"
	"github.com/febster16/go-auth/internal/models"
)

func MigrateDB() {
	database.DB.AutoMigrate(&models.User{})
}
