# Kubernetes Deployment

## Prerequisites

- kubectl
- Docker (to build the app image)
- Kubernetes cluster (minikube, kind, or cloud provider)

## Build and load the app image

### Minikube
```bash
eval $(minikube docker-env)
docker build -t todo-app:latest .
```

### Kind
```bash
docker build -t todo-app:latest .
kind load docker-image todo-app:latest
```

### Use a container registry
```bash
docker build -t your-registry/todo-app:latest .
docker push your-registry/todo-app:latest
```
Then edit `app-deployment.yaml` to use your image and set `imagePullPolicy: Always`.

## Deploy

```bash
kubectl apply -k k8s/
```

Or apply individually:
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/postgres-secret.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml
kubectl apply -f k8s/app-secret.yaml
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/app-service.yaml
kubectl apply -f k8s/ingress.yaml
```

## Access

**Port-forward (no Ingress):**
```bash
kubectl port-forward -n todo svc/todo-app 8080:80
```
Then open http://localhost:8080

**With Ingress (nginx ingress controller):**
Add `127.0.0.1 todo.local` to `/etc/hosts`, then:
```bash
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 8080:80
```
Open http://todo.local:8080

## Delete

```bash
kubectl delete -k k8s/
```
