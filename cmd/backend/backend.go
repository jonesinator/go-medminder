package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"jonesinator/go-medminder/internal/config"
	"jonesinator/go-medminder/internal/database"
)

var GlobalDB *database.Database

func HandleReadPrescriptions(c *gin.Context) {
	prescriptions, err := database.ReadAllPrescriptions(GlobalDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]string, len(prescriptions))
	for i, v := range prescriptions {
		result[i] = v.Name
	}

	c.JSON(http.StatusOK, result)
}

func HandleReadPrescription(c *gin.Context) {
	name := c.Params.ByName("name")
	prescription, err := database.ReadPrescription(GlobalDB, name)
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
}

func HandleDeletePrescription(c *gin.Context) {
	name := c.Params.ByName("name")
	err := database.DeletePrescription(GlobalDB, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func HandleCreatePrescription(c *gin.Context) {
	log.Print(c.Request.Body)

	var json struct {
		Quantity float64 `json:"quantity" binding:"required"`
		Rate     float64 `json:"rate" binding:"required"`
	}

	name := c.Params.ByName("name")
	err := c.Bind(&json)
	if err != nil {
		log.Print("AAA")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var rx *database.Prescription
	rx, err = database.CreatePrescription(GlobalDB, name, json.Quantity, json.Rate)
	if err != nil {
		log.Print("BBB")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"name":     rx.Name,
		"count":    rx.ExpectedCount(),
		"refill":   rx.RefillDate(),
		"rate":     rx.Rate,
		"quantity": rx.Quantity,
		"updated":  rx.Updated,
	})
}

func HandleUpdatePrescription(c *gin.Context) {
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
		err = database.UpdatePrescriptionQuantity(GlobalDB, name, json.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if json.Rate != 0 {
		err = database.UpdatePrescriptionRate(GlobalDB, name, json.Rate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var rx *database.Prescription
	rx, err = database.ReadPrescription(GlobalDB, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"name":     rx.Name,
		"count":    rx.ExpectedCount(),
		"refill":   rx.RefillDate(),
		"rate":     rx.Rate,
		"quantity": rx.Quantity,
		"updated":  rx.Updated,
	})
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/rx", HandleReadPrescriptions)
	router.GET("/rx/:name", HandleReadPrescription)
	router.DELETE("/rx/:name", HandleDeletePrescription)
	router.POST("/rx/:name", HandleCreatePrescription)
	router.PATCH("/rx/:name", HandleUpdatePrescription)

	return router
}

func SetupDatabase() (*database.Database, error) {
	databaseFlag := flag.String("db", "", "Path to database.")
	flag.Parse()

	var databaseFilePath = *databaseFlag
	if databaseFilePath == "" {
		configDir, err := config.GetConfigDir("go-medminder")
		if err != nil {
			return nil, err
		}
		databaseFilePath = filepath.Join(configDir, "db.sqlite3")
	}

	db, err := database.OpenDatabase(databaseFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var err error
	GlobalDB, err = SetupDatabase()
	defer database.CloseDatabase(GlobalDB)
	if err != nil {
		log.Fatal(err)
	}

	router := SetupRouter()
	err = router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
