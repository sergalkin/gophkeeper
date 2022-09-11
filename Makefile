.PHONY: binary proto build-builder build-server run-server

build-builder:
	cd build/package && docker compose build builder

binary:
	cd build/package && docker compose run --rm builder

build-server:
	cd build/package && docker compose build server

run-server:
	 cd build/package && docker compose up server pgsql

#example: make proto name=user
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/$(name).proto

#example: make migrate-create name=users
migrate-create:
	migrate create -ext sql -dir internal/migrations -seq $(name)

#example: make migrate type=up
migrate:
	migrate -path internal/migrations -database "postgresql://root:root@localhost:5432/gophkeeper?sslmode=disable" -verbose $(type)

# example: make release V=0.0.0
release:
	echo v$(V)
	@read -p "Press enter to confirm and push to origin ..."