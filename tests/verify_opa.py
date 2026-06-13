import unittest
import time

# OPA DCC Causal Enforcement Logic Verification
# Ensures that Rego policies can effectively block non-causal mutations.

class TestOPADCCEnforcement(unittest.TestCase):
    def setUp(self):
        self.dcc_map = {}
        self.CAUSALITY_WINDOW_NS = 500 * 1000 * 1000 # 500ms

    def issue_token(self, request_id):
        self.dcc_map[request_id] = {
            "ts": time.time_ns(),
            "consumed": False
        }

    def rego_dcc_is_verified(self, request_id):
        # Simulation of the dcc.is_verified() Rego built-in
        now = time.time_ns()
        
        if request_id not in self.dcc_map:
            return False # REJECT (No causal origin)
            
        token = self.dcc_map[request_id]
        if now - token["ts"] > self.CAUSALITY_WINDOW_NS:
            return False # REJECT (Stale intent)
            
        if token["consumed"]:
            return False # REJECT (Replay)
            
        token["consumed"] = True # OPA evaluation consumes the intent
        return True # ALLOW (Causal chain closed)

    def test_reject_unauthorized_mutation(self):
        # Request with no causal token must be rejected by Rego
        request_id = "req-123"
        result = self.rego_dcc_is_verified(request_id)
        self.assertFalse(result)

    def test_allow_verified_mutation(self):
        # Request with fresh causal token must be allowed
        request_id = "req-123"
        self.issue_token(request_id)
        result = self.rego_dcc_is_verified(request_id)
        self.assertTrue(result)

    def test_prevent_stale_admission(self):
        # Admission requests outside the causality window must fail
        request_id = "req-123"
        self.issue_token(request_id)
        self.dcc_map[request_id]["ts"] -= (self.CAUSALITY_WINDOW_NS + 1000)
        result = self.rego_dcc_is_verified(request_id)
        self.assertFalse(result)

if __name__ == "__main__":
    print("--- Running OPA DCC Causal Enforcement Tests ---")
    unittest.main()
