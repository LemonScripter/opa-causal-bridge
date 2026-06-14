# DCC Integration Roadmap: From Research to Production

This roadmap outlines the path from the current functional bridge to a fully hardware-anchored Digital Causal Closure (DCC) enforcement system.

## Phase 1: Functional Bridge & Simulation (Current)
- [x] OPA custom built-in registration (`dcc.is_verified`)
- [x] Fail-closed logic implementation
- [x] Unix Domain Socket protocol definition
- [x] Mock service for logic verification
- [x] Python-based causal window simulation

## Phase 2: Kernel-Space Integration (Q3 2026)
- [ ] **eBPF Map Synchronization:** Implement the background sync between Tetragon's `TracingPolicy` and the `global_dcc_map`.
- [ ] **LSM Hook Hardening:** Move from `kprobes` to `LSM` (Linux Security Modules) hooks for `socket_connect` and `bpffs_map` access.
- [ ] **Atomic Token Consumption:** Ensure thread-safe, non-blocking token consumption in kernel space using `bpf_spin_lock`.

## Phase 3: Hardware Anchoring (Q4 2026)
- [ ] **IRQ Correlation:** Bind causal tokens to physical hardware interrupts (Keyboard, Network IRQ) to ensure human-in-the-loop/verified event provenance.
- [ ] **TPM/TEE Integration:** Signed causal tokens anchored in Trusted Execution Environments.

## Phase 4: Ecosystem Scale (2027)
- [ ] **Multi-Tool Bridges:** Standardize the DCC protocol for Falco, Istio, and SPIRE.
- [ ] **BioOS Native Support:** Deep integration into the BioOS kernel distribution.

---
*Maintained by MetaSpace BioOS Team*
