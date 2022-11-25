kubectl apply -f db/postgres-config.yaml
kubectl apply -f db/postgres-deployment.yaml
kubectl apply -f db/postgres-pvc-pv.yaml

kubectl apply -f environment/prod-config.yaml
kubectl apply -f environment/prod-secret.yaml

kubectl apply -f adapter.yaml
kubectl apply -f third-party/redis/redis-deployment.yaml

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
