# Hanzo Bridge — fork template

Reference configuration + deployment scaffold for white-label bridge
consumers. **NOT** a separate deployment image. Hanzo's own production
bridge runs the upstream `ghcr.io/luxfi/bridge` image directly with
tenant config supplied at runtime via env + ConfigMap.

## What lives here

| File | Purpose |
|---|---|
| `tenant.yaml` | Declarative tenant config — brand, IAM endpoint `hanzo.id`, KMS endpoint `kms.hanzo.ai`, MPC cluster, strict-pq profile, supported chains, basket allowlist, fee receiver, domain `bridge.hanzo.network`, per-family release-pool sizing. |
| `contracts/tenant.json` + `contracts/Deploy.sh` | Wraps `@luxfi/standard v1.7.5+`'s `DeployTenant.s.sol` with the Hanzo manifest. |
| `k8s/` | Deployment / Service / IngressRoute / ConfigMap manifests targeting `bridge.hanzo.network`. |
| `tenant_test.go` | Validates `tenant.yaml` against the upstream `github.com/luxfi/bridge/pkg/tenant` schema. |
| `Dockerfile` | **Reference only.** Documents how a downstream fork that needs its own image (e.g. compliance bake-in like `partner/bridge`) would assemble one. Hanzo production does not build or publish this image. |

## Production deployment — config, not image

Hanzo production pulls the canonical upstream image and supplies
tenant config at deploy time:

```yaml
# k8s/deployment.yaml — runtime config injection, not image bake
spec:
  template:
    spec:
      containers:
        - name: bridge
          image: ghcr.io/luxfi/bridge:v1.1.40        # clean upstream semver
          args: ["--tenant-config", "/etc/bridge/tenant.yaml"]
          volumeMounts:
            - name: tenant-config
              mountPath: /etc/bridge
      volumes:
        - name: tenant-config
          configMap:
            name: hanzo-bridge-tenant   # tenant.yaml in this ConfigMap
```

ConfigMap rotation is hot-reloadable — change `hanzo-bridge-tenant`,
no image rebuild required. Same image runs every Hanzo environment
(testnet, mainnet, dev) — environments differ only by ConfigMap.

## Why this repo exists

White-label by composition, never by fork. New OSS features land in
`luxfi/bridge`; downstream consumers pick them up by bumping the image
tag in their deployment manifest. This repo is the **template** any
new consumer can clone:

```
git clone https://github.com/hanzoai/bridge       # or zooai/bridge
# edit tenant.yaml for your brand + endpoints
# point your k8s manifests at ghcr.io/luxfi/bridge:vX.Y.Z
# done.
```

If a downstream needs to bake config into an image for compliance
reasons (the Liquidity pattern — US ATS/BD/TA requires region-locked
GAR image with config bake-in), the `Dockerfile` here shows the
minimal scaffold.

## Upstream pin

`luxfi/bridge` `v1.1.40` — the image is at `ghcr.io/luxfi/bridge:v1.1.40`.
Bump it in `k8s/deployment.yaml` to pick up new OSS features. Use clean
semver — no tenant suffix.
