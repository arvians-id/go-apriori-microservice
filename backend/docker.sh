docker image rm arvians/apriori-adapter
docker image rm arvians/apriori-category
docker image rm arvians/apriori-apriori
docker image rm arvians/apriori-comment
docker image rm arvians/apriori-message-broker
docker image rm arvians/apriori-notification
docker image rm arvians/apriori-password-reset
docker image rm arvians/apriori-payment
docker image rm arvians/apriori-product
docker image rm arvians/apriori-transaction
docker image rm arvians/apriori-user-order
docker image rm arvians/apriori-user

docker build -t arvians/apriori-adapter -f adapter/Dockerfile adapter
docker build -t arvians/apriori-category -f services/category-service/Dockerfile services/category-service
docker build -t arvians/apriori-apriori -f services/apriori-service/Dockerfile services/apriori-service
docker build -t arvians/apriori-comment -f services/comment-service/Dockerfile services/comment-service
docker build -t arvians/apriori-message-broker -f services/message-service/Dockerfile services/message-service
docker build -t arvians/apriori-notification -f services/notification-service/Dockerfile services/notification-service
docker build -t arvians/apriori-password-reset -f services/password-reset-service/Dockerfile services/password-reset-service
docker build -t arvians/apriori-payment -f services/payment-service/Dockerfile services/payment-service
docker build -t arvians/apriori-product -f services/product-service/Dockerfile services/product-service
docker build -t arvians/apriori-transaction -f services/transaction-service/Dockerfile services/transaction-service
docker build -t arvians/apriori-user-order -f services/user-order-service/Dockerfile services/user-order-service
docker build -t arvians/apriori-user -f services/user-service/Dockerfile services/user-service