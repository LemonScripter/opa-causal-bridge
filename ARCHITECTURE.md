# DCC Causal Bridge: Architectural Specification

This document describes the flow of information and the security boundaries within the Digital Causal Closure (DCC) integration for Open Policy Agent.

## Operational Flow

1.  **Policy Evaluation (Rego):** A request arrives at OPA. The policy includes a call to `dcc.is_verified(input.request_id)`.
2.  **Go Built-in Execution:** The custom OPA extension intercepts the call. It extracts the `request_id` and initiates a synchronous query.
3.  **Unix Socket Bridge:** The extension connects to `/var/run/bioos/dcc.sock`. This provides a high-performance, local-only communication path.
4.  **DCC Service (Daemon):** A privileged daemon receives the request. It acts as the gatekeeper between userspace OPA and the Linux Kernel.
5.  **Kernel Verification (eBPF/LSM):** The daemon performs a `bpf_map_lookup_elem` on the `global_dcc_map`. 
    *   The map is populated by **Tetragon LSM hooks** triggered by physical IRQs (Human/Hardware intent).
    *   The kernel verifies the PID, Intent ID, and the 500ms Causality Window.
6.  **Response Path:** The result (Verified/Denied) propagates back through the socket to OPA, where the Rego policy completes its decision.

## Security Boundaries

| Component | Privilege | Responsibility |
| :--- | :--- | :--- |
| **Rego Policy** | Unprivileged | High-level logical decision |
| **Go Extension** | OPA Runtime | Protocol bridge & Fail-closed enforcement |
| **DCC Service** | Root / CAP_BPF | Kernel map interaction |
| **eBPF Module** | Ring 0 | Hardware-anchored truth & Atomic closure |

## Data Flow Diagram

```text
[ User Request ] -> [ OPA Engine (Rego) ]
                          |
                  [ dcc.is_verified() ]  <-- (Go Built-in)
                          |
               (Unix Domain Socket)
                          |
                  [ DCC Daemon ]         <-- (Verification Logic)
                          |
               (eBPF Map Lookup)
                          |
               [ Linux Kernel (LSM) ]    <-- (Hardware Anchor)
```
