package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
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

// FOR LOCAL TESTING ONLY! CLEARS THE CURRENT DATABASE SO THE TESTS WORK AS INTENDED
/*
func TestClearDb(t *testing.T) {
	url := "http://localhost:8080/clear"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal("Failed to reach API:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 201, got %v", resp.StatusCode)
	}
}*/

func TestAPIUserPostEndpoint(t *testing.T) {
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
			Id:           1,
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

func TestAPIAdminsPostEndpoint(t *testing.T) {
	url := "http://localhost:8080/admins"
	jsonBody := []byte(`{"email":"admin@test.com", "password":"test"}`)

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
			Id:           2,
			Username:     "admin@test.com",
			Name:         "",
			Surname:      "",
			Password:     "test",
			Email:        "admin@test.com",
			Location:     "",
			Admin:        true,
			BlockedUser:  false,
			ProfilePhoto: 0,
			Description:  ""},
	}
	assert.Equal(t, responseBody, expectedBody)
}

func TestAPIUsersGet(t *testing.T) {
	url := "http://localhost:8080/users"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal("Failed to reach API:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %v", resp.StatusCode)
	}

	var responseBody ResponseUsers
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal("Failed to decode JSON response:", err)
	}

	var expectedUser1 = User{
		Id:           1,
		Username:     "test@test.com",
		Name:         "",
		Surname:      "",
		Password:     "test",
		Email:        "test@test.com",
		Location:     "",
		Admin:        false,
		BlockedUser:  false,
		ProfilePhoto: 0,
		Description:  "",
	}

	var expectedUser2 = User{
		Id:           2,
		Username:     "admin@test.com",
		Name:         "",
		Surname:      "",
		Password:     "test",
		Email:        "admin@test.com",
		Location:     "",
		Admin:        true,
		BlockedUser:  false,
		ProfilePhoto: 0,
		Description:  ""}

	users := append([]User{}, expectedUser1, expectedUser2)
	expectedResponse := ResponseUsers{
		Users: users,
	}

	assert.Equal(t, responseBody, expectedResponse)

}

func TestAPILoginPost(t *testing.T) {
	url := "http://localhost:8080/login"
	jsonBody := []byte(`{"email":"test@test.com", "password":"test"}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal("Failed to reach API:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %v", resp.StatusCode)
	}

	var responseBody ResponseUser

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal("Failed to decode JSON response:", err)
	}
	var expectedBody = ResponseUser{
		User: User{
			Id:           1,
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
	assert.Equal(t, expectedBody, responseBody)
}

func TestAPIUserGet(t *testing.T) {
	url := "http://localhost:8080/users/1"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal("Failed to reach API:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %v", resp.StatusCode)
	}

	var responseBody ResponseUser
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatal("Failed to decode JSON response:", err)
	}
	var expectedBody = ResponseUser{
		User: User{
			Id:           1,
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
	assert.Equal(t, expectedBody, responseBody)
}
