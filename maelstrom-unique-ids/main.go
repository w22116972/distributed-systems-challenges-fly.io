package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"strconv"
	"time"
)

func main() {
	n := maelstrom.NewNode()
	// Change the type to "generate"
	n.Handle("generate", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		// Change the type to "generate_ok" for responding
		body["type"] = "generate_ok"
		// Use like snowflake method, use node id with current timestamp to create unique ids
		body["id"] = msg.Src + ":" + msg.Dest + ":" + strconv.FormatInt(time.Now().UnixNano(), 10)

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		log.Fatal(err)
	}
}
