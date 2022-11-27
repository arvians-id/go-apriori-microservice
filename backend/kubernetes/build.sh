kubectl apply -f environment/prod-config.yaml
kubectl apply -f environment/prod-secret.yaml

kubectl apply -f adapter.yaml
kubectl apply -f third-party/redis/redis-deployment.yaml
kubectl apply -f third-party/db/postgres-deployment.yaml

kubectl apply -f services/apriori-deployment.yaml
kubectl apply -f services/category-deployment.yaml
kubectl apply -f services/comment-deployment.yaml
kubectl apply -f services/notification-deployment.yaml
kubectl apply -f services/password-reset-deployment.yaml
kubectl apply -f services/payment-deployment.yaml
kubectl apply -f services/product-deployment.yaml
kubectl apply -f services/transaction-deployment.yaml
kubectl apply -f services/user-order-deployment.yaml
kubectl apply -f services/user-deployment.yaml
kubectl apply -f services/message-broker-deployment.yaml

helm uninstall nsq-release
helm uninstall prometheus-release
helm uninstall grafana-release

helm install nsq-release nsqio/nsq --set nsqlookupd.replicaCount=1,nsqd.replicaCount=1,nsqadmin.replicaCount=1,nsqadmin.service.type=NodePort,nsqadmin.service.nodePort=30000
helm install prometheus-release prometheus/prometheus --set prometheus-node-exporter.hostRootFsMount.enabled=false
helm install grafana-release grafana/grafana --set persistence.enabled=false