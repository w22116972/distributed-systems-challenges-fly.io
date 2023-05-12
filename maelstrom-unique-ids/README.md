# Unique IDs challenge

### How 

Use like snowflake method: source node id + destination node id + current timestamp

### Test
```shell
./maelstrom test -w unique-ids --bin ~/go/bin/maelstrom-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```
