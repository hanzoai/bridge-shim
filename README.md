# Hanzo Bridge — DEPRECATED REPO

This repo is **no longer needed**. Hanzo's bridge deployment runs the
canonical `ghcr.io/luxfi/bridge` image directly with tenant config
supplied at runtime via environment variables (or a ConfigMap-mounted
tenant.yaml when one is convenient — both are first-class in
luxfi/bridge's declarative Bridge SDK pattern).

There is no Hanzo-specific Go code, no Hanzo-specific Docker image,
no Hanzo-specific anything that justifies a separate repo. Tenant
identity is config, not a code artifact.

## Where Hanzo's bridge config lives now

In Hanzo's platform infrastructure repo (`~/work/hanzo/platform` →
`platform.hanzo.ai`), as ConfigMap + Secret manifests applied to the
`hanzo-bridge` namespace on whichever cluster runs Hanzo's bridge.
The deployment YAML pulls the canonical upstream image:

```yaml
spec:
  containers:
    - name: bridge
      image: ghcr.io/luxfi/bridge:v1.1.40
      envFrom:
        - configMapRef:
            name: hanzo-bridge-config
        - secretRef:
            name: hanzo-bridge-secrets
```

Same image every Hanzo environment runs (testnet, mainnet, dev) —
environments differ only by ConfigMap.

## What was here (historical)

The repo formerly tried to be a "shim" — a Dockerfile that did
`FROM ghcr.io/luxfi/bridge:vX.Y.Z` plus a baked tenant.yaml. That
pattern is the right move only when compliance requires the config
baked into a region-locked image (US ATS/BD/TA regulated case).
Hanzo has no such constraint; runtime config wins.

The Dockerfile + tenant.yaml + k8s manifests remain in tree as a
reference for any future contributor who wants to read how a
compliance-baked variant composes. None of this is the Hanzo
deployment path.

## Looking for the Hanzo bridge SDK?

That's the browser-side `@luxfi/bridge` package — same package every
Lux tenant consumes. There is no `@hanzoai/bridge`. White-label
branding flows through `@luxfi/brand` at build time; the SDK itself
is brand-neutral by design.

```ts
import { mountBridge } from '@luxfi/bridge'
import hanzoBrand from '@hanzoai/brand/brand.json'

mountBridge({
  config: {
    apiHost: 'https://api.bridge.hanzo.network',
    env: 'mainnet',
    brand: {
      name: `${hanzoBrand.brand.shortName} Bridge`,
      primaryColor: hanzoBrand.brand.primaryColor,
    },
  },
})
```

## Repo disposition

Recommended: archive this repo. The naming (`hanzoai/bridge-shim`)
implies an artifact that no longer exists. Nothing here is load-
bearing for Hanzo production.
