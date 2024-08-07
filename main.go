package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Model representing your data
type User struct {
	gorm.Model
	Name  string
	Email string
}

var db *gorm.DB

func init() {
	// Replace "db-service" with the name of your Kubernetes service
	password := ""
	dsn := "root:" + password + "@tcp(db-service:3306)/kubernetes?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{})
}

func main() {
	r := gin.Default()

	// Create a new user
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	// Get all users
	r.GET("/users", func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	// Get a single user
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// Update a user
	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&User{}).Where("id = ?", id).Updates(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// Delete a user
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := db.Delete(&user, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}
		c.JSON(http.StatusNoContent, gin.H{})
	})

	r.Run(":8080")
}
