# zms-devops-notification


Build and push to DockerHub
```shell
docker build -t devopszms2024/zms-devops-notification-service:latest .
docker build -t devopszms2024/zms-devops-angular-app:3.15 .
docker push devopszms2024/zms-devops-notification-service:latest
```

Create namespace & setup keycloak & notification-service infrastructure

```shell
minikube addons enable ingress
istioctl install --set profile=demo -y
```
First time you can use apply
```shell
kubectl apply -R -f notification-k8s 
kubectl apply -R -f notification-istio
```
When you want to replace existing pod, svc... you should use this command
```shell
kubectl replace --force -f notification-k8s
kubectl replace --force -f notification-istio
```

```shell
kubectl get pods -n backend
kubectl describe pods POD -n backend
```

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-kafka bitnami/kafka --namespace backend

helm install my-kafka bitnami/kafka --set persistence.size=8Gi,logPersistence.size=8Gi,replicaCount=3,volumePermissions.enabled=true,persistence.enabled=true,logPersistence.enabled=true,auth.clientProtocol=plaintext,serviceAccount.create=true,rbac.create=true,image.tag=latest

```