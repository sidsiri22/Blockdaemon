package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

var db *gorm.DB

func main() {
	
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=transactions port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

		db.AutoMigrate(&Transaction{})

	r := gin.Default()
	
	r.POST("/api/transaction/", func(c *gin.Context) {
		var tx Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		tx.ID = uuid.New() 
		tx.Timestamp = time.Now()
		
		if err := db.Create(&tx).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
			return
		}

				c.JSON(http.StatusOK, gin.H{
			"transaction_id": tx.ID,
			"amount":         tx.Amount,
			"timestamp":      tx.Timestamp,
		})
	})
	
	r.Run(":8080")
}
