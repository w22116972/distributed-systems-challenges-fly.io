package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"os"
)

func main() {
	n := maelstrom.NewNode()
	// register a handler callback function
	n.Handle("echo", func(msg maelstrom.Message) error {
		// Unmarshal the message body as a loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		// Update the message type to return back.
		body["type"] = "echo_ok"
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})
	// delegate execution to the Node by calling its Run()
	// continuously reads messages from STDIN and fires off a goroutine for each one to the associated handler
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		log.Fatal(err)
		os.Exit(1)
	}

}
