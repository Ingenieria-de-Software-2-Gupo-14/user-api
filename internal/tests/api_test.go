package tests

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	. "ing-soft-2-tp1/internal/models"
	"net/http"
	"testing"
)

func TestAPIHealthEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatal("Failed to reach API:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", resp.StatusCode)
	}
}

func TestAPIUserPostEndopint(t *testing.T) {
	url := "http://localhost:8080/users"
	jsonBody := []byte(`{"email":"test@test.com", "password":"test"}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal("Failed to reach API:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %v", resp.StatusCode)
	}

	var responseBody ResponseUser

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal("Failed to decode JSON response:", err)
	}
	var expectedBody = ResponseUser{
		User: User{
			Id:           0,
			Username:     "test@test.com",
			Name:         "",
			Surname:      "",
			Password:     "test",
			Email:        "test@test.com",
			Location:     "",
			Admin:        false,
			BlockedUser:  false,
			ProfilePhoto: 0,
			Description:  ""},
	}
	assert.Equal(t, responseBody, expectedBody)
}
