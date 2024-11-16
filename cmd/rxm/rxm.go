package main

import (
	"errors"
	"flag"
	"fmt"
	"jonesinator/go-medminder/internal/config"
	"jonesinator/go-medminder/internal/database"
	"log"
	"path/filepath"
	"strconv"
)

func handleLsAll(db *database.Database) error {
	prescriptions, err := database.ReadAllPrescriptions(db)
	if err != nil {
		return err
	}

	for _, value := range prescriptions {
		fmt.Printf(
			"%s - %.2f - %s\n",
			value.Name, value.ExpectedCount(), value.RefillDate().Format("2006-01-02"))
	}

	return nil
}

func handleLsOne(db *database.Database, name string) error {
	rx, err := database.ReadPrescription(db, name)
	if err != nil {
		return err
	}

	fmt.Printf("Name:     %s\n", rx.Name)
	fmt.Printf("Expected: %.2f\n", rx.ExpectedCount())
	fmt.Printf("Refill:   %s\n", rx.RefillDate().Format("2006-01-02"))
	fmt.Printf("Updated:  %.2f on %s\n", rx.Quantity, rx.Updated.Format("2006-01-02"))
	return nil
}

func handleLs(db *database.Database, name string) error {
	if name == "" {
		return handleLsAll(db)
	} else {
		return handleLsOne(db, name)
	}
}

func handleAdd(db *database.Database, name string, quantity float64, rate float64) error {
	err := database.CreatePrescription(db, name, quantity, rate)
	if err != nil {
		return err
	}

	return nil
}

func handleRm(db *database.Database, name string) error {
	err := database.DeletePrescription(db, name)
	if err != nil {
		return err
	}

	return nil
}

func handleUp(db *database.Database, name string, field string, value float64) error {
	switch field {
	case "quantity":
		err := database.UpdatePrescriptionQuantity(db, name, value)
		if err != nil {
			return err
		}
	case "rate":
		err := database.UpdatePrescriptionRate(db, name, value)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown field")
	}

	return nil
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

	var actionErr error
	action := flag.Arg(0)
	switch action {
	case "ls":
		actionErr = handleLs(db, flag.Arg(1))
	case "add":
		quantity, err := strconv.ParseFloat(flag.Arg(2), 64)
		if err != nil {
			actionErr = err
			break
		}
		rate, err := strconv.ParseFloat(flag.Arg(3), 64)
		if err != nil {
			actionErr = err
			break
		}
		actionErr = handleAdd(db, flag.Arg(1), quantity, rate)
	case "rm":
		actionErr = handleRm(db, flag.Arg(1))
	case "up":
		value, err := strconv.ParseFloat(flag.Arg(3), 64)
		if err != nil {
			actionErr = err
			break
		}
		actionErr = handleUp(db, flag.Arg(1), flag.Arg(2), value)
	default:
		actionErr = errors.New("unknown action")
	}

	if actionErr != nil {
		log.Fatal("Error:", actionErr)
	}
}
