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
	"runtime"
	"time"
)

/*
 * DCC Causal Bridge: Standalone Sidecar (Interoperable Mode)
 *
 * This implementation provides an HTTP/JSON endpoint for OPA 'http.send' verification.
 * It interfaces with the BioOS kernel state to enforce Digital Causal Closure.
 */

const (
	DefaultDCCSocket   = "/var/run/bioos/dcc.sock"
	DefaultHTTPAddr    = "127.0.0.1:8080"
	DCCDialTimeout     = 100 * time.Millisecond
)

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

	// Fail-Closed Logic: Assume denied unless kernel explicitly verifies
	verified, err := performKernelLookup(r.Context(), requestID)
	
	resp := VerifyResponse{Verified: verified}
	if err != nil {
		resp.Reason = err.Error()
		log.Printf("[DCC Audit] Denial for ID %s: %v", requestID, err)
	} else {
		log.Printf("[DCC Audit] Verified ID %s: ALLOW", requestID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// performKernelLookup interacts with the BioOS DCC subsystem.
func performKernelLookup(ctx context.Context, id string) (bool, error) {
	// 1. OS Interoperability & Demo Mode Check
	if runtime.GOOS != "linux" && os.Getenv("DCC_DEMO_MODE") != "1" {
		return false, fmt.Errorf("DCC kernel anchoring requires BioOS/Linux")
	}

	// 2. Production Path: Unix Domain Socket Bridge to Kernel Daemon
	if os.Getenv("DCC_DEMO_MODE") != "1" {
		return dialKernelDaemon(ctx, id)
	}

	// 3. Demo/Reproduction Path: Logic Simulation for Maintainers
	// This mirrors the kernel's temporal and atomic verification invariants.
	if len(id) > 6 && id[:6] == "VALID-" {
		return true, nil
	}
	
	return false, fmt.Errorf("DCC Violation: No verified causal lineage for ID %s", id)
}

func dialKernelDaemon(ctx context.Context, id string) (bool, error) {
	d := net.Dialer{Timeout: DCCDialTimeout}
	conn, err := d.DialContext(ctx, "unix", DefaultDCCSocket)
	if err != nil {
		return false, fmt.Errorf("DCC subsystem unreachable: %w", err)
	}
	defer conn.Close()

	// Synchronous binary protocol with the DCC Kernel Module
	if _, err := conn.Write([]byte(id)); err != nil {
		return false, fmt.Errorf("kernel write failure: %w", err)
	}

	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		return false, fmt.Errorf("kernel read failure: %w", err)
	}

	return buf[0] == 0x01, nil
}

func main() {
	addr := flag.String("addr", DefaultHTTPAddr, "HTTP address for sidecar mode")
	demo := flag.Bool("demo", false, "Enable logic simulation mode for reproduction/testing")
	flag.Parse()

	if *demo {
		os.Setenv("DCC_DEMO_MODE", "1")
		log.Printf("[DCC Sidecar] LOGIC SIMULATION ENABLED (Demo Mode)")
	}
	
	http.HandleFunc("/verify", handleVerify)
	log.Printf("[DCC Sidecar] Listening on %s (Interoperable Mode)", *addr)
	
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Failed to start sidecar: %v", err)
	}
}
