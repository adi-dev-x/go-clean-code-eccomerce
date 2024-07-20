package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	clientID     = "AdZ4ZRsP5DaT-B8DT4mfWTEY-c42yF9DHnJvju5e3tMSrta6IAP2SXD8-xkknlU8uDBbXG4E0cvWkXKU"
	clientSecret = "EKnlirnIYfYljO331c6APaxN9GpZLumceWVJywzgUowrwE7gMXmVpVQlUZ3mnC1gWqKFZPI2b7H4qI3m"
)

const (
	PayPalAPIBaseURL = "https://api.sandbox.paypal.com" // Sandbox base URL, change to live URL for production
)

type AccessTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
}

type CreatePaymentRequest struct {
	Intent       string        `json:"intent"`
	Payer        Payer         `json:"payer"`
	Transactions []Transaction `json:"transactions"`
	RedirectURLs RedirectURLs  `json:"redirect_urls"`
}

type Payer struct {
	PaymentMethod string `json:"payment_method"`
}

type Transaction struct {
	Amount Amount `json:"amount"`
}

type Amount struct {
	Total    string `json:"total"`
	Currency string `json:"currency"`
}

type RedirectURLs struct {
	ReturnURL string `json:"return_url"`
	CancelURL string `json:"cancel_url"`
}

type CreatePaymentResponse struct {
	ID    string `json:"id"`
	Links []Link `json:"links"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

func main() {
	http.HandleFunc("/create_payment", handleCreatePayment)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	// Get access token
	accessToken, err := getAccessToken()
	if err != nil {
		log.Println("Error getting access token:", err)
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// Prepare request body for creating payment
	createPaymentReq := CreatePaymentRequest{
		Intent: "sale",
		Payer: Payer{
			PaymentMethod: "paypal",
		},
		Transactions: []Transaction{
			{
				Amount: Amount{
					Total:    "10.00", // Adjust amount as needed
					Currency: "USD",
				},
			},
		},
		RedirectURLs: RedirectURLs{
			ReturnURL: "http://localhost:8080/execute_payment",
			CancelURL: "http://localhost:8080/cancel_payment",
		},
	}

	// Convert request body to JSON
	reqBody, err := json.Marshal(createPaymentReq)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Create HTTP request to PayPal API
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/payments/payment", PayPalAPIBaseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	var createPaymentResp CreatePaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&createPaymentResp)
	if err != nil {
		log.Println("Error decoding response:", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	// Redirect to PayPal approval URL
	for _, link := range createPaymentResp.Links {
		if link.Rel == "approval_url" {
			http.Redirect(w, r, link.Href, http.StatusFound)
			return
		}
	}

	http.Error(w, "No approval_url found in response", http.StatusInternalServerError)
}

func getAccessToken() (string, error) {
	// Create basic auth string
	authString := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Create HTTP request to get access token
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/oauth2/token", PayPalAPIBaseURL), nil)
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authString)

	// Set form data
	q := req.URL.Query()
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	// Perform HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	var accessTokenResp AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&accessTokenResp)
	if err != nil {
		return "", err
	}

	return accessTokenResp.AccessToken, nil
}
