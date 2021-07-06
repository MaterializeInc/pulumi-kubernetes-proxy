# Kubernetes Proxy Pulumi Provider

A [Pulumi](https://pulumi.com) provider that proxies a port on the host to a
port on a pod in the cluster during provisioning.

**Warning:** The way this fits into the Pulumi architecture is sketchy. It might
break in future Pulumi releases. Use at your own risk!

This can be useful to use Pulumi to provision services that run within the
cluster that are not publicly accessible, like a database.

## Usage example

To provision a PostgreSQL database running in Kubernetes:

```python
import pulumi_kubernetes_proxy as k8s_proxy

# Hardcode a port that is likely to be free. This needs to be stable (i.e., we
# can't just let the kernel allocate a free port) so that Pulumi doesn't
# perpetually show a diff in the PostgreSQL provider.
PORT = 32123

eks_cluster = eks.Cluster(...)

k8s_proxy.Provider(
    "postgresql-proxy",
    kubeconfig=eks_cluster.kubeconfig,
    host_port=PORT,
    remote_port=5432,   # Target the PostgreSQL port.
    namespace="default",
    pod_selector="workload=postgresql",
)

provider = postgresql.Provider(
    base_name,
    host="localhost",
    port=PORT,
    connect_timeout=10,
    database="db",
    username="user",
    password="pass",
    opts=pulumi.ResourceOptions(depends_on=[rds_proxy_provider]),
)
```
