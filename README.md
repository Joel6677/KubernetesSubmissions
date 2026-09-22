# KubernetesSubmissions

## Exercises

### Chapter 2

- [1.1.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.1)
- [1.2.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.2)
- [1.3.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.3)
- [1.4.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.4)
- [1.5.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.5)
- [1.6.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.6)
- [1.7.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.7)
- [1.8.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.8)
- [1.9.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.9)
- [1.10.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.10)
- [1.11.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.11)
- [1.12.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.12)
- [1.13.](https://github.com/Joel6677/KubernetesSubmissions/tree/1.13)

### Chapter 3

- [2.1.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.1)
- [2.2.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.2)
- [2.3.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.3)
- [2.4.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.4)
- [2.5.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.5)
- [2.6.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.6)
- [2.7.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.7)
- [2.8.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.8)
- [2.9.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.9)
- [2.10.](https://github.com/Joel6677/KubernetesSubmissions/tree/2.10)

### Chapter 4

- [3.1.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.1)
- [3.2.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.2)
- [3.3.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.3)
- [3.4.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.4)
- [3.5.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.5)
- [3.6.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.6)
- [3.7.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.7)
- [3.8.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.8)
- [3.9.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.9)
- [3.10.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.10)
- [3.11.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.11)
- [3.12.](https://github.com/Joel6677/KubernetesSubmissions/tree/3.12)

### Chapter 5

- [4.1.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.1)
- [4.2.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.2)
- [4.3.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.3/README.md)
- [4.4.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.4/ping-pong)
- [4.5.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.5)
- [4.6.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.6)
- [4.7.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.7)
- [4.8.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.8)
- [4.9.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.9)
- [4.10.](https://github.com/Joel6677/KubernetesSubmissions/tree/4.10)
- [4.10.](https://github.com/Joel6677/KubernetesSubmissions-project-config/tree/4.10)

### Chapter 6
- [5.1.](https://github.com/Joel6677/KubernetesSubmissions/tree/5.1/chapter_6/5.1)
- [5.2.](https://github.com/Joel6677/KubernetesSubmissions/tree/5.2)




## Prerequisites

- [k3d](https://k3d.io/) with a running cluster
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/) (used to install NATS, and optionally Prometheus/Grafana/ArgoCD) 
- [Envoy Gateway](https://gateway.envoyproxy.io/)

## Cluster setup (k3d)

```bash
k3d cluster create --agents 2 -p 8081:80@loadbalancer --k3s-arg '--disable=traefik@server:0'
```

```bash
kubectl apply --server-side -f https://github.com/envoyproxy/gateway/releases/latest/download/install.yaml
kubectl -n envoy-gateway-system rollout status deployment/envoy-gateway --timeout=180s
```

### NATS

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
helm upgrade --install my-nats nats/nats \
  --namespace nats \
  --create-namespace \
  --set promExporter.enabled=true
```

## Running the exercises (Log output & Ping-pong)

```bash
cd exercises_config/overlays/k3d
kubectl apply -k .
```

## Running the project (Todo App)

```bash
cd project/overlays/production   # or /staging
kubectl apply -k .
```

### Required secrets

```bash
kubectl create secret generic todo-postgres-credentials \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=<pick-a-password> \
  --from-literal=POSTGRES_DB=todos \
  -n <namespace>

kubectl create secret generic broadcaster-webhook \
  --from-literal=WEBHOOK_URL='<your Discord/Slack/generic webhook url>' \
  -n <namespace>
```

### 3.9. DBaaS vs DIY

#### DBaaS (Google Cloud SQL) vs. DIY (Postgres on GKE + PVC)

| Dimension | DBaaS (Google Cloud SQL) | DIY (Postgres on GKE + PVC) |
| :--- | :--- | :--- |
| **Initialization Work** | **Low:** Can be initialized within minutes with a couple of commands via `gcloud` . | **High:** Requires writing custom Kubernetes manifests (`StatefulSet`, `PVC`, `StorageClass`, `Service`, `Secret`). |
| **Maintenance Effort** | **Low:** Provider handles node rollouts, PostgreSQL updates, automated replication/failover, and hardware health. | **High:** Node rollouts could cause downtime unless a complex multinode system is setup. Minor/major engine upgrades have to be done manually. Storage scaling, and HA failovers have to be taken care of. |
| **Backup Methods** | **Seamless:** Native point-in-time recovery (PITR) and scheduled automated snapshots via GCP console/API and zero downtime restores. | **Complex:** Custom backup cronjobs or backup tools are required |
| **Infrastructure Cost** | **Higher:** Infrastructure cost is higher due to included management layer and vendor markup. | **Lower:** Direct infrastructure cost is lower because it uses standard GKE worker node pool compute and persistent disks without management overhead markups |
| **Total Cost (for small teams)** | **Lower:** Higher hosting fees are offset by no labor hours required for upkeep. | **Higher:** Infrastructure savings are consumed by engineering overhead spent maintaining and recovering database state. |

### 4.3. Prometheus query

sum(kube_pod_info{namespace="monitoring", created_by_kind="StatefulSet"})



