package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"jonesinator/go-medminder/internal/config"
	"jonesinator/go-medminder/internal/database"
)

var global_db *database.Database

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/rx", func(c *gin.Context) {
		prescriptions, err := database.ReadAllPrescriptions(global_db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]string, len(prescriptions))
		for i, v := range prescriptions {
			result[i] = v.Name
		}

		c.JSON(http.StatusOK, result)
	})

	r.GET("/rx/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		prescription, err := database.ReadPrescription(global_db, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":     prescription.Name,
			"count":    prescription.ExpectedCount(),
			"refill":   prescription.RefillDate(),
			"rate":     prescription.Rate,
			"quantity": prescription.Quantity,
			"updated":  prescription.Updated,
		})
	})

	r.DELETE("/rx/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		err := database.DeletePrescription(global_db, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	r.POST("/rx/:name", func(c *gin.Context) {
		var json struct {
			Quantity float64 `json:"quantity" binding:"required"`
			Rate     float64 `json:"rate" binding:"required"`
		}

		name := c.Params.ByName("name")
		err := c.Bind(&json)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = database.CreatePrescription(global_db, name, json.Quantity, json.Rate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.PATCH("/rx/:name", func(c *gin.Context) {
		var json struct {
			Quantity float64 `json:"quantity"`
			Rate     float64 `json:"rate"`
		}

		name := c.Params.ByName("name")
		err := c.Bind(&json)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if json.Quantity != 0 {
			err = database.UpdatePrescriptionQuantity(global_db, name, json.Quantity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if json.Rate != 0 {
			err = database.UpdatePrescriptionRate(global_db, name, json.Rate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	return r
}

func main() {
	databaseFlag := flag.String("db", "", "Path to database.")
	flag.Parse()

	var databaseFilePath = *databaseFlag
	if databaseFilePath == "" {
		configDir, err := config.GetConfigDir("go-medminder")
		if err != nil {
			return
		}
		databaseFilePath = filepath.Join(configDir, "db.sqlite3")
	}

	db, err := database.OpenDatabase(databaseFilePath)
	if err != nil {
		return
	}
	defer database.CloseDatabase(db)
	global_db = db

	r := setupRouter()
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
