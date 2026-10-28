> **This package has moved.** `github.com/bonzo-signalhouse/SignalHouseSDK-GO` is now [`github.com/signalhouse-sms/sdk-go`](https://pkg.go.dev/github.com/signalhouse-sms/sdk-go), developed at
> https://github.com/signalhouse-sms/sdk-go. This package is a thin shim that installs `github.com/signalhouse-sms/sdk-go` at the
> same version and re-exports it. It receives every release until **2027-05-01**, then stops.
>
> **Migrate:** `go get github.com/signalhouse-sms/sdk-go`, then change imports from `github.com/bonzo-signalhouse/SignalHouseSDK-GO/signalhouse` to `github.com/signalhouse-sms/sdk-go/signalhouse`. The API is identical.

# SignalHouse Go SDK

Official Go SDK for the SignalHouse API. Provides a type-safe client for managing SMS/MMS messaging, phone numbers, 10DLC brands/campaigns, billing, and more.

Full API reference and guides: https://app2.signalhouse.io/docs

## Requirements

- Go 1.21+

## Installation

```bash
go get github.com/signalhouse-sms/sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/signalhouse-sms/sdk-go/signalhouse"
)

func main() {
    client, err := signalhouse.NewClient(signalhouse.Config{
        APIKey:  "your-api-key",
        BaseURL: "https://v2.signalhouse.io",
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Login
    resp, err := client.Auth.Login(ctx, signalhouse.LoginParams{
        Email:    "user@example.com",
        Password: "secret",
    }, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Success: %v\n", resp.Success)
}
```

## Authentication

The SDK supports two authentication modes:

1. **API Key** (set at client creation): Used by default for all requests.
2. **Per-request token** (via `RequestOptions`): Overrides the API key for a single request.

```go
// Use a per-request token (e.g., JWT from login)
opts := &signalhouse.RequestOptions{Token: "user-jwt-token"}
resp, err := client.Numbers.GetPhoneNumbers(ctx, signalhouse.GetPhoneNumbersParams{
    GroupID: "G12345678",
}, opts)
```

Set `RequestOptions.IdempotencyKey` to make a create, update or delete safe to retry: the SDK sends it as the `Idempotency-Key` header, and a retry with the same key and request replays the first response instead of running it again. Keys are scoped to the authenticated user and kept for 24 hours.

## API Domains

All API methods are organized by domain:

| Service | Access | Description |
|---------|--------|-------------|
| `client.Auth` | `AuthService` | Login, password reset, auth history |
| `client.Billing` | `BillingService` | Wallets, payments, transactions, invoices, fees |
| `client.Brands` | `BrandsService` | 10DLC brand registration and management, toll-free business verification |
| `client.Campaigns` | `CampaignsService` | 10DLC campaign management |
| `client.Groups` | `GroupsService` | Group management |
| `client.Landings` | `LandingsService` | Landing page management (with file upload) |
| `client.Messages` | `MessagesService` | SMS, MMS, group messaging, analytics |
| `client.Templates` | `TemplatesService` | Versioned message templates (`/template`): drafts, publish, validate, media |
| `client.Notifications` | `NotificationsService` | Notification management |
| `client.Numbers` | `NumbersService` | Phone number search, purchase, assignment |
| `client.Shortlinks` | `ShortlinksService` | URL shortener |
| `client.Subgroups` | `SubgroupsService` | Subgroup management |
| `client.Subscriptions` | `SubscriptionsService` | Subscription management |
| `client.Users` | `UsersService` | User management |
| `client.Webhooks` | `WebhooksService` | Webhook management |

## Admin Methods

Some services have admin-only methods accessible via the `.Admin` field. Enable them at client creation:

```go
client, _ := signalhouse.NewClient(signalhouse.Config{
    APIKey:      "admin-api-key",
    BaseURL:     "https://v2.signalhouse.io",
    EnableAdmin: true,
})

// Admin-only: approve a campaign
resp, err := client.Campaigns.Admin.ApproveCampaign(ctx, signalhouse.ApproveCampaignParams{
    CampaignID: "campaign-id",
}, opts)
```

Admin services are available on: `Campaigns`, `Groups`, `Subscriptions`, and `Users`.

## Response Format

All methods return a `*signalhouse.Response`:

```go
type Response struct {
    Success bool            // true for 2xx responses
    Data    json.RawMessage // response body (unmarshal as needed)
    Status  int             // HTTP status code
    Error   string          // error message for non-2xx responses
}
```

## File Uploads

Landing pages and MMS messages support file uploads via `FileUpload`:

```go
file, _ := os.Open("logo.png")
defer file.Close()

resp, err := client.Landings.CreateLanding(ctx, signalhouse.CreateLandingParams{
    LandingData: signalhouse.CreateLandingData{
        BrandID:     "brand-id",
        Description: "My landing page",
    },
    File: &signalhouse.FileUpload{
        Filename: "logo.png",
        Reader:   file,
    },
}, opts)
```

## Error Handling

The SDK returns two types of errors:

- `*signalhouse.ValidationError`: Missing required parameter (client-side validation)
- `*signalhouse.APIError`: Non-2xx API response

```go
resp, err := client.Billing.GetWallet(ctx, signalhouse.GetWalletParams{}, opts)
if err != nil {
    var valErr *signalhouse.ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Missing field: %s\n", valErr.Field)
    }
}
```

## Custom HTTP Client

```go
client, _ := signalhouse.NewClient(signalhouse.Config{
    APIKey:     "your-api-key",
    BaseURL:    "https://v2.signalhouse.io",
    HTTPClient: &http.Client{Timeout: 30 * time.Second},
})
```

## Examples

See the [examples/basic](examples/basic/) directory for a complete usage example.
