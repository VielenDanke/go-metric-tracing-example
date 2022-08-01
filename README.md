# go-metric-tracing-example
Jaeger + Prometheus + Kubernetes + Application

## Jaeger

1. Install cert-manager (https://cert-manager.io/docs/installation/)  
2. Enable ingress for minikube (minikube addons enable ingress)  
3. Install jaeger-operator (kubectl apply -f jaeger-operator.yaml -n "namespace")  
4. Install jaeger (kubectl apply -f jaeger.yaml -n "namespace")  
5. Build image to minikube (minikube build image -t users:1.0.0 .)  
6. Install application to minikube (helm install users -f chart/values.yaml chart/)
7. Modify /etc/hosts file (Add _minikube ip_ to file and associate it with users.local)  