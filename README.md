#KubernetesSubmissions

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

### 3.9. DBaaS vs DIY

#### DBaaS (Google Cloud SQL) vs. DIY (Postgres on GKE + PVC)

| Dimension | DBaaS (Google Cloud SQL) | DIY (Postgres on GKE + PVC) |
| :--- | :--- | :--- |
| **Initialization Work** | **Low:** Can be initialized within minutes with a couple of commands via `gcloud` . | **High:** Requires writing custom Kubernetes manifests (`StatefulSet`, `PVC`, `StorageClass`, `Service`, `Secret`). |
| **Maintenance Effort** | **Low:** Provider handles node rollouts, PostgreSQL updates, automated replication/failover, and hardware health. | **High:** Node rollouts could cause downtime unless a complex multinode system is setup. Minor/major engine upgrades have to be done manually. Storage scaling, and HA failovers have to be taken care of. |
| **Backup Methods** | **Seamless:** Native point-in-time recovery (PITR) and scheduled automated snapshots via GCP console/API and zero downtime restores. | **Complex:** Custom backup cronjobs or backup tools are required |
| **Infrastructure Cost** | **Higher:** Infrastructure cost is higher due to included management layer and vendor markup. | **Lower:** Direct infrastructure cost is lower because it uses standard GKE worker node pool compute and persistent disks without management overhead markups |
| **Total Cost (for small teams)** | **Lower:** Higher hosting fees are offset by no labor hours required for upkeep. | **Higher:** Infrastructure savings are consumed by engineering overhead spent maintaining and recovering database state. |
