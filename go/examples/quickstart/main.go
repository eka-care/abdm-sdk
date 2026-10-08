package main

import (
	"context"
	"fmt"
	"log"

	ekasdk "github.com/eka-care/eka-sdk-go"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
	"github.com/eka-care/eka-sdk-go/services/abdm/abha/login"
)

func main() {
	ctx := context.Background()

	// Step 1: Create SDK client from environment variables
	// The SDK looks for credentials in this order:
	// 1. EKA_CLIENT_ID and EKA_CLIENT_SECRET (required)
	// 2. Explicit configuration via options
	//
	// Set these environment variables:
	// export EKA_ENVIRONMENT=production
	// export EKA_CLIENT_ID=your-client-id
	// export EKA_CLIENT_SECRET=your-client-secret
	client := ekasdk.NewFromEnv()

	// Step 2: Do client login (authenticate with Eka platform)
	if err := client.Login(ctx); err != nil {
		log.Fatalf("❌ Client authentication failed: %v", err)
	}
	fmt.Println("✅ Client authenticated with Eka Care platform!")

	// Step 3: Use ABDM login APIs
	headers := interfaces.Headers{
		PatientID:     "eka-user-oid",
		PartnerUserID: "your-user-id",
		HipID:         "your-hip-id",
	}

	// Generate OTP for ABDM login
	fmt.Println("📱 Generating OTP for ABDM login...")
	otpReq := &login.InitLoginRequest{
		Identifier: "demo@abdm",
		Method:     login.LoginMethodPhrAddress,
	}

	otpResp, err := client.ABDM.Login().LoginInit(ctx, headers, otpReq)
	if err != nil {
		log.Printf("⚠️  OTP generation failed: %v", err)
		return
	}
	fmt.Printf("✅ OTP generated. Transaction ID: %s\n", otpResp.TxnID)

	// Verify OTP for ABDM login
	fmt.Println("� Verifying OTP...")
	verifyReq := &login.VerifyLoginOTPRequest{
		OTP:   "123456", // In real app, get this from user input
		TxnID: otpResp.TxnID,
	}

	verifyResp, err := client.ABDM.Login().LoginVerify(ctx, headers, verifyReq)
	if err != nil {
		log.Printf("⚠️  OTP verification failed: %v", err)
		return
	}
	fmt.Printf("✅ ABDM login successful! Transaction ID: %s\n", verifyResp.TxnID)

	fmt.Println("🎉 Complete flow: Client login → ABDM login successful!")
}
