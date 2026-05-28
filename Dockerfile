# Hanzo Bridge production image.
#
# The shim ships zero Go code at runtime — the binary is the upstream
# OSS bridge. We overlay our tenant.yaml at /etc/bridge/tenant.yaml
# and instruct the entrypoint to load it via --tenant-config.
#
# Pin the upstream version explicitly. NEVER use :latest, :main, or
# floating tags in production manifests per global CLAUDE.md.
FROM ghcr.io/luxfi/bridge:v1.1.39

# Tenant configuration. ConfigMap mount in K8s overrides this for
# rotating prod credentials without an image rebuild; the bake-in is
# the default and dev fallback.
COPY tenant.yaml /etc/bridge/tenant.yaml

# Re-declare the entrypoint with the tenant-config flag pinned.
ENTRYPOINT ["/usr/local/bin/bridge", "--tenant-config", "/etc/bridge/tenant.yaml"]
