![Build and test](https://github.com/okr-go-club/blockchain/actions/workflows/go.yml/badge.svg)
![golangci-lint](https://github.com/okr-go-club/blockchain/actions/workflows/golangci-lint.yml/badge.svg)

# blockchain

Educational project for learning golang and blockchain technology

### Run
```shell
go run main.go p2p.go chain.go -address localhost:8080

go run main.go p2p.go chain.go -address localhost:8081 -peers localhost:8080

go run main.go p2p.go chain.go -address localhost:8082 -peers localhost:8080,localhost:8081
```
