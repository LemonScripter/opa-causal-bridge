# DCC Causal Enforcement for Open Policy Agent (OPA)

## Overview
The **DCC Causal Enforcement** module is a professional extension for the Open Policy Agent (OPA). It introduces **Digital Causal Closure (DCC)** to the Rego policy engine, allowing for hardware-anchored, runtime-aware decision making.

## The Problem: Contextless Mutations
Standard OPA policies (Rego) are stateless and context-blind regarding the *causal origin* of an API request. A policy can verify "What" is being requested, but cannot verify "Why" (the intent chain). This leads to a security gap where authorized identities can be used to perform unauthorized, autonomous mutations within a cluster.

## The Solution: Hardware-Anchored Causal Context
This module extends OPA with a custom **Rego Built-in Function** (`dcc.is_verified`) that queries the BioOS DCC kernel module.
- **Causal Verification:** Rego policies can now require a valid, non-expired Causal Token for any sensitive mutation (e.g., Kubernetes Admission Control).
- **Runtime Awareness:** Decisions are no longer based solely on static manifest data, but on real-time evidence of a verified causal chain.

## Scientific Background
This integration is based on the following formal research:
- [The Causal Operating System: Digital Causal Closure for Autonomous Systems](https://doi.org/10.5281/zenodo.20384700)
- [BioOS Causal Constitution (PDF)](https://bioos.metaspace.bio/bioos_causal_constitution_en.pdf)

## Components
- **`dcc_builtin.go`**: OPA Go extension implementing the `dcc.is_verified` Rego function.
- **`verify_opa.py`**: Logic verification suite ensuring 100% enforcement accuracy for admission requests.

## Upstreaming Proposal
We propose the inclusion of Causal Context as a standard extension for OPA deployments in high-security environments, linking policy evaluation to the physical reality of user intent.

---
*Created by MetaSpace BioOS | [metaspace.bio](https://metaspace.bio) | [admin@metaspace.bio](mailto:admin@metaspace.bio)*
