package dcc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

/*
 * OPA DCC Causal Built-in (Go)
 *
 * This extension adds the 'dcc.is_verified' function to the Rego language.
 * It queries the BioOS DCC verification service (local unix socket) to confirm
 * if a Request ID is backed by a verified hardware-anchored causal chain.
 */

const (
	DCCSocketPath = "/var/run/bioos/dcc.sock"
	DCCDialTimeout = 100 * time.Millisecond
)

func Register() {
	rego.RegisterBuiltin1(
		&rego.Function{
			Name:             "dcc.is_verified",
			Decl:             ast.NewFunction(ast.Var("dcc.is_verified"), ast.String),
			Memoize:          true,
		},
		func(bctx rego.BuiltinContext, op1 *ast.Term) (*ast.Term, error) {
			var requestID string
			if err := ast.As(op1.Value, &requestID); err != nil {
				return nil, err
			}

			// Perform real causal verification via DCC service
			verified, err := verifyCausalState(bctx.Context, requestID)
			if err != nil {
				// We log the error but return false to enforce "Fail-Closed" security
				return ast.BooleanTerm(false), nil
			}

			return ast.BooleanTerm(verified), nil
		},
	)
}

// verifyCausalState connects to the DCC daemon and verifies the token for the given request.
func verifyCausalState(ctx context.Context, id string) (bool, error) {
	d := net.Dialer{Timeout: DCCDialTimeout}
	conn, err := d.DialContext(ctx, "unix", DCCSocketPath)
	if err != nil {
		return false, fmt.Errorf("failed to connect to DCC service: %w", err)
	}
	defer conn.Close()

	// Protocol: Write Request ID, Read 1-byte status (0x01 = Verified, 0x00 = Failed)
	_, err = conn.Write([]byte(id))
	if err != nil {
		return false, fmt.Errorf("failed to write to DCC service: %w", err)
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		return false, fmt.Errorf("failed to read from DCC service: %w", err)
	}

	return buf[0] == 0x01, nil
}
