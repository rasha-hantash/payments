#!/bin/bash

# Set the base URL for your API
BASE_URL="http://localhost:8080"  # Change this to your actual API base URL

# Function to make a POST request
make_post_request() {
    endpoint=$1
    data=$2
    echo "Making POST request to $endpoint"
    curl -X POST -H "Content-Type: application/json" -d "$data" "$BASE_URL$endpoint"
    echo -e "\n"
}

# Function to make a GET request
make_get_request() {
    endpoint=$1
    params=$2
    echo "Making GET request to $endpoint"
    curl -X GET "$BASE_URL$endpoint$params"
    echo -e "\n"
}

# Create a user
user_data='{"name": "John Doe", "email": "john@example.com"}'
make_post_request "/create_user" "$user_data"

# Create an account
account_data='{"user_id": "user_123", "account_type": "checking"}'
make_post_request "/create_account" "$account_data"

# Deposit funds
deposit_data='{"account_id": "account_123", "amount": 1000}'
make_post_request "/deposit_funds" "$deposit_data"

# Withdraw funds
withdraw_data='{"account_id": "account_123", "amount": 500}'
make_post_request "/withdraw_funds" "$withdraw_data"

# Transfer funds
transfer_data='{"from_account_id": "account_123", "to_account_id": "account_456", "amount": 250}'
make_post_request "/transfer_funds" "$transfer_data"

# List transactions
make_get_request "/list_transactions" "?account_id=account_123&limit=10"

# Get account balance
make_get_request "/get_account_balance" "?account_id=account_123"

echo "All requests completed."