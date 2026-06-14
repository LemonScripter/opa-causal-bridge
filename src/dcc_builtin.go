package dcc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

/*
 * DCC OPA Built-in: Hardened Implementation
 */

const (
	DCCSocketPath  = "/var/run/bioos/dcc.sock"
	DCCDialTimeout = 100 * time.Millisecond
	DCCReadTimeout = 50 * time.Millisecond
)

// Custom error types for fine-grained logging and security auditing
var (
	ErrServiceOffline = errors.New("DCC service unreachable")
	ErrProtocolError  = errors.New("DCC protocol violation")
	ErrTimeout        = errors.New("DCC verification timeout")
	ErrUnsupportedOS  = errors.New("DCC kernel anchoring is only supported on Linux")
)

func Register() {
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

			// OS Interoperability: DCC is currently kernel-anchored on Linux.
			// On other OSs, we enforce a Fail-Closed posture.
			if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
				// Note: Darwin added for potential future expansion, but currently unsupported.
				return ast.BooleanTerm(false), nil
			}

			// Fail-Closed Logic: Any error in verification returns false
			verified, err := verifyCausalState(bctx.Context, requestID)
			if err != nil {
				// Audit log entry would go here in a production BioOS environment
				return ast.BooleanTerm(false), nil
			}

			return ast.BooleanTerm(verified), nil
		},
	)
}

func verifyCausalState(ctx context.Context, id string) (bool, error) {
	// Re-verify OS support at the dialer level
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return false, ErrUnsupportedOS
	}

	d := net.Dialer{Timeout: DCCDialTimeout}
	
	conn, err := d.DialContext(ctx, "unix", DCCSocketPath)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrServiceOffline, err)
	}
	defer conn.Close()

	// Set deadline for the entire operation to prevent OPA hanging
	deadline := time.Now().Add(DCCReadTimeout)
	conn.SetDeadline(deadline)

	// Protocol: RequestID (max 256 bytes)
	_, err = conn.Write([]byte(id))
	if err != nil {
		return false, fmt.Errorf("%w: failed to write request", ErrProtocolError)
	}

	// Expect 1-byte response: 0x01 (Verified) or 0x00 (Denied)
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false, ErrTimeout
		}
		return false, fmt.Errorf("%w: failed to read response", ErrProtocolError)
	}

	if n != 1 {
		return false, fmt.Errorf("%w: invalid response length", ErrProtocolError)
	}

	return buf[0] == 0x01, nil
}
