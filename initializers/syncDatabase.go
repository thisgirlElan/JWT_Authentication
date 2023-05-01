package initializers

import "github.com/thisgirlElan/jwt_auth/models"

func SyncDatabase() {
	Db.AutoMigrate(&models.User{})
}
