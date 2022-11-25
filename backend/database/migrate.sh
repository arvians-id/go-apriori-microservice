# Apriori Local Migration
migrate -path postgres/migration -database "postgres://root:root@localhost:5432/apriori?sslmode=disable" -verbose up

# Apriori Local Test Migration
migrate -path postgres/migration -database "postgres://root:root@localhost:5432/apriori_test?sslmode=disable" -verbose up