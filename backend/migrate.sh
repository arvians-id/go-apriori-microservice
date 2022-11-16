# Apriori Migration
migrate -path database/postgres/migration -database "postgres://root:root@localhost:5432/apriori?sslmode=disable" -verbose up

# Apriori Test Migration
migrate -path database/postgres/migration -database "postgres://root:root@localhost:5432/apriori_test?sslmode=disable" -verbose up