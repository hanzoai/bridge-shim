# Hanzo Bridge Shim

White-label tenant shim for the canonical bridge — `bridge.hanzo.network`.

The runtime binary is `ghcr.io/luxfi/bridge` (the OSS upstream). This
repo ships:

- `tenant.yaml` — the declarative Hanzo Bridge configuration (brand,
  IAM endpoint `hanzo.id`, KMS endpoint `kms.hanzo.ai`, MPC cluster,
  strict-pq profile, supported chains, basket allowlist, fee receiver,
  domain `bridge.hanzo.network`, per-family release-pool sizing).
- `Dockerfile` — overlays `tenant.yaml` on top of the upstream image
  and pins `--tenant-config /etc/bridge/tenant.yaml` at the entrypoint.
- `contracts/tenant.json` + `contracts/Deploy.sh` — wraps
  `@luxfi/standard v1.7.5+`'s `DeployTenant.s.sol` with the Hanzo
  manifest.
- `k8s/` — Deployment, Service, IngressRoute, ConfigMap manifests
  targeting `bridge.hanzo.network`.
- `tenant_test.go` — single Go test that imports
  `github.com/luxfi/bridge/pkg/tenant` and validates `tenant.yaml`.

## Why a shim and not a fork

White-label by composition, never by fork. New OSS features land in
`luxfi/bridge` and we pick them up via a Dockerfile `FROM` tag bump —
no merge conflicts, no drift, no extra audit surface.

## Deploy

```bash
go test ./...                                # validate config
docker build -t ghcr.io/hanzoai/bridge-shim:v0.1.0 .
kubectl apply -k k8s/
HANZO_PRIVATE_KEY=0x... ./contracts/Deploy.sh https://rpc.hanzo.network
```

## Upstream pins

- `ghcr.io/luxfi/bridge`: v2.0.0
- `@luxfi/standard`: v1.7.5
