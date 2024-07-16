# chariot-takehome


# Design Specification: Custom ID Structure

## 1. Technical Specification

### 1.1 ID Structure

The ID is a string of 20 characters with the following components:

1. **Prefix**: Variable length, ending with an underscore (`_`).
2. **Timestamp**: 12 characters, representing YYMMDDhhmmss.
3. **Random Part**: Remaining characters (variable length), alternating between letters and digits.

Example: `usr_230715123059A1B2C3D4E5`

### 1.2 Components

#### 1.2.1 Prefix
- Consists of lowercase letters.
- Ends with an underscore (`_`).
- Represents the type of entity (e.g., `usr` for user, `txn` for transaction).

#### 1.2.2 Timestamp
- 12 characters long.
- Format: YYMMDDhhmmss (Year, Month, Day, Hour, Minute, Second).
- Based on UTC time.

#### 1.2.3 Random Part
- Variable length (fills the remainder of the 20-character ID).
- Alternates between uppercase letters and digits.
- Uses nanosecond precision for added entropy.

### 1.3 Generation Process

1. Start with the predefined prefix.
2. Append the current timestamp in YYMMDDhhmmss format.
3. Generate the random part:
   - Use nanosecond precision for entropy.
   - Alternate between letters and digits.
   - Fill the remaining length of the ID.

### 1.4 Validation Process

1. Check the total length (must be 20 characters).
2. Identify the prefix by locating the underscore.
3. Extract and validate the 12-character timestamp.
4. Verify the alternating pattern of letters and digits in the random part.

## 2. Rationale for Design Decisions

### 2.1 Overall Structure

The ID structure is inspired by Stripe's object IDs and Segment's KSUIDs, combining readability, sortability, and uniqueness.

**Rationale:**
- Human-readable: The prefix and timestamp make the ID meaningful at a glance.
- Sortable: The timestamp allows chronological sorting.
- Unique: The combination of timestamp and random part ensures uniqueness.

### 2.2 Prefix

**Decision:** Include a variable-length prefix ending with an underscore.

**Rationale:**
- Improves readability and debugging by immediately identifying the entity type.
- Flexible length allows for various entity types without changing the overall ID structure.
- The underscore clearly separates the prefix from the rest of the ID.

### 2.3 Timestamp

**Decision:** Use a 12-character timestamp in YYMMDDhhmmss format.

**Rationale:**
- Provides chronological sorting capability.
- Second-level precision is sufficient for most applications.
- 12 characters balance between precision and ID length.
- Using UTC avoids timezone complications.

### 2.4 Random Part

**Decision:** Use alternating uppercase letters and digits, influenced by nanosecond precision.

**Rationale:**
- Alternating pattern increases the character set without using special characters, improving readability.
- Uppercase letters are distinct from the lowercase prefix, enhancing visual parsing

### 4. Potential Improvements and Considerations
- Base Encoding: Explore base32 or base64 encoding to increase information density.
- Collision Handling: Implement retry logic in case of extremely rare collisions.


## Setup Instructions

1. Install go-task:
   ```
   brew install go-task
   ```

2. Set up the project:
   ```
   task setup
   ```

3. Generate API proto files:
   ```
   task proto:gen:api
   ```

4. Build and start Docker containers:
   ```
   docker compose up -d --no-deps --build
   ```

5. Run database migrations:
   ```
   task migrate:up:local
   ```

## Cleanup Instructions

To stop and remove Docker containers:
```
docker compose down
```

## Testing the API

You can use the following curl commands to test the API endpoints:

```bash
# Base URL - replace with your actual API URL
BASE_URL="http://localhost:8080"

# Create a user
curl -X POST -H "Content-Type: application/json" -d '{"name": "John Doe", "email": "john@example.com"}' "$BASE_URL/create_user"

# Create an account
curl -X POST -H "Content-Type: application/json" -d '{"user_id": "usr_[your-returned-id]", "account_type": "checking", "account_state":"open"}' "$BASE_URL/create_account"

# Deposit funds only allows to deposit from external ledger account to internal ledger account
curl -X POST -H "Content-Type: application/json" -d '{"debit_account_id": "acct_[your-ext-account-id]", "credit_account_id": "acct_[your-int-account-id]", "amount": 1000,  "idempotency_key": "blah"}' "$BASE_URL/deposit_funds"

# Withdraw funds only allows to deposit from internal ledger account to external ledger account
curl -X POST -H "Content-Type: application/json" -d '{"debit_account_id": "acct_[your-int-account-id]", "credit_account_id": "acct_[your-ext-account-id]", "amount": 500,  "idempotency_key": "blah"}' "$BASE_URL/withdraw_funds"

# Transfer funds between internal accounts only 
curl -X POST -H "Content-Type: application/json" -d '{"debit_account_id": "acct_[your-int-account-id]", "credit_account_id": "acct_[another-persons-int-account-id]", "amount": 250,  "idempotency_key": "tr_123"}' "$BASE_URL/transfer_funds"

# List transactions
curl -X GET "$BASE_URL/list_transactions?account_id=acct_[your-account-id]"

# Get account balance
curl -X GET "$BASE_URL/get_account_balance?account_id=acct_[your-acct-id]"
```

Note: Replace `usr_[your-user-id]` and `acct_[your-account-id]` with actual IDs from your system.

## Idempotency Implementation

Idempotency is achieved using a simple in-memory storage of idempotency keys. Here's a brief explanation:

```go
var (
    idempotencyKeys = make(map[string]bool)
    mu              sync.Mutex
)

func CreateUserHandler(grpcClient *client.ApiClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... [request parsing code] ...

        // Check for idempotency key
        idempotencyKey := r.Header.Get("Idempotency-Key")
        if idempotencyKey != "" {
            mu.Lock()
            if _, exists := idempotencyKeys[idempotencyKey]; exists {
                mu.Unlock()
                http.Error(w, "Duplicate request", http.StatusConflict)
                return
            }
            idempotencyKeys[idempotencyKey] = true
            mu.Unlock()
        }

        // ... [rest of the handler code] ...
    }
}
```

This implementation uses a map to store idempotency keys. When a request comes in with an idempotency key:
1. We check if the key exists in our map.
2. If it does, we return a conflict error, preventing duplicate operations.
3. If it doesn't, we add the key to the map and proceed with the operation.

Note: This simple implementation is not suitable for production use as it doesn't handle key expiration or persistence across server restarts. For production, consider using a distributed cache or database to store idempotency keys.

## Concurrency Handling

Concurrency is managed using database transactions with serializable isolation level:

```go
tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
if err != nil {
    return "", err
}
```

This approach ensures that concurrent operations on the same data are executed sequentially, preventing race conditions and maintaining data consistency. The serializable isolation level provides the highest level of isolation, ensuring that concurrent transactions appear to execute sequentially.

Key points:
1. Serializable isolation prevents read and write skew anomalies.
2. If a concurrent transaction conflicts, one will be rolled back and can be retried.
3. This approach may impact performance under high concurrency, so proper monitoring and optimization may be necessary.
 
I did not include benchmark tests for the Identifier library. 
The only method that is untested is the ListTransactions. 


## Miscellaneous
### Logging
This approach to logging closely follows BetterStack's recommendations for structured logging, which greatly enhances the ability to monitor and debug the application efficiently.

Structured logging provides a consistent format for log messages, making it easier to parse and analyze logs, especially in cloud environments. This is particularly useful for handling complex operations and concurrency in financial applications like the ones you're involved with.

Example Implementation:

### Setup the Logger:

```go
h := &logger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, opts)}
slogLogger := slog.New(h)
slog.SetDefault(slogLogger)
err := env.Parse(&c)
if err != nil {
    slog.Error("failed to parse default config", "error", err)
    os.Exit(1)
}
slog.Info("starting grpc", "port", c.ServerPort)
```

### Context Handler:

```go
type ctxKey string

const (
    slogFields ctxKey = "slog_fields"
)

type ContextHandler struct {
    slog.Handler
}

func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // Get the metadata from the incoming context
        md, ok := metadata.FromIncomingContext(ctx)
        if (!ok) {
            return nil, fmt.Errorf("couldn't parse incoming context metadata")
        }

        for k, v := range md {
            if len(v) > 1 {
                ctx = AppendCtx(ctx, slog.Any(k, v))
            } else {
                ctx = AppendCtx(ctx, slog.String(k, v[0]))
            }
        }
        slog.InfoContext(ctx, "gRPC request")
        return handler(ctx, req)
    }
}

// Handle adds contextual attributes to the Record before calling the underlying handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
    if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
        for _, v := range attrs {
            r.AddAttrs(v)
        }
    }

    return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr ...slog.Attr) context.Context {
    if parent == nil {
        parent = context.Background()
    }

    var newAttrs []slog.Attr
    if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
        newAttrs = append(v, attr...)
    } else {
        newAttrs = append([]slog.Attr{}, attr...)
    }

    return context.WithValue(parent, slogFields, newAttrs)
}
```
## Benefits:


- Improved Filtering: Structured logs make it easier to filter specific log entries based on attributes, improving troubleshooting and monitoring.
- Enhanced Readability: Logs are formatted in a consistent manner, making them easier to read and analyze.
- Seamless Integration with AWS CloudWatch: Structured logs can be easily filtered using CloudWatch's filter and pattern syntax, enabling more efficient log management.

For more information on AWS CloudWatch log filtering, refer to the following resources:

- [AWS CloudWatch Filter and Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html)
- [The Importance of Structured Logging in AWS](https://medium.com/@connorbutch/the-importance-of-structured-logging-in-aws-and-anywhere-else-52a4534c53aa)
- [BetterStack guide to slog](https://betterstack.com/community/guides/logging/logging-in-go/)



# TODO for later
- update readme
- complete readme: draw out image of architecture
- fix benchmark testing for the identifier library
- add comment on the necessity for double-entry accounting
- add an example walkthrough of how to debug logging in cloudwatch
- improve cursor to be the next row not last row
- extra: add middleware to gateway to check for authentication
- extra: add rate limiter
- extra: add health monitor 
- extra: add a pr branch for how to do an ach transaction with a 3p like moov  
### Other resources 
- https://dev.to/stripe/designing-apis-for-humans-object-ids-3o5a
- https://medium.com/@RobertKhou/double-entry-accounting-in-a-relational-database-2b7838a5d7f8 
- https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html
- https://clerk.com/blog/generating-sortable-stripe-like-ids-with-segment-ksuids
- https://segment.com/blog/a-brief-history-of-the-uuid/
- https://github.com/paralleldrive/cuid2
