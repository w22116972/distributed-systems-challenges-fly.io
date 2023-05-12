# Challenge #1: Echo

### After developing program
```shell
go install .
```

Then our `~/go/bin` will have our program `~/go/bin/maelstrom-echo`

### Test it on Maelstrom
```shell
../maelstrom/maelstrom test -w echo --bin ~/go/bin/maelstrom-echo --node-count 1 --time-limit 10
# Ok if seeing "Everything looks good! ヽ(‘ー`)ノ"
```

### Debug
```shell
../maelstrom/maelstrom serve
# Go to localhost:8080
```
