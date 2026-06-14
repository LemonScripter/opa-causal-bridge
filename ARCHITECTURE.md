# ARCHITECTURE.md: Causal Enforcement Specification

This document defines the architectural boundaries and the information flow of the DCC Causal Bridge for OPA.

## 1. System Data Flow

The verification process follows a deterministic, four-layer vertical integration:

1.  **Policy Layer (Rego):** OPA evaluates a policy containing `dcc.is_verified(input.request_id)`.
2.  **Extension Layer (Go):** The custom built-in intercepts the call and initiates a synchronous UDS dialer.
3.  **Service Layer (DCC Daemon):** A privileged daemon receives the request and performs an atomic lookup in the kernel map.
4.  **Enforcement Layer (eBPF/LSM):** The Linux kernel verifies the PID and intent ID against the `global_dcc_map`, populated by Tetragon LSM hooks.

## 2. Security Boundaries

| Layer | Component | Privilege | Responsibility |
| :--- | :--- | :--- | :--- |
| **User (L1)** | Rego Policy | Unprivileged | High-level logical access control |
| **Bridge (L2)** | Go Extension | OPA Runtime | Protocol fail-closed enforcement |
| **Control (L3)** | DCC Daemon | Root / CAP_BPF | Kernel IPC & Map management |
| **Trust (L4)** | eBPF Module | Ring 0 | Hardware-anchored causal truth |

## 3. Fail-Closed Implementation

The bridge implements a strict fail-closed posture at Layer 2. Any failure in the communication or verification path results in a `false` (DENY) response to the OPA engine.

- **DCC_OFFLINE:** Immediate denial if UDS is missing.
- **DCC_TIMEOUT:** 100ms hard limit on verification latency.
- **DCC_PROTOCOL_ERR:** Denial on malformed or partial status reads.

## 4. Logical Interaction Diagram

```text
[ Admission Request ] -> [ OPA Engine (Rego) ]
                                |
                      [ dcc.is_verified() ]    <-- (L2: Go Extension)
                                |
                    (Unix Domain Socket / IPC)
                                |
                        [ DCC Daemon ]         <-- (L3: Verification Logic)
                                |
                     (bpf_map_lookup_elem)
                                |
                     [ Linux Kernel (LSM) ]    <-- (L4: Hardware Anchor)
```
