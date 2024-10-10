run:
	go run cmd/blockchain/main.go -address localhost:8080 -http localhost:8090

generate_key:
	go run cmd/private_key_generator/main.go
