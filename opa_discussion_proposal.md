# Proposal: Hardware-Anchored Causal Context for OPA Rego Policies

Hello OPA/Styra Community,

Open Policy Agent (OPA) is the gold standard for decoupled policy enforcement. However, as workloads become more autonomous, we see a growing need for policies to verify the **Causal Origin** of a request—the evidence of intent that justifies a mutation.

### The Problem: Stateless Admission Control
Currently, OPA makes decisions based on the provided `input` JSON. It can verify "What" is being requested, but it cannot verify "Why." An authorized Service Account could be used by a compromised agent to send a valid, but unauthorized mutation request. OPA has no way to distinguish this from a legitimate, human-initiated action.

### Proposal: DCC Causal Built-in for Rego
I propose introducing a **DCC (Digital Causal Closure)** extension for OPA. This extension adds a new Rego built-in function, `dcc.is_verified(request_id)`, which allows policies to check if a request is backed by a hardware-anchored, kernel-level causal chain.

### Proof of Concept
I have developed a PoC showing how a DCC built-in can be integrated into the OPA engine to enforce causality-aware admission control:
[https://github.com/LemonScripter/opa-causal-bridge](https://github.com/LemonScripter/opa-causal-bridge)

### Key Benefits:
1.  **Context-Aware Decisions:** Policies can now require proof of intent (Causal Token) for sensitive actions like `Deployment` or `Secret` mutations.
2.  **Runtime Integrity:** Links high-level policy evaluation to the physical reality of the kernel's causal state.
3.  **Peer-Reviewed Foundation:** Based on research in Causal Operating Systems (DOI: 10.5281/zenodo.20384700).

We believe this paradigm can significantly harden Kubernetes security and autonomous agent orchestration. We would love to discuss the feasibility of exposing such causal metadata as a first-class citizen in the Rego engine.

Best regards,

**MetaSpace BioOS Team**
[metaspace.bio](https://metaspace.bio) | [admin@metaspace.bio](mailto:admin@metaspace.bio)
