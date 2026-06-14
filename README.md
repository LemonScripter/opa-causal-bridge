# DCC Causal Bridge for Open Policy Agent (OPA)

[![Status](https://img.shields.io/badge/Status-Research--Prototype-blue)](ROADMAP.md)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.20384700.svg)](https://doi.org/10.5281/zenodo.20384700)

## Overview

The **DCC Causal Bridge** is a research-grade extension for the Open Policy Agent (OPA). It introduces **Digital Causal Closure (DCC)**—a security paradigm that binds Rego policy decisions to hardware-anchored causal events. By bridging the semantic gap between unprivileged policy evaluation and kernel-level execution lineage, it physically eliminates unauthorized autonomous mutations.

### Key Value Propositions

- **Fail-Closed Security Architecture:** Deterministic denial of all requests if the DCC verification service is unreachable or if a protocol violation occurs.
- **Causal Lineage Enforcement:** Every mutation must be backed by a verified, non-expired Causal Token anchored in kernel-space eBPF/LSM hooks.
- **Low-Latency Integration:** Implemented as a native Go built-in utilizing high-performance Unix Domain Socket (UDS) IPC for O(1) lookups.
- **Empirical Validation:** Verified on a live research node in Tokyo (GCP `asia-northeast1`) against replay attacks and stale intent hijacking.

## Technical Components

- **Built-in:** `dcc.is_verified(request_id)` registered via OPA `rego.RegisterBuiltin1`.
- **Enforcement:** Synchronous binary protocol over `/var/run/bioos/dcc.sock`.
- **Verification:** Logic-verified against the [BioOS Causal Constitution](https://doi.org/10.5281/zenodo.20384700).

## Empirical Proof

Detailed execution logs and system environment details from our Tokyo Node can be found in [VERIFICATION.md](VERIFICATION.md).

---
*Developed by the MetaSpace BioOS Team | Security Engineering & Research*
