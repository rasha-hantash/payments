package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"net/url"
	// "time"
)

const baseURL = "http://localhost:8080" // Change this to your actual API base URL

// Colors for output
const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	YELLOW = "\033[1;33m"
	NC     = "\033[0m" // No Color
)

// User struct to hold user response data
type User struct {
	ID                 string `json:"id"`
	IntLedgerAccountID string `json:"int_ledger_account_id"`
	ExtLedgerAccountID string `json:"ext_ledger_account_id"`
}

// Account struct to hold account response data
type Account struct {
	ID string `json:"id"`
}

// Balance struct to hold balance response data
type Balance struct {
	Balance float64 `json:"balance"`
}

// Generic response struct
type Response struct {
	Balance float64 `json:"balance,omitempty"`
	ID      string  `json:"id,omitempty"`
}

// Function to print error messages
func fail(message string) {
	fmt.Fprintf(os.Stderr, "%sError: %s%s\n", RED, message, NC)
}

// Function to print success messages
func success(message string) {
	fmt.Printf("%sSuccess: %s%s\n", GREEN, message, NC)
}

// Function to print info messages
func info(message string) {
	fmt.Printf("%sInfo: %s%s\n", YELLOW, message, NC)
}

// Function to make a POST request and return the response
func makePostRequest(endpoint string, data interface{}) ([]byte, error) {
	url := baseURL + endpoint
	info(fmt.Sprintf("Making POST request to %s", url))
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d, response: %s", resp.StatusCode, body)
	}

	return body, nil
}


// Function to make a POST request and return the response

// Function to make a GET request and return the response
func makeGetRequest(endpoint string, params map[string]interface{}) ([]byte, error) {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value.(string))
	}
	u.RawQuery = q.Encode()

	info(fmt.Sprintf("Making GET request to %s", u.String()))
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d, response: %s", resp.StatusCode, body)
	}

	return body, nil
}

// Main execution
func main() {
	info("Starting banking API test script")
	info(fmt.Sprintf("Using base URL: %s", baseURL))

	// Create a user and store the user_id
	userData := map[string]string{"name": "John Doe", "email": "john@example.com"}
	userResponse, err := makePostRequest("/create_user", userData)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	var user User
	err = json.Unmarshal(userResponse, &user)
	if err != nil {
		log.Fatalf("Failed to parse user response: %v", err)
	}

	success(fmt.Sprintf("Created user with ID: %s. Internal Ledger Account ID: %s. External Ledger Account ID: %s", user.ID, user.IntLedgerAccountID, user.ExtLedgerAccountID))

	// Create an account using the stored user_id
	accountData := map[string]string{"account_type": "internal"}
	accountResponse, err := makePostRequest("/create_account", accountData)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}

	var account Account
	err = json.Unmarshal(accountResponse, &account)
	if err != nil {
		log.Fatalf("Failed to parse account response: %v", err)
	}

	success(fmt.Sprintf("Created an internal account with ID: %s", account.ID))

	// Deposit funds using the stored account_id
	depositData := map[string]interface{}{
		"debit_account_id":   user.ExtLedgerAccountID,
		"credit_account_id":  user.IntLedgerAccountID,
		"user_id":            user.ID,
		"idempotency_key":    "dep_idempotent_john@example.com",
		"amount":             1000,
	}

	depositResponse, err := makePostRequest("/deposit_funds", depositData)
	if err != nil {
		log.Fatalf("Failed to deposit funds: %v", err)
	}
	success(fmt.Sprintf("Deposited funds: %s", string(depositResponse)))

	// Deposit funds again to test idempotency
	depositResponse2, err := makePostRequest("/deposit_funds", depositData)
	if err != nil {
		log.Fatalf("Failed to deposit funds again: %v", err)
	}
	success(fmt.Sprintf("Second deposit attempt: %s", string(depositResponse2)))

	// Get balance to check idempotency
	getBalanceData := map[string]interface{}{
		"account_id": user.IntLedgerAccountID,
		// "at_time":   "2024-09-01T00:00:00Z",
	}

	getBalanceResponse, err := makeGetRequest("/get_account_balance", getBalanceData)
	if err != nil {
		log.Fatalf("Failed to get account balance: %v", err)
	}
	success(fmt.Sprintf("Account balance after deposits: %s", string(getBalanceResponse)))

	// Transfer funds
	transferData := map[string]interface{}{
		"debit_account_id":  user.IntLedgerAccountID,
		"credit_account_id": account.ID,
		"amount":            250,
	}
	transferResponse, err := makePostRequest("/transfer_funds", transferData)
	if err != nil {
		log.Fatalf("Failed to transfer funds: %v", err)
	}
	success(fmt.Sprintf("Transferred funds: %s", string(transferResponse)))

	// Withdraw funds
	withdrawData := map[string]interface{}{
		"account_id": user.IntLedgerAccountID,
		"amount":     500,
	}
	withdrawResponse, err := makePostRequest("/withdraw_funds", withdrawData)
	if err != nil {
		log.Fatalf("Failed to withdraw funds: %v", err)
	}
	success(fmt.Sprintf("Withdrew funds: %s", string(withdrawResponse)))

	// List transactions
	listTransactionsData := map[string]interface{}{
		"account_id": user.IntLedgerAccountID,
		"limit":      10,
	}
	transactionsResponse, err := makePostRequest("/list_transactions", listTransactionsData)
	if err != nil {
		log.Fatalf("Failed to list transactions: %v", err)
	}
	success(fmt.Sprintf("List of transactions: %s", string(transactionsResponse)))

	// Get final account balance
	finalBalanceData := map[string]interface{}{
		"account_id": user.IntLedgerAccountID,
	}
	finalBalanceResponse, err := makePostRequest("/get_account_balance", finalBalanceData)
	if err != nil {
		log.Fatalf("Failed to get final account balance: %v", err)
	}
	success(fmt.Sprintf("Final account balance: %s", string(finalBalanceResponse)))

	info("All requests completed.")

	// Print summary
	var finalBalance Balance
	err = json.Unmarshal(finalBalanceResponse, &finalBalance)
	if err != nil {
		log.Fatalf("Failed to parse final balance response: %v", err)
	}

	fmt.Println()
	info("Summary of operations:")
	fmt.Printf("- Created user: %s\n", user.ID)
	fmt.Printf("- Created internal account: %s\n", account.ID)
	fmt.Printf("- Deposited funds: 1000\n")
	fmt.Printf("- Transferred funds: 250\n")
	fmt.Printf("- Withdrew funds: 500\n")
	fmt.Printf("- Final balance: %.2f\n", finalBalance.Balance)
}
