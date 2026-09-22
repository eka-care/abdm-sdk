# Eka Care SDK for Go

A Go SDK for integrating with [Eka Care's healthcare APIs](https://developer.eka.care), including ABDM (Ayushman Bharat Digital Mission) services.

## Getting Started

### 1. Get Your Credentials

1. Register at [Eka Care Developer Portal](https://developer.eka.care)
2. Get your `client_id` and `client_secret` from the portal

### 2. Install the SDK

```bash
go get github.com/eka-care/eka-sdk-go
```

### 3. Set Up Authentication

The SDK uses environment variables for secure credential management following industry best practices.

#### Recommended Method (Environment Variables)

Create a `.env` file or set environment variables:

```bash
# Required for production
export EKA_ENVIRONMENT=production
export EKA_CLIENT_ID=your-client-id
export EKA_CLIENT_SECRET=your-client-secret

# For development/testing
export EKA_ENVIRONMENT=development
export EKA_CLIENT_ID=your-dev-client-id
export EKA_CLIENT_SECRET=your-dev-client-secret
```

#### Alternative Method

**Explicit Configuration (Not Recommended for Production)**
```go
client := ekasdk.New(
    ekasdk.WithEnvironment(ekasdk.EnvironmentProduction),
    ekasdk.WithClientID("your-client-id"),
    ekasdk.WithClientSecret("your-client-secret"),
)
```

**Security Best Practices:**
- ✅ Use environment variables in production
- ✅ Never commit credentials to version control
- ✅ Use different credentials for dev/staging/production
- ❌ Avoid hardcoding credentials in source code

### 4. Quick Start Example

#### Method 1: Environment Variables (Recommended)

```go
package main

import (
    "context"
    "log"
    
    ekasdk "github.com/eka-care/eka-sdk-go"
    "github.com/eka-care/eka-sdk-go/services/abdm/abha/login"
    "github.com/eka-care/eka-sdk-go/internal/interfaces"
)

func main() {
    ctx := context.Background()

    // Create SDK client from environment variables
    // Automatically reads EKA_CLIENT_ID, EKA_CLIENT_SECRET, EKA_ENVIRONMENT
    client := ekasdk.NewFromEnv()

    // Authenticate with Eka platform
    if err := client.Login(ctx); err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    
    // Now you can use any Eka Care API
    headers := interfaces.Headers{
		PatientID:     "eka-user-oid",
		PartnerUserID: "your-user-id",
		HipID:         "your-hip-id",
	}
    
    // Example: ABDM login
    otpReq := &login.InitLoginRequest{
        Identifier: "demo@abdm",
        Method:     login.LoginMethodPhrAddress,
    }
    
    otpResp, err := client.ABDM.Login().LoginInit(ctx, headers, otpReq)
    if err != nil {
        log.Printf("ABDM login failed: %v", err)
        return
    }
    
    log.Printf("Success! Transaction ID: %s", otpResp.TxnID)
}
```

#### Method 2: Explicit Configuration

```go
package main

import (
    "context"
    "log"
    
    ekasdk "github.com/eka-care/eka-sdk-go"
)

func main() {
    ctx := context.Background()

    // Create SDK client with explicit configuration
    client := ekasdk.New(
        ekasdk.WithEnvironment(ekasdk.EnvironmentProduction),
        ekasdk.WithClientID("your-client-id"),
        ekasdk.WithClientSecret("your-client-secret"),
    )

    // Authenticate with Eka platform
    if err := client.Login(ctx); err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    
    // Use the APIs as normal...
}
```

## Configuration

### Environment Variables

The SDK supports the following environment variables:

#### Authentication (Required)
```bash
EKA_CLIENT_ID       # Your client ID from developer portal
EKA_CLIENT_SECRET   # Your client secret from developer portal
EKA_ENVIRONMENT     # "production" or "development"
```

#### Optional Configuration
```bash
EKA_TIMEOUT         # Request timeout in seconds (default: 30)
EKA_MAX_RETRIES     # Maximum retry attempts (default: 3)
EKA_USER_AGENT      # Custom User-Agent header
EKA_LOG_LEVEL       # Logging level: "debug", "info", "warn", "error"
EKA_DISABLE_SSL     # Disable SSL verification (default: false)
EKA_REGION          # API region (default: "us")
```

### Configuration Priority

The SDK resolves configuration in this order:
1. **Environment variables** (highest priority)
2. **Explicit options** via `WithXxx()` functions
3. **Default values** (lowest priority)

### Error Messages

The SDK provides clear error messages for common configuration issues:

```go
client := ekasdk.NewFromEnv()
if err := client.Login(ctx); err != nil {
    // Example error messages:
    // "client ID is required for authentication. Set EKA_CLIENT_ID environment variable or use WithClientID() option"
    // "client secret is required for authentication. Set EKA_CLIENT_SECRET environment variable or use WithClientSecret() option"
    // "failed to authenticate with provided credentials: invalid client credentials"
    log.Fatalf("Authentication failed: %v", err)
}
```

## Available Services

Once authenticated, you can access:

- **ABDM Services**: `client.ABDM.Login()`, `client.ABDM.Registration()`, `client.ABDM.Profile()`
- **More services** will be added as they become available

## Need Help?

- **Documentation**: [developer.eka.care](https://developer.eka.care)
- **API Reference**: [API Docs](https://developer.eka.care/api-reference)
- **Issues**: [GitHub Issues](https://github.com/eka-care/eka-sdk-go/issues)