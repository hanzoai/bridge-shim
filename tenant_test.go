// Package hanzobridgeshim tests the Hanzo Bridge tenant configuration
// against the canonical upstream tenant.Config schema. The shim itself
// ships no Go binary — the production image overlays this tenant.yaml
// on top of the upstream OSS bridge container.
package hanzobridgeshim

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/luxfi/bridge/pkg/tenant"
)

// TestTenantConfigValidates is the single guardrail this shim ships.
// Loads tenant.yaml via the upstream tenant.Load and asserts every
// required field is set + every constraint is satisfied.
func TestTenantConfigValidates(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	cfgPath := filepath.Join(filepath.Dir(thisFile), "tenant.yaml")

	cfg, err := tenant.Load(cfgPath)
	if err != nil {
		t.Fatalf("tenant.Load(%s): %v", cfgPath, err)
	}

	if cfg.Brand.Slug != "hanzo-bridge" {
		t.Errorf("Brand.Slug: got %q, want hanzo-bridge", cfg.Brand.Slug)
	}
	if cfg.Network.ID != 36963 {
		t.Errorf("Network.ID: got %d, want 36963 (Hanzo Mainnet)", cfg.Network.ID)
	}
	if cfg.IAM.Endpoint != "https://hanzo.id" {
		t.Errorf("IAM.Endpoint: got %q, want https://hanzo.id", cfg.IAM.Endpoint)
	}
	if cfg.KMS.Endpoint != "https://kms.hanzo.ai" {
		t.Errorf("KMS.Endpoint: got %q, want https://kms.hanzo.ai", cfg.KMS.Endpoint)
	}
	if cfg.PQProfile != "strict-pq" {
		t.Errorf("PQProfile: got %q, want strict-pq", cfg.PQProfile)
	}
	if cfg.Domain != "bridge.hanzo.network" {
		t.Errorf("Domain: got %q, want bridge.hanzo.network", cfg.Domain)
	}
}

// TestSupportedChainsScopedToHanzo asserts the chain allowlist
// matches Hanzo's deployment scope.
func TestSupportedChainsScopedToHanzo(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	cfgPath := filepath.Join(filepath.Dir(thisFile), "tenant.yaml")
	cfg, err := tenant.Load(cfgPath)
	if err != nil {
		t.Fatalf("tenant.Load: %v", err)
	}
	for _, c := range []string{"eth", "btc", "sol", "ton", "xrp", "dot", "hanzo"} {
		if !cfg.IsChainSupported(c) {
			t.Errorf("expected chain %q to be supported", c)
		}
	}
	for _, c := range []string{"liquid", "zoo", "pars", "lux"} {
		if cfg.IsChainSupported(c) {
			t.Errorf("chain %q must NOT be supported by Hanzo Bridge", c)
		}
	}
}
