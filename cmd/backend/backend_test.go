package main

import (
	"bytes"
	"encoding/json"
	"jonesinator/go-medminder/internal/database"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAPI(t *testing.T) {
	os.Remove("./test.db")
	db, err := database.OpenDatabase("./test.db")
	if err != nil {
		t.Error("unable to create db")
	}
	GlobalDB = db

	router := SetupRouter()

	data := gin.H{"quantity": 123.45, "rate": 543.21}
	jsonValue, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", "/rx/foo", bytes.NewBuffer(jsonValue))
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Error("unexpected status")
	}

	request, _ = http.NewRequest("GET", "/rx", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("unexpected status")
	}

	var names []string
	_ = json.Unmarshal(recorder.Body.Bytes(), &names)
	if len(names) != 1 || names[0] != "foo" {
		t.Error("unexpected response")
	}

	request, _ = http.NewRequest("GET", "/rx/foo", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("unexpected status")
	}

	var prescription database.Prescription
	_ = json.Unmarshal(recorder.Body.Bytes(), &prescription)
	if prescription.Quantity != 123.45 || prescription.Rate != 543.21 {
		t.Error("unexpected response")
	}

	data = gin.H{"quantity": 333.33}
	jsonValue, _ = json.Marshal(data)
	request, _ = http.NewRequest("PATCH", "/rx/foo", bytes.NewBuffer(jsonValue))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("unexpected status")
	}

	data = gin.H{"rate": 222.22}
	jsonValue, _ = json.Marshal(data)
	request, _ = http.NewRequest("PATCH", "/rx/foo", bytes.NewBuffer(jsonValue))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("unexpected status")
	}

	request, _ = http.NewRequest("DELETE", "/rx/foo", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Error("unexpected status")
	}

	os.Remove("./test.db")
}
