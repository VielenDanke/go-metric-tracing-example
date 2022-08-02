# go-metric-tracing-example
Jaeger + Prometheus + Kubernetes + Application

## Minikube

Install: https://minikube.sigs.k8s.io/docs/  

## Jaeger

1. Install cert-manager (https://cert-manager.io/docs/installation/)  
2. Enable ingress for minikube (minikube addons enable ingress)  
3. Install jaeger-operator (kubectl apply -f jaeger-operator.yaml -n "namespace").   
   1. It is in Cluster Wide mode. All Jaeger resources in all namespaces will be found.  
   2. To change Cluster Wide behavior - need to change ClusterRole and ClusterRoleBindings to Role, RoleBindings. Also need to place namespaces in jaeger-operator.yaml file under WATCH_NAMESPACE property comma-separated      
4. Install jaeger (kubectl apply -f jaeger.yaml -n "namespace")  
5. Build image to minikube (minikube build image -t users:1.0.0 .)  
6. Install application to minikube (helm install users -f chart/values.yaml chart/)
7. Modify /etc/hosts file (Add _minikube ip_ to file and associate it with users.local)  
8. Forward Jaeger UI to local port (kubectl port-forward -n "namespace" "pod-name" 9090:16686)  