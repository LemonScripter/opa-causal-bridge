# DCC Causal Bridge for Open Policy Agent (OPA)

[![Verified](https://img.shields.io/badge/Verified-Tokyo--Node-green)](VERIFICATION.md)
[![Status](https://img.shields.io/badge/Status-Research--Prototype-blue)](ROADMAP.md)
[![Project](https://img.shields.io/badge/BioOS-Causal--Security-green)](https://metaspace.bio)
[![DOI](https://img.shields.io/badge/DOI-10.5281%2Fzenodo.20384700-purple)](https://doi.org/10.5281/zenodo.20384700)

## Overview

The **DCC Causal Bridge** is a research-grade extension for the Open Policy Agent (OPA). It introduces **Digital Causal Closure (DCC)**—a security paradigm that binds Rego policy decisions to hardware-anchored causal events. By bridging the semantic gap between unprivileged policy evaluation and kernel-level execution lineage, it physically eliminates unauthorized autonomous mutations.

### 🚀 [Run the Proof: Quickstart Guide](QUICKSTART.md)

### Key Value Propositions

- **Fail-Closed Security Architecture:** Deterministic denial of all requests if the DCC verification service is unreachable or if a protocol violation occurs.
- **Causal Lineage Enforcement:** Every mutation must be backed by a verified, non-expired Causal Token anchored in kernel-space eBPF/LSM hooks.
- **Low-Latency Integration:** Implemented as a native Go built-in utilizing high-performance Unix Domain Socket (UDS) IPC for O(1) lookups.
- **Empirical Validation:** Verified on a live research node in Tokyo (GCP `asia-northeast1`) against replay attacks and stale intent hijacking.

## Technical Components

- **Built-in:** `dcc.is_verified(request_id)` registered via OPA `rego.RegisterBuiltin1`.
- **Enforcement:** Synchronous binary protocol over `/var/run/bioos/dcc.sock`.
- **Verification:** Logic-verified against the BioOS DCC standards.

### Scientific & Technical Foundation

This implementation is based on the following formal specifications and research:

- **Research Paper:** [The Causal Operating System: Digital Causal Closure for Autonomous Systems](https://doi.org/10.5281/zenodo.20384700) (DOI: 10.5281/zenodo.20384700)
- **Formal Specification:** [BioOS Causal Constitution (PDF)](https://bioos.metaspace.bio/bioos_causal_constitution_en.pdf)

## Empirical Proof

Detailed execution logs and system environment details from our Tokyo Node can be found in [VERIFICATION.md](VERIFICATION.md).

---
*MetaSpace.Bio Logic Project | [metaspace.bio](https://metaspace.bio) | [admin@metaspace.bio](mailto:admin@metaspace.bio)*
