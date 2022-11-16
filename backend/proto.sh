# Run Proto
protoc adapter/proto/*.proto --go_out=plugins=grpc:.
protoc services/apriori-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/category-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/comment-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/notification-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/password-reset-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/payment-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/product-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/transaction-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/user-order-service/proto/*.proto --go_out=plugins=grpc:.
protoc services/user-service/proto/*.proto --go_out=plugins=grpc:.

#Inject Services:
protoc-go-inject-tag -input=./services/*/pb/*.pb.go

#Inject Adapter:
protoc-go-inject-tag -input=./adapter/pb/*.pb.go