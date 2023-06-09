package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndSendMessage(t *testing.T) {
	// Create a mock request body
	requestBody := []byte(`{"from": "sender", "to": "recipient", "message": "Hello, World!"}`)

	// Create a new request with the mock body
	req, err := http.NewRequest("POST", "/message", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(createAndSendMessage)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	// Parse the response body
	var response MessageResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response data
	if response.MessageID == "" {
		t.Error("expected non-empty message ID")
	}
	if response.Status != "pending" {
		t.Errorf("expected status pending; got %s", response.Status)
	}
}

func TestProcessMessage(t *testing.T) {
	// Create a mock request with the message ID
	req, err := http.NewRequest("POST", "/message/process?messageId=abc123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(processMessage)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status Accepted; got %d", rr.Code)
	}
}

func TestGetAccountState(t *testing.T) {
	// Create a mock request with the account address
	req, err := http.NewRequest("GET", "/account/state?accountAddress=abc123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(getAccountState)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status InternalServerError; got %d", rr.Code)
	}
}

func TestQueryBlockchainData(t *testing.T) {
	// Create a mock request
	req, err := http.NewRequest("GET", "/blockchain/data?dataType=blocks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(queryBlockchainData)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	// Parse the response body
	var response interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Add additional assertions for the response data if needed
}

func TestSubscribeToUpdates(t *testing.T) {
	// Create a mock request
	req, err := http.NewRequest("POST", "/updates/subscribe?eventType=messages", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(subscribeToUpdates)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}

	// Check the response body
	expectedBody := "Subscribed to updates successfully"
	if rr.Body.String() != expectedBody {
		t.Errorf("expected body '%s'; got '%s'", expectedBody, rr.Body.String())
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()

	// Perform any necessary cleanup or teardown here
	db.Close()

	// Exit with the test code
	os.Exit(code)
}
