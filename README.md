![Build and test](https://github.com/okr-go-club/blockchain/actions/workflows/go.yml/badge.svg)
![golangci-lint](https://github.com/okr-go-club/blockchain/actions/workflows/golangci-lint.yml/badge.svg)

# Blockchain

An educational project for learning `Golang` and `blockchain` technology. This project uses `Golang`, `net/http` for networking, and `Badger` as persistent storage.

### Run
An example of running multiple nodes on a single machine:
```shell
go run cmd/blockchain/main.go -address localhost:8080 -http localhost:8090 -storage chain_storage

go run cmd/blockchain/main.go -address localhost:8081 -peers localhost:8080 -http localhost:8091 -storage chain_storage_2

go run cmd/blockchain/main.go -address localhost:8082 -peers localhost:8080,localhost:8081 -http localhost:8092 -storage chain_storage_3
```

### Generate Private Key for Testing
You can generate a private and public key for testing purposes:
```shell
go run cmd/private_key_generator/main.go
```

### Architecture
The blockchain implements a Bitcoin-like model. The blockchain and wallet entities are implemented in the `chain` package. Each node stores its own copy of the blockchain (in the `storage` package) and synchronizes it with others via peer-to-peer connections (using the `p2p` package).

### UI
The UI is built with TypeScript, React, and Chakra UI. To run a local instance:
```shell
cd frontend
npm start
```

The default port is `8090`, but you can change it in `frontend/src/axiosConfig.ts`.
