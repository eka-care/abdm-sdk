# Environment Configuration

The Eka SDK supports comprehensive environment-based configuration following industry best practices. This approach provides secure, flexible, and deployment-friendly credential management.

## Quick Setup

Set these environment variables:

```bash
export EKA_CLIENT_ID="your-client-id"
export EKA_CLIENT_SECRET="your-client-secret"
export EKA_ENVIRONMENT="production"  # or "development"
```

Then use the SDK:

```go
client := ekasdk.NewFromEnv()
err := client.Login(ctx)
```

## Complete Environment Variables Reference

### Authentication (Required)

| Variable | Description | Example |
|----------|-------------|---------|
| `EKA_CLIENT_ID` | Your client ID from developer portal | `abc123` |
| `EKA_CLIENT_SECRET` | Your client secret from developer portal | `secret456` |
| `EKA_ENVIRONMENT` | Target environment | `production` or `development` |

### Optional Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `EKA_TIMEOUT` | Request timeout in seconds | `30` | `60` |
| `EKA_MAX_RETRIES` | Maximum retry attempts | `3` | `5` |
| `EKA_USER_AGENT` | Custom User-Agent header | `eka-sdk-go/1.0.0` | `MyApp/1.0` |
| `EKA_LOG_LEVEL` | Logging level | `info` | `debug` |
| `EKA_DISABLE_SSL` | Disable SSL verification | `false` | `true` |
| `EKA_REGION` | API region | `us` | `us` |

## Environments

| Environment | Base URL | Use For |
|-------------|----------|---------|
| `production` | `https://api.eka.care` | Live applications |
| `development` | `https://api-dev.eka.care` | Testing and development |

## Configuration Examples

### Production Setup
```bash
# .env file for production
EKA_ENVIRONMENT=production
EKA_CLIENT_ID=prod_abc123
EKA_CLIENT_SECRET=prod_secret456
EKA_LOG_LEVEL=warn
EKA_MAX_RETRIES=5
```

### Development Setup
```bash
# .env file for development
EKA_ENVIRONMENT=development
EKA_CLIENT_ID=dev_abc123
EKA_CLIENT_SECRET=dev_secret456
EKA_LOG_LEVEL=debug
EKA_TIMEOUT=60
```

### CI/CD Setup
```bash
# Environment variables in CI/CD
export EKA_ENVIRONMENT=production
export EKA_CLIENT_ID="${VAULT_EKA_CLIENT_ID}"
export EKA_CLIENT_SECRET="${VAULT_EKA_CLIENT_SECRET}"
```

## Security Best Practices

✅ **Do:**
- Use environment variables for all credentials
- Use different credentials for different environments
- Store secrets in secure credential management systems
- Rotate credentials regularly
- Use `.env` files for local development (add to `.gitignore`)

❌ **Don't:**
- Hardcode credentials in source code
- Commit credentials to version control
- Use production credentials in development
- Share credentials via insecure channels

## Deployment Considerations

### Docker
```dockerfile
ENV EKA_CLIENT_ID=${EKA_CLIENT_ID}
ENV EKA_CLIENT_SECRET=${EKA_CLIENT_SECRET}
ENV EKA_ENVIRONMENT=production
```

### Kubernetes
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: eka-credentials
data:
  EKA_CLIENT_ID: <base64-encoded-client-id>
  EKA_CLIENT_SECRET: <base64-encoded-secret>
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        env:
        - name: EKA_ENVIRONMENT
          value: "production"
        envFrom:
        - secretRef:
            name: eka-credentials
```

### Cloud Platforms

**Google Cloud:**
- Use Secret Manager
- Set environment variables in Cloud Run/App Engine/GKE

**Azure:**
- Use Key Vault
- Set environment variables in App Service/Container Instances

## Troubleshooting

### Common Issues

**Missing Credentials:**
```
Error: client ID is required for authentication. Set EKA_CLIENT_ID environment variable or use WithClientID() option
```
Solution: Set the `EKA_CLIENT_ID` environment variable.

**Invalid Credentials:**
```
Error: failed to authenticate with provided credentials: invalid client credentials
```
Solution: Verify your credentials in the developer portal.

**Environment Detection:**
The SDK will default to `production` if `EKA_ENVIRONMENT` is not set. Always set this explicitly.
