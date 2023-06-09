package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// tron "github.com/TRON-US/go-tron"
	// "github.com/TRON-US/go-tron/api"
	tron "github.com/go-chain/go-tron"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	Status    string `json:"status"`
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
	privateKey := "your-private-key"

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

	func GetTransactionStatus(transactionID string) (string, error) {
		// Create a new Tron client
		client := tron.NewEasyClient(api.DefaultFullNode)

		// Get the transaction status from the blockchain.
		// transactionStatus, err := client.GetTransactionStatus(transactionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Update the message status in the database.
		err = updateMessageStatus(messageID, transactionStatus)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Write a success message to the response writer.
		response := MessageResponse{
			MessageID: messageID,
			Status:    transactionStatus,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
}
func processMessage(w http.ResponseWriter, r *http.Request) {
	// Get the message ID from the request.
	messageID := r.URL.Query().Get("messageId")

	// Get the transaction status from the blockchain.
	transactionStatus, err := client.GetTransactionStatus(messageID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the message status in the database.
	err = updateMessageStatus(messageID, transactionStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check the transaction status.
	switch transactionStatus {
	case "pending":
		// The message is still pending confirmation.
		w.WriteHeader(http.StatusAccepted)
	case "confirmed":
		// The message has been confirmed.
		w.WriteHeader(http.StatusOK)
	case "failed":
		// The message has failed.
		w.WriteHeader(http.StatusBadRequest)
	default:
		// Unknown transaction status.
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Update the message status in the database
func updateMessageStatus(messageID, status string) error {
	// Replace the database configuration with your actual database connection
	dsn := "user:password@tcp(localhost:3306)/database"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	defer db.Close()

	// Find the message record by message ID
	var message Message
	result := db.First(&message, "id = ?", messageID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Handle the case when the message record is not found
			return errors.New("Message not found")
		}
		// Handle other database errors
		return result.Error
	}

	// Update the message status
	message.Status = status
	result = db.Save(&message)
	if result.Error != nil {
		// Handle the database save error
		return result.Error
	}

	log.Println("Updated message status:", messageID, "->", status)
	return nil
}

func getAccountState(w http.ResponseWriter, r *http.Request) {
	// Get the account address from the request.
	accountAddress := r.URL.Query().Get("accountAddress")

	// Get the account state from the blockchain.
	accountState, err := client.GetAccountState(accountAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the account state to the response writer.
	json.NewEncoder(w).Encode(accountState)
}

// Define the handler function for querying blockchain data
func queryBlockchainData(w http.ResponseWriter, r *http.Request) {
	// Perform necessary validation
	// Check if the requested data type is valid
	dataType := r.URL.Query().Get("dataType")
	if dataType != "blocks" && dataType != "transactions" && dataType != "messages" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data type"})
		return
	}

	// Query the blockchain data based on the requested data type
	var data interface{}
	var err error

	switch dataType {
	case "blocks":
		// Query blocks data
		data, err = client.GetBlocks(10)
	case "transactions":
		// Query transactions data
		data, err = client.GetTransactions(10)
	case "messages":
		// Query messages data
		data, err = client.GetMessages(10)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to query blockchain data"})
		return
	}

	// Write the data to the response writer
	json.NewEncoder(w).Encode(data)
}


func subscribeToUpdates(w http.ResponseWriter, r *http.Request) {
	// Get the event or update type from the request.
	eventType := r.URL.Query().Get("eventType")

	// Subscribe to the event or update.
	err := client.SubscribeToUpdates(eventType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write a success message to the response writer.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Subscribed to updates successfully"))
}
func main() {
	// Create a new router using gorilla/mux
	router := mux.NewRouter()

	// Define the API endpoints and their corresponding handler functions
	router.HandleFunc("/message", createAndSendMessage).Methods("POST")
	router.HandleFunc("/message/process", processMessage).Methods("POST")
	router.HandleFunc("/account/state", getAccountState).Methods("GET")
	router.HandleFunc("/blockchain/data", queryBlockchainData).Methods("GET")
	router.HandleFunc("/updates/subscribe", subscribeToUpdates).Methods("POST")



	// Start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}