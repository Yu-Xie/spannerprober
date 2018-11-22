### Useful Kubernetes Commands
#### Check Current Kubernetes Cluster Context
```
kubectl config current-context
```
#### Switch to Another Zone K8S Cluster
See a complete list of GCP regions and zones [here](https://cloud.google.com/compute/docs/regions-zones/).
```
gcloud config set compute/zone us-central1-a
gcloud container clusters create [CLUSTER_NAME]
gcloud container clusters get-credentials u[CLUSTER_NAME]
```
#### Upload Secret to GKE cluster
Assuming your local copy of credentials JSON file is at `config/credentials.json`,
```
kubectl create secret generic credentials-key --from-file=key.json=config/credentials.json
```
#### Deploy/Delete Spaner Prober using Current K8S
```
kubectl create -f config/kube.yaml
kubectl delete deployment spanner-prober
```
#### Check Deployment(Pod) Status
```
➜  spannerprober git:(master) ✗ kubectl get pods
NAME                              READY     STATUS    RESTARTS   AGE
spanner-prober-65b5b4445d-cc6wr   1/1       Running   0          7s
```

#### Check Logs
```
kubectl logs -lapp="spanner-prober-app"
```
