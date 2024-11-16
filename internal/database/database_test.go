package database

import (
	"os"
	"testing"
)

func TestDatabaseFile(t *testing.T) {
	// Create an initially-empty database.
	tempFile, err := os.CreateTemp("", "db.")
	if err != nil {
		t.Error(err)
	}
	db, err := OpenDatabase(tempFile.Name())
	if err != nil {
		t.Error(err)
	}

	// Reading the list of prescriptions should yield nothing.
	prescriptions, err := ReadAllPrescriptions(db)
	if err != nil {
		t.Error(err)
	}
	if len(prescriptions) != 0 {
		t.Error("database is not empty")
	}

	// Reading a particular prescription should yield nothing.
	_, err = ReadPrescription(db, "foo")
	if err.Error() != "prescription not found" {
		t.Error(err)
	}

	// Creating a prescription should succeed.
	err = CreatePrescription(db, "foo", 123.45, 12.34)
	if err != nil {
		t.Error(err)
	}

	// Reading it back by reading all prescriptions should succeed.
	prescriptions, err = ReadAllPrescriptions(db)
	if err != nil {
		t.Error(err)
	}
	if len(prescriptions) != 1 {
		t.Error("database is empty")
	}
	if prescriptions[0].Name != "foo" {
		t.Error("incorrect name", prescriptions[0].Name)
	}
	if prescriptions[0].Quantity != 123.45 {
		t.Error("incorrect quantity", prescriptions[0].Quantity)
	}
	if prescriptions[0].Rate != 12.34 {
		t.Error("incorrect rate", prescriptions[0].Rate)
	}

	// Creating a second prescription should work.
	err = CreatePrescription(db, "bar", 543.21, 43.21)
	if err != nil {
		t.Error(err)
	}

	// Reading it back by reading all prescriptions should succeed.
	prescriptions, err = ReadAllPrescriptions(db)
	if err != nil {
		t.Error(err)
	}
	if len(prescriptions) != 2 {
		t.Error("database has incorrect number of records")
	}

	// Reading it back by reading that particluar prescription should succeed.
	prescription, err := ReadPrescription(db, "bar")
	if err != nil {
		t.Error(err)
	}
	if prescription.Name != "bar" {
		t.Error("incorrect name", prescription.Name)
	}
	if prescription.Quantity != 543.21 {
		t.Error("incorrect quantity", prescription.Quantity)
	}
	if prescription.Rate != 43.21 {
		t.Error("incorrect rate", prescription.Rate)
	}

	// Trying to create a prescription that already exists should fail.
	err = CreatePrescription(db, "bar", 1, 1)
	if err.Error() != "prescription already exists" {
		t.Error(err)
	}

	// Deleting it should succeed.
	err = DeletePrescription(db, "bar")
	if err != nil {
		t.Error(err)
	}

	// Trying to delete it again should fail.
	err = DeletePrescription(db, "bar")
	if err.Error() != "prescription not found" {
		t.Error(err)
	}

	// Updating the quantity should succeed, and we should be able to read it
	// back.
	err = UpdatePrescriptionQuantity(db, "foo", 444.44)
	if err != nil {
		t.Error(err)
	}
	prescription, err = ReadPrescription(db, "foo")
	if err != nil {
		t.Error(err)
	}
	if prescription.Name != "foo" {
		t.Error("incorrect name", prescription.Name)
	}
	if prescription.Quantity != 444.44 {
		t.Error("incorrect quantity", prescription.Quantity)
	}
	if prescription.Rate != 12.34 {
		t.Error("incorrect rate", prescription.Rate)
	}

	// Updating the rate should succeed, and we should be able to read it
	// back.
	err = UpdatePrescriptionRate(db, "foo", 33.33)
	if err != nil {
		t.Error(err)
	}
	prescription, err = ReadPrescription(db, "foo")
	if err != nil {
		t.Error(err)
	}
	if prescription.Name != "foo" {
		t.Error("incorrect name", prescription.Name)
	}
	if prescription.Quantity != 444.44 {
		t.Error("incorrect quantity", prescription.Quantity)
	}
	if prescription.Rate != 33.33 {
		t.Error("incorrect rate", prescription.Rate)
	}

	// Delete the last record. It should succeed.
	err = DeletePrescription(db, "foo")
	if err != nil {
		t.Error(err)
	}

	// Updating quantities on non-existent records should fail.
	err = UpdatePrescriptionQuantity(db, "foo", 555.55)
	if err.Error() != "prescription not found" {
		t.Error(err)
	}
	err = UpdatePrescriptionRate(db, "foo", 555.55)
	if err.Error() != "prescription not found" {
		t.Error(err)
	}

	// Clean up.
	CloseDatabase(db)
	os.Remove(tempFile.Name())
}
