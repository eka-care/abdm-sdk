# Authentication

The SDK handles authentication automatically using OAuth 2.0 client credentials flow with automatic token management, following industry best practices.

## Quick Start

### Method 1: Environment Variables (Recommended)

1. Set your credentials:
```bash
export EKA_CLIENT_ID="your-client-id"
export EKA_CLIENT_SECRET="your-client-secret"
export EKA_ENVIRONMENT="production"  # or "development"
```

2. Use the SDK:
```go
client := ekasdk.NewFromEnv()
err := client.Login(ctx)  // Automatically authenticates
```

### Method 2: Explicit Configuration

```go
client := ekasdk.New(
    ekasdk.WithClientID("your-client-id"),
    ekasdk.WithClientSecret("your-client-secret"),
    ekasdk.WithEnvironment(ekasdk.EnvironmentProduction),
)
err := client.Login(ctx)
```

## How Authentication Works

### 1. Credential Resolution
The SDK looks for credentials in this order:
1. `EKA_CLIENT_ID` and `EKA_CLIENT_SECRET` environment variables
2. Explicit configuration via options

### 2. OAuth 2.0 Client Credentials Flow
- The SDK exchanges your client credentials for access tokens
- Tokens are automatically cached and refreshed before expiration
- All API calls include the current valid token

### 3. Automatic Token Management
- **Initial Authentication**: Happens during `client.Login(ctx)`
- **Token Refresh**: Automatic when tokens are about to expire
- **Retry Logic**: Built-in retry for authentication failures
- **Thread Safety**: Concurrent access is handled safely

## Advanced Usage

### Custom Credentials Provider

```go
// Create a custom credentials provider
provider := auth.NewStaticCredentialsProvider(
    "access-token", 
    "refresh-token", 
    3600,  // expires in seconds
    7200,  // refresh expires in seconds
)

client := ekasdk.New(
    ekasdk.WithCredentialsProvider(provider),
)
```

### Manual Token Management

```go
// Get current credentials
creds, err := client.GetCredentials(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Access Token: %s\n", creds.AccessToken)
fmt.Printf("Expires At: %s\n", creds.ExpiresAt)
```

### Custom Authentication Flow

```go
// Create a custom login request
loginReq := &auth.ClientLoginRequest{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
}

// Create and set a custom provider
provider := client.NewClientCredentialsProvider(loginReq)
client.SetCredentialsProvider(provider)

// Authenticate
creds, err := provider.Retrieve(ctx)
```

## Error Handling

The SDK provides detailed error messages for authentication issues:

```go
client := ekasdk.NewFromEnv()
if err := client.Login(ctx); err != nil {
    switch {
    case strings.Contains(err.Error(), "client ID is required"):
        // Missing EKA_CLIENT_ID environment variable
        log.Fatal("Set EKA_CLIENT_ID environment variable")
        
    case strings.Contains(err.Error(), "client secret is required"):
        // Missing EKA_CLIENT_SECRET environment variable
        log.Fatal("Set EKA_CLIENT_SECRET environment variable")
        
    case strings.Contains(err.Error(), "failed to authenticate with provided credentials"):
        // Invalid credentials
        log.Fatal("Check your client ID and secret in the developer portal")
        
    default:
        log.Fatalf("Authentication failed: %v", err)
    }
}
```

## Security Features

### Automatic Token Refresh
- Tokens are refreshed 5 minutes before expiration
- Refresh happens automatically in the background
- Failed refresh triggers re-authentication

### Secure Storage
- Tokens are stored in memory only
- No persistent storage of credentials
- Automatic cleanup on client disposal

### Thread Safety
- Multiple goroutines can safely use the same client
- Concurrent authentication requests are handled correctly
- Token refresh is protected by mutexes

## Best Practices

✅ **Do:**
- Use environment variables for credentials
- Call `client.Login(ctx)` once after client creation
- Let the SDK handle token management automatically
- Use different credentials for different environments

❌ **Don't:**
- Hardcode credentials in source code
- Manually manage token refresh
- Share clients between different credential sets
- Call `Login()` repeatedly unless re-authentication is needed

## Compatibility

This authentication system works with:
- **All Eka Care APIs**: ABDM, Profile, and future services
- **All environments**: Production and development
- **All deployment methods**: Docker, Kubernetes, serverless functions
- **All Go versions**: 1.19+

The authentication is completely transparent to your application code - once authenticated, all API calls automatically include the proper authorization headers.
