# Quickstart: DCC Causal Enforcement for OPA

This guide provides a self-contained, 3-step environment to verify the **Digital Causal Closure (DCC)** integration for OPA.

## Prerequisites
- **Go** (1.21+)
- **Curl** (for API verification)
- **Linux** or **macOS**

## Step 1: Build the DCC Environment
Clone the repository and build the binary:
```bash
git clone https://github.com/LemonScripter/opa-causal-bridge.git
cd opa-causal-bridge
make build
```

## Step 2: Spin up the DCC Sidecar
Start the standalone verification service. This service simulates the BioOS kernel bridge:
```bash
./bin/dcc-sidecar --mode sidecar --addr 127.0.0.1:8080
```

## Step 3: Verify the Logic
In a separate terminal, run the automated reproduction script to see the "Fail-Closed" behavior in action:
```bash
./reproduce.sh
```

### Integration Paths

You can integrate DCC with OPA using two different paradigms:

1.  **Standalone Sidecar (Recommended for Interoperability):**
    OPA calls the DCC service via the native `http.send` built-in.
    ```rego
    allow {
        resp := http.send({"method": "get", "url": "http://localhost:8080/verify?request_id=" + input.request_id})
        resp.body.verified == true
    }
    ```
2.  **Native Go Plugin:**
    For custom OPA builds, use the high-performance `dcc.is_verified(id)` built-in registered in `src/main.go`.

## Verification Scenarios
- **Authorized:** Request IDs starting with `VALID-` simulate a hardware-verified intent.
- **Fail-Closed:** Any network error, timeout, or missing token results in a `verified: false` status, ensuring maximum security.

---
*Production-Grade Research Prototype by MetaSpace BioOS Team*
