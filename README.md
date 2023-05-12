# distributed-systems-challenges-fly.io
Gossip Glomers: A series of distributed systems challenges, https://fly.io/dist-sys/

### Install

```shell
brew install graphviz gnuplot
```

Download from https://github.com/jepsen-io/maelstrom/releases/tag/v0.2.3

Install go dependency
```shell
# Same directory with go.mod
go get github.com/jepsen-io/maelstrom/demo/go
```

### To start new challenge

```shell
mkdir "xxxChallenge"
cd "xxxChallenge"
go mod init "xxxChallenge"
go mod tidy
go get github.com/jepsen-io/maelstrom/demo/go
# create `main.go` with `func main() { }`
```
