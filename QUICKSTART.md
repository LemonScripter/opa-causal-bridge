# Quickstart: Verifying Digital Causal Closure (DCC) in OPA

This guide provides step-by-step instructions to run the DCC Causal Bridge and verify the results yourself.

## Prerequisites
- **Go** (1.21 or later)
- **Linux** or **macOS** (for Unix Domain Socket support)
- **OPA CLI** (Optional, for manual evaluation)

## Step 1: Clone and Build the Proof
```bash
git clone https://github.com/LemonScripter/opa-causal-bridge.git
cd opa-causal-bridge
make build
```

## Step 2: Start the DCC Mock Service
In a separate terminal, start the service that simulates the BioOS kernel bridge:
```bash
go run tests/dcc_mock_service.go
```
*Note: The mock service listens on `/tmp/dcc_test.sock`. For this demo, we'll point the extension to this path.*

## Step 3: Run the Automated Verification
The included Python suite performs a comprehensive logic check (temporal integrity, anti-replay):
```bash
make test-integration
```

## Step 4: Manual Verification with OPA
You can test the logic closure using the OPA Go SDK or by running a custom OPA build that includes the `dcc` package.

### Logic Flow Verification:
1. **Verified Request:** Use an ID starting with `VALID-`.
   - Result: `dcc.is_verified("VALID-123")` -> **true**
2. **Orphaned Request:** Use any other ID.
   - Result: `dcc.is_verified("ghost-456")` -> **false** (Fail-Closed)
3. **Replay Protection:** Use the same `VALID-` ID twice.
   - Result: Second attempt -> **false** (Atomic Consumption)

## Why a Custom Built-in vs. `http.send`?
While OPA's `http.send` can communicate over Unix sockets, a dedicated `dcc.is_verified` built-in provides:
1. **Atomic Security Semantic:** Hardcoded fail-closed behavior and nanosecond-accurate temporal checks.
2. **Reduced Complexity:** Policy authors don't need to know the UDS binary protocol; they use a single, high-level primitive.
3. **Auditability:** Specific error types (`ErrServiceOffline`, `ErrTimeout`) allow for fine-grained security auditing at the OPA runtime level.

---
*MetaSpace.Bio Logic Project | [metaspace.bio](https://metaspace.bio)*
