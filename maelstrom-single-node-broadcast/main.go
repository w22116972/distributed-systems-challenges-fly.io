package main

import (
	"encoding/json"
	"fmt"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync"
)

func main() {
	node := maelstrom.NewNode()
	server := &Server{node: node, messages: make([]int, 0)}

	node.Handle("broadcast", server.handleBroadcastMessage)
	node.Handle("read", server.handleReadMessage)
	node.Handle("topology", server.handleTopologyMessage)

	if err := node.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		log.Fatal(err)
	}
}

type Server struct {
	node        *maelstrom.Node
	messages    []int
	messageLock sync.RWMutex

	topology     map[string]interface{}
	topologyLock sync.RWMutex
}

//type topologyMessage struct {
//	Topology map[string][]string `json:"topology"`
//}

func (server *Server) handleBroadcastMessage(message maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(message.Body, &body); err != nil {
		return err
	}
	server.messageLock.Lock()
	server.messages = append(server.messages, int(body["message"].(float64)))
	server.messageLock.Unlock()

	return server.node.Reply(message, map[string]any{
		"type": "broadcast_ok",
	})
}

func (server *Server) handleReadMessage(message maelstrom.Message) error {
	server.messageLock.RLock()
	messages := server.messages
	server.messageLock.RUnlock()
	return server.node.Reply(message, map[string]any{
		"type":     "read_ok",
		"messages": messages,
	})
}

func (server *Server) handleTopologyMessage(message maelstrom.Message) error {
	var topologyMessageBody map[string]interface{}
	if err := json.Unmarshal(message.Body, &topologyMessageBody); err != nil {
		return err
	}
	topology, ok := topologyMessageBody["topology"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid topology format")
	}
	server.topologyLock.Lock()
	server.topology = topology
	server.topologyLock.Unlock()

	return server.node.Reply(message, map[string]any{
		"type": "topology_ok",
	})
}
