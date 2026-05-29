# Hanzo Bridge

Tenant deployment of the canonical bridge — `bridge.hanzo.network`.

> The "shim" framing in earlier revs of this repo was misleading. This
> is not a shim — it's the Hanzo Bridge **tenant deployment**: brand,
> network endpoints, per-tenant infra wiring. The runtime binary IS
> `ghcr.io/luxfi/bridge` (the OSS upstream), composed at deploy time
> via tenant.yaml. **Rename target: `hanzoai/bridge-shim` → `hanzoai/bridge`.**

This repo ships:

- `tenant.yaml` — declarative Hanzo Bridge configuration (brand,
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

## Why composition, not a fork

White-label by composition, never by fork. New OSS features land in
`luxfi/bridge` and we pick them up via a Dockerfile `FROM` tag bump —
no merge conflicts, no drift, no extra audit surface.

## Image tagging convention

Tag the tenant image by the upstream version it composes:

```
ghcr.io/hanzoai/bridge:v1.1.40-hanzo
```

NOT an independent semver. The Dockerfile `FROM` pin IS the contract
— bumping a local `v0.x.y` every time upstream bumps is double-
bookkeeping with no extra information. (Legacy `v0.1.x` and `v0.2.x`
tags exist on this repo from the older convention and remain in
history; new deploys use the `vX.Y.Z-hanzo` form.)

## Deploy

```bash
go test ./...                                          # validate config
docker build -t ghcr.io/hanzoai/bridge:v1.1.40-hanzo . # tag matches upstream
kubectl apply -k k8s/
HANZO_PRIVATE_KEY=0x... ./contracts/Deploy.sh https://rpc.hanzo.network
```

## Upstream pins

- `ghcr.io/luxfi/bridge`: v1.1.40
- `@luxfi/standard`: v1.7.5
