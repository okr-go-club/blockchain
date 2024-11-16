![Build](https://github.com/okr-go-club/blockchain/actions/workflows/go.yml/badge.svg)
![Tests](https://github.com/okr-go-club/blockchain/actions/workflows/go-test/badge.svg)
![golangci-lint](https://github.com/okr-go-club/blockchain/actions/workflows/golangci-lint.yml/badge.svg)

# Blockchain

An educational project for learning `Golang` and `blockchain` technology. This project uses `Golang`, `net/http` for networking, and `Badger` as persistent storage.

### Contents:

- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [Running app](#Run)
- [Generate Private Key for Testing](#Generate-Private-Key-for-Testing)
- [UI](#UI)
- [Contributing](#Contributing)

### Prerequisites
Before running the project, make sure you have the following installed:
- [Golang](https://golang.org/doc/install) (version 1.22.2 or higher)
- [Node.js](https://nodejs.org/) (version 14 or higher)
- [npm](https://www.npmjs.com/get-npm) (comes with Node.js)

### Installation
1. Clone the repository:
    ```shell
    git clone https://github.com/okr-go-club/blockchain.git
    cd blockchain
    ```
2. Install dependencies for the backend:
    ```shell
    go mod tidy
    ```
3. Install dependencies for the frontend:
    ```shell
    cd frontend
    npm install
    ```
4. Run the project (as described in the "Run" section below).

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

### Contributing
Contributions are welcome! If you would like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

Please make sure your code follows the existing style and passes all linting checks before submitting a PR.
