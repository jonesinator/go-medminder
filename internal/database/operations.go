package database

import (
	"database/sql"
	"errors"
	"time"
)

type Prescription struct {
	Name     string
	Quantity float64
	Rate     float64
	Updated  time.Time
}

func (p Prescription) ExpectedCount() float64 {
	daysSince := int(time.Since(p.Updated).Hours()) / int(24)
	amountConsumed := float64(daysSince) * p.Rate
	return max(p.Quantity-amountConsumed, 0)
}

func (p Prescription) RefillDate() time.Time {
	days := p.Quantity / p.Rate
	endDay := p.Updated.Add(time.Duration(days) * 24 * time.Hour)
	return endDay
}

func expectOneAffected(result sql.Result, err error, message string) error {
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New(message)
	}

	return nil
}

func CreatePrescription(d *Database, name string, quantity float64, rate float64) (*Prescription, error) {
	updated := time.Now().UTC()
	result, err := d.db.Exec(
		"INSERT INTO prescriptions (name, quantity, rate, updated) VALUES (?, ?, ?, ?)",
		name, quantity, rate, updated)
	if err != nil && err.Error() == "UNIQUE constraint failed: prescriptions.name" {
		return nil, errors.New("prescription already exists")
	}
	err = expectOneAffected(result, err, "prescription already exists")
	if err != nil {
		return nil, err
	}
	return &Prescription{name, quantity, rate, updated}, nil
}

func ReadPrescription(d *Database, name string) (*Prescription, error) {
	rows := d.db.QueryRow(
		"SELECT name, quantity, rate, updated FROM prescriptions WHERE name = ?", name)
	var rx Prescription
	err := rows.Scan(&rx.Name, &rx.Quantity, &rx.Rate, &rx.Updated)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("prescription not found")
		}
		return nil, err
	}
	return &rx, nil
}

func ReadAllPrescriptions(d *Database) ([]*Prescription, error) {
	rows, err := d.db.Query("SELECT name, quantity, rate, updated FROM prescriptions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []*Prescription
	for rows.Next() {
		var p Prescription
		err := rows.Scan(&p.Name, &p.Quantity, &p.Rate, &p.Updated)
		if err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, &p)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return prescriptions, nil
}

func UpdatePrescriptionQuantity(d *Database, name string, quantity float64) error {
	result, err := d.db.Exec(
		"UPDATE prescriptions SET quantity = ?, updated = ? WHERE name = ?",
		quantity, time.Now().UTC(), name)
	return expectOneAffected(result, err, "prescription not found")
}

func UpdatePrescriptionRate(d *Database, name string, rate float64) error {
	result, err := d.db.Exec("UPDATE prescriptions SET rate = ? WHERE name = ?", rate, name)
	return expectOneAffected(result, err, "prescription not found")
}

func DeletePrescription(d *Database, name string) error {
	result, err := d.db.Exec("DELETE FROM prescriptions WHERE name = ?", name)
	return expectOneAffected(result, err, "prescription not found")
}
