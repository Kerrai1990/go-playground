package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "password", "todo")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s instead", body)
	}
}

func TestNonExistingUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users/99", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be 'User not found'. Got %s instead", m["error"])
	}
}
func TestCreatNewUser(t *testing.T) {
	clearTable()

	payload := []byte(`{"name" : "test", "age" : 30}`)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test" {
		t.Errorf("Expected name to be 'test'. got %s instead", m["name"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected ID to be 1.0. Got %f instead", m["id"])
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	return response
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected a %d, but got a %d", expected, actual)
	}
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL
)`
