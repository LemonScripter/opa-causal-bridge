import unittest
import time
import os

# OPA DCC Causal Enforcement Logic Verification
# This suite simulates the OPA built-in interaction with the DCC Causal Service.
# It ensures that Rego policies enforce Digital Causal Closure through fail-closed logic.

class TestOPADCCEnforcement(unittest.TestCase):
    def setUp(self):
        self.dcc_map = {}
        self.CAUSALITY_WINDOW_NS = 500 * 1000 * 1000 # 500ms

    def issue_token(self, request_id):
        # Simulation of the DCC kernel bridge population (global_dcc_map)
        self.dcc_map[request_id] = {
            "ts": time.time_ns(),
            "consumed": False
        }

    def rego_dcc_is_verified(self, request_id):
        """
        Simulation of the dcc.is_verified() Rego built-in.
        In production, this queries /var/run/bioos/dcc.sock.
        """
        now = time.time_ns()
        
        # 1. Verification Path: Fail-Closed if no token exists
        if request_id not in self.dcc_map:
            print(f"DCC [REJECT]: No causal origin for {request_id}")
            return False 
            
        token = self.dcc_map[request_id]
        
        # 2. Kernel Anchor: Check for expiration (Causality Window)
        if now - token["ts"] > self.CAUSALITY_WINDOW_NS:
            print(f"DCC [REJECT]: Stale intent for {request_id}")
            return False 
            
        # 3. Atomic Verification: Prevent Replay
        if token["consumed"]:
            print(f"DCC [REJECT]: Intent already consumed for {request_id}")
            return False 
            
        token["consumed"] = True 
        print(f"DCC [ALLOW]: Causal chain closed for {request_id}")
        return True 

    def test_reject_unauthorized_mutation(self):
        # Scenario: Autonomous/Orphaned request with no causal token
        request_id = "req-unauthorized"
        result = self.rego_dcc_is_verified(request_id)
        self.assertFalse(result)

    def test_allow_verified_mutation(self):
        # Scenario: Legitimate request backed by hardware-anchored intent
        request_id = "req-verified"
        self.issue_token(request_id)
        result = self.rego_dcc_is_verified(request_id)
        self.assertTrue(result)

    def test_prevent_stale_admission(self):
        # Scenario: Request initiated outside the 500ms causality window
        request_id = "req-stale"
        self.issue_token(request_id)
        # Manually aging the token
        self.dcc_map[request_id]["ts"] -= (self.CAUSALITY_WINDOW_NS + 1000000)
        result = self.rego_dcc_is_verified(request_id)
        self.assertFalse(result)

    def test_prevent_replay_attack(self):
        # Scenario: Attempting to reuse a single-use causal token
        request_id = "req-atomic"
        self.issue_token(request_id)
        self.assertTrue(self.rego_dcc_is_verified(request_id))
        self.assertFalse(self.rego_dcc_is_verified(request_id)) # Second attempt must fail

if __name__ == "__main__":
    print("--- Running OPA DCC Causal Enforcement Integration Tests ---")
    unittest.main()
