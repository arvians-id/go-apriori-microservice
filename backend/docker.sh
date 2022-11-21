docker image rm apriori-adapter
docker image rm apriori-apriori
docker image rm apriori-category
docker image rm apriori-comment
docker image rm apriori-message-broker
docker image rm apriori-notification
docker image rm apriori-password-reset
docker image rm apriori-payment
docker image rm apriori-product
docker image rm apriori-transaction
docker image rm apriori-user-order
docker image rm apriori-user

docker build . -t apriori-adapter -f adapter/Dockerfile
docker build . -t apriori-apriori -f services/apriori-service/Dockerfile
docker build . -t apriori-category -f services/category-service/Dockerfile
docker build . -t apriori-comment -f services/comment-service/Dockerfile
docker build . -t apriori-message-broker -f services/message-service/Dockerfile
docker build . -t apriori-notification -f services/notification-service/Dockerfile
docker build . -t apriori-password-reset -f services/password-reset-service/Dockerfile
docker build . -t apriori-payment -f services/payment-service/Dockerfile
docker build . -t apriori-product -f services/product-service/Dockerfile
docker build . -t apriori-transaction -f services/transaction-service/Dockerfile
docker build . -t apriori-user-order -f services/user-order-service/Dockerfile
docker build . -t apriori-user -f services/user-service/Dockerfile

