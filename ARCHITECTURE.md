# ARCHITECTURE.md: Interoperable Causal Enforcement

The DCC OPA Bridge supports a modular architecture designed for high interoperability and "Fail-Closed" security.

## Deployment Models

The system supports two primary integration paradigms to bridge high-level policies with kernel-level causal truth.

### 1. Standalone Sidecar Model (Loosely Coupled)
In this model, the DCC service runs as a sidecar process. OPA interacts with it using standard HTTP/JSON protocols.

**Flow:**
`OPA (http.send) -> DCC Sidecar (HTTP) -> DCC Daemon (UDS) -> eBPF Map -> Kernel`

- **Pros:** Maximum reach, cross-OS support for the policy engine, no custom OPA build required.
- **Cons:** Slight network overhead compared to native plugins.

### 2. Native Plugin Model (Tightly Coupled)
The DCC verification logic is compiled directly into OPA as a custom Go built-in.

**Flow:**
`OPA (dcc.is_verified) -> Unix Domain Socket -> DCC Daemon -> eBPF Map -> Kernel`

- **Pros:** Zero-overhead, atomic security semantic, nanosecond-accurate checks.
- **Cons:** Requires a custom OPA build.

## Security Boundaries

| Layer | Responsibility | Posture |
| :--- | :--- | :--- |
| **Policy (L1)** | Logic evaluation | Reject by default |
| **Sidecar (L2)** | Protocol translation | Fail-Closed on error |
| **Kernel (L3)** | Hardware-anchored truth | Atomic consumption |

## Data Flow Diagram

```text
[ External Request ]
        |
[ OPA Policy Engine ] ---------------------+
        |                                  |
   (Integration Choice)                    |
        |                                  |
        +---> [ DCC Sidecar (HTTP) ] ------+
        |                                  |
        +---> [ DCC Plugin (Native) ] -----+
                                           |
                                [ Unix Domain Socket ]
                                           |
                                   [ DCC Kernel Module ]
                                           |
                                    [ eBPF LSM Hooks ]
```
