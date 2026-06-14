# DCC Causal Enforcement for Open Policy Agent (OPA)

## Technical Specification: Causal Built-in

This repository provides a concrete implementation of **Digital Causal Closure (DCC)** for the Open Policy Agent. It extends the Rego policy engine with a hardware-anchored causal verification primitive.

### Architectural Overview

Traditional identity-based policies are vulnerable to autonomous session hijacking. The DCC bridge addresses this by enforcing **Runtime Causal Lineage**.

1.  **Enforcement Point:** OPA Rego policies call `dcc.is_verified(request_id)`.
2.  **Verification Path:** The Go extension performs a synchronous, low-latency query via a Unix Domain Socket to the BioOS DCC service.
3.  **Kernel Anchor:** The DCC service interacts with the `global_dcc_map` (eBPF HASH map), which is populated by kernel-level IRQ/LSM hooks (`dcc_bridge.bpf.c`).
4.  **Atomic Verification:** The kernel verifies that the process initiating the request is currently within a valid **Causality Window** authorized by a physical event.

### Implementation Details

*   **`src/dcc_builtin.go`**: OPA Go extension using `rego.RegisterBuiltin1`. Implements fail-closed logic for unix socket communication.
*   **Protocol:** Synchronous binary protocol over `/var/run/bioos/dcc.sock`.
*   **Complexity:** O(1) lookup in eBPF kernel space.

### Scientific Reference
Formal proof for the Causal OS paradigm: [DOI: 10.5281/zenodo.20384700](https://doi.org/10.5281/zenodo.20384700).

---
*Technical Implementation by MetaSpace BioOS Team*
