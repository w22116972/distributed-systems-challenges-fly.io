#!/bin/bash
mkdir "maelstrom-$1"
cd "maelstrom-$1" || exit
go mod init "maelstrom-$1"
go mod tidy
go get github.com/jepsen-io/maelstrom/demo/go
echo "package main" >> main.go
echo "#!/bin/bash" >> test.sh
echo "go install ." >> test.sh
chmod +x test.sh
