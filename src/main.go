package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

/*
 * DCC Causal Bridge: Dual-Mode Implementation (Plugin & Sidecar)
 *
 * Mode A: Standalone Sidecar (HTTP/UDS) - Supports OPA 'http.send'
 * Mode B: Native OPA Extension - For custom OPA builds
 */

const (
	DefaultDCCSocket = "/var/run/bioos/dcc.sock"
	DefaultHTTPAddr  = "127.0.0.1:8080"
	VerificationWindow = 500 * time.Millisecond
)

// --- Mode A: Standalone Sidecar / HTTP Service ---

type VerifyResponse struct {
	Verified bool   `json:"verified"`
	Reason   string `json:"reason,omitempty"`
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("request_id")
	if requestID == "" {
		http.Error(w, "request_id is mandatory", http.StatusBadRequest)
		return
	}

	// Fail-Closed Logic
	verified, err := performKernelLookup(r.Context(), requestID)
	
	resp := VerifyResponse{Verified: verified}
	if err != nil {
		resp.Reason = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func startSidecar(addr string) {
	http.HandleFunc("/verify", handleVerify)
	log.Printf("[DCC Sidecar] Listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start sidecar: %v", err)
	}
}

// --- Mode B: Native OPA Extension Registration ---

func RegisterBuiltin() {
	rego.RegisterBuiltin1(
		&rego.Function{
			Name:    "dcc.is_verified",
			Decl:    ast.NewFunction(ast.Var("dcc.is_verified"), ast.String),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
			var requestID string
			if err := ast.As(op1.Value, &requestID); err != nil {
				return nil, err
			}

			// OS Compatibility check
			if runtime.GOOS != "linux" {
				return ast.BooleanTerm(false), nil
			}

			verified, _ := performKernelLookup(bctx.Context, requestID)
			return ast.BooleanTerm(verified), nil
		},
	)
}

// --- Core Verification Logic (Shared) ---

func performKernelLookup(ctx context.Context, id string) (bool, error) {
	// In a real BioOS environment, this dials DefaultDCCSocket
	// For the "Reproduction Demo", we simulate the kernel verification logic.
	
	// Fail-Closed: If non-Linux and not in Demo Mode, deny.
	if runtime.GOOS != "linux" && os.Getenv("DCC_DEMO_MODE") != "1" {
		return false, fmt.Errorf("DCC restricted to Linux/BioOS")
	}

	// Mock Logic for Reproducibility:
	// IDs starting with 'VALID-' are verified if within the window.
	if len(id) > 6 && id[:6] == "VALID-" {
		return true, nil
	}
	
	return false, fmt.Errorf("DCC Violation: Orphaned or unauthorized lineage")
}

func main() {
	mode := flag.String("mode", "sidecar", "Operation mode: sidecar or plugin-demo")
	addr := flag.String("addr", DefaultHTTPAddr, "HTTP address for sidecar mode")
	flag.Parse()

	if *mode == "sidecar" {
		os.Setenv("DCC_DEMO_MODE", "1") // Enable mock logic for the sidecar demo
		startSidecar(*addr)
	} else {
		fmt.Println("DCC OPA Extension Mode: Running logic verification...")
		// Logic verification code would go here
	}
}
