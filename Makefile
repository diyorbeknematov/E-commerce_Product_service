CURRENT_DIR = $(shell pwd)

DB_URL := postgres://postgres:123321@localhost:5432/ecommerce_auth_service?sslmode=disable

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

swag-init:
	swag init -g api/router.go --output api/handler/docs


mig-up:
	migrate -path db/migrations -database '${DB_URL}' -verbose up

mig-down:
	migrate -path db/migrations -database '${DB_URL}' -verbose down

mig-force:
	migrate -path db/migrations -database '${DB_URL}' -verbose force 1

mig-create-users:
	migrate create -ext sql -dir db/migrations -seq auth_service_table

swag-gen:
	~/go/bin/swag init -g api/router.go -o api/handler/docs
