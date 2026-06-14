package authz

# Default deny
default allow = false

# Allow mutation only if the request ID is backed by a verified causal chain
allow {
    # Custom DCC built-in call
    # This synchronizes with the host kernel via Unix Domain Socket
    dcc.is_verified(input.request_id)
}
