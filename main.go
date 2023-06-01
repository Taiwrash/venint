package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	tron "github.com/go-chain/go-tron"
)

// Define the request structure for creating and sending messages
type MessageRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

// Define the response structure for a successful message creation and sending
type MessageResponse struct {
	MessageID string `json:"messageId"`
}

// CreateAndSendMessage creates and sends a message to the Tron blockchain
func CreateAndSendMessage(request MessageRequest) (string, error) {
	// Perform necessary validation
	if len(request.Message) == 0 {
		return "", errors.New("Empty message")
	}

	// Create a new Tron client
	client := tron.NewEasyClient(api.DefaultFullNode)

	// Get the account address and private key from your configuration or wallet
	account := tron.Address(request.From)
	privateKey := "venom-private-key"

	// Create a new transaction
	transaction := tron.NewTransaction(
		account,
		tron.Address(request.To),
		0, // Set amount to 0 for message transactions
		&client,
	)

	// Set the message data
	transaction.SetMessage([]byte(request.Message))

	// Sign the transaction with the private key
	err := transaction.Sign(privateKey)
	if err != nil {
		return "", err
	}

	// Broadcast the transaction to the network
	transactionID, err := client.BroadcastTransaction(transaction)
	if err != nil {
		return "", err
	}

	return transactionID, nil
}

// Define the handler function for creating and sending a message
func createAndSendMessage(w http.ResponseWriter, r *http.Request) {
	var request MessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageID, err := CreateAndSendMessage(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := MessageResponse{
		MessageID: messageID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Define the handler function for processing messages reliably
func processMessage(w http.ResponseWriter, r *http.Request) {
	var request MessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Implement logic for processing the message reliably
	maxRetries := 3
	retryInterval := time.Second * 5

	// Retry processing the message for a certain number of times
	for i := 0; i < maxRetries; i++ {
		err := processSingleMessage(request)
		if err == nil {
			// Processing successful
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log the error or handle it as per your requirement
		log.Printf("Error processing message: %s. Retrying...", err)

		// Wait for the retry interval before attempting again
		time.Sleep(retryInterval)
	}

	// Message processing failed even after retries
	w.WriteHeader(http.StatusInternalServerError)
}

// processSingleMessage processes a single message
func processSingleMessage(request MessageRequest) error {
	// Check if the message has expired
	if isMessageExpired(request) {
		return errors.New("Message has expired")
	}

	// Thinking of integrating some sdk or call some APIs but No Idea for now
	log.Printf("Processing message: %s", request.Message)

	return nil
}

// isMessageExpired checks if the given message has expired
func isMessageExpired(request MessageRequest) bool {
	// Get the current time
	now := time.Now()

	// Parse the message creation time
	creationTime, err := time.Parse(time.RFC3339, request.CreatedAt)
	if err != nil {
		// Unable to parse the creation time, consider it expired
		return true
	}

	// Calculate the duration since the message was created
	duration := now.Sub(creationTime)

	// Compare the duration with the maximum expiration duration
	if duration > maxMessageExpiration {
		return true
	}

	return false
}


// Define the handler function for getting account state
func getAccountState(w http.ResponseWriter, r *http.Request) {
	// Add your implementation here for getting account state
}

// Define the handler function for querying blockchain data
func queryBlockchainData(w http.ResponseWriter, r *http.Request) {
	// Add your implementation here for querying blockchain data
}

// Define the handler function for subscribing to events and updates
func subscribeToUpdates(w http.ResponseWriter, r *http.Request) {
	// Add your implementation here for subscribing to events and updates
}

// ... Define additional handler functions for other functionalities ...

func main() {
	// Create a new router using gorilla/mux
	router := mux.NewRouter()

	// Define the API endpoints and their corresponding handler functions
	router.HandleFunc("/message", createAndSendMessage).Methods("POST")
	router.HandleFunc("/message/process", processMessage).Methods("POST")
	router.HandleFunc("/account", getAccountState).Methods("GET")
	router.HandleFunc("/blockchain/query", queryBlockchainData).Methods("GET")
	router.HandleFunc("/updates/subscribe", subscribeToUpdates).Methods("POST")

	// Start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
