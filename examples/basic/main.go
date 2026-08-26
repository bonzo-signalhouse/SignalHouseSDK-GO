package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bonzo-signalhouse/SignalHouseSDK-GO/signalhouse"
)

func main() {
	// Initialize the client
	client, err := signalhouse.NewClient(signalhouse.Config{
		APIKey:  "your-api-key",
		BaseURL: "https://api.signalhouse.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Example: Login
	loginResp, err := client.Auth.Login(ctx, signalhouse.LoginParams{
		Email:    "user@example.com",
		Password: "password123",
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Login success: %v, status: %d\n", loginResp.Success, loginResp.Status)

	// Example: Use a token from login for subsequent requests
	var loginData struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(loginResp.Data, &loginData); err != nil {
		log.Fatal(err)
	}

	opts := &signalhouse.RequestOptions{Token: loginData.Token}

	// Example: Get phone numbers
	numbersResp, err := client.Numbers.GetPhoneNumbers(ctx, signalhouse.GetPhoneNumbersParams{
		GroupID: "G12345678",
		Page:   "1",
		Limit:  "10",
	}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Numbers response success: %v\n", numbersResp.Success)

	// Example: Send an SMS
	smsResp, err := client.Messages.SendSMS(ctx, signalhouse.SendSMSParams{
		SenderPhoneNumber:     "15551234567",
		RecipientPhoneNumbers: []string{"15559876543"},
		MessageBody:           "Hello from SignalHouse Go SDK!",
	}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SMS sent: %v, status: %d\n", smsResp.Success, smsResp.Status)

	// Example: Get wallet balance
	walletResp, err := client.Billing.GetWallet(ctx, signalhouse.GetWalletParams{
		GroupID: "G12345678",
	}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wallet response: %s\n", string(walletResp.Data))

	// Example: Search available phone numbers
	availableResp, err := client.Numbers.GetAvailablePhoneNumbers(ctx, signalhouse.GetAvailablePhoneNumbersParams{
		Country: "US",
		State:   "CA",
		Limit:   "5",
	}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Available numbers: %s\n", string(availableResp.Data))
}
