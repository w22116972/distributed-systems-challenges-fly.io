package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync"
)

func main() {
	node := maelstrom.NewNode()
	server := &Server{node: node, messagesMap: make(map[int]interface{})}

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
	messagesMap map[int]interface{} // Hashset of golang version
	messageLock sync.RWMutex

	topology     map[string][]string // mapping from the node ID to list of its neighbor IDs
	topologyLock sync.RWMutex
}

func (server *Server) handleBroadcastMessage(message maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(message.Body, &body); err != nil {
		return err
	}

	// Update messageMap
	messageValue := int(body["message"].(float64))
	server.messageLock.Lock()
	if _, exists := server.messagesMap[messageValue]; exists {
		server.messageLock.Unlock()
		return server.node.Reply(message, map[string]any{
			"type": "broadcast_ok",
		})
	}
	server.messagesMap[messageValue] = struct{}{}
	server.messageLock.Unlock()

	// broadcast to neighbor node ID in the same cluster
	for _, neighborIDs := range server.topology[server.node.ID()] {
		// Since goroutine only created at ending of loop, so neighborIDs will only have one value
		// that is why we need a new local value at every iteration
		targetNodeID := neighborIDs
		go func() {
			if err := server.node.Send(targetNodeID, body); err != nil {
				panic(err)
			}
		}()
	}

	return server.node.Reply(message, map[string]any{
		"type": "broadcast_ok",
	})
}

func (server *Server) handleReadMessage(message maelstrom.Message) error {
	server.messageLock.Lock()
	// To extract keys as array from hashset in server.messagesMap
	messages := make([]int, 0, len(server.messagesMap))
	for key := range server.messagesMap {
		messages = append(messages, key)
	}
	server.messageLock.Unlock()
	return server.node.Reply(message, map[string]any{
		"type":     "read_ok",
		"messages": messages,
	})
}

type topologyMessage struct {
	Topology map[string][]string `json:"topology"`
}

func (server *Server) handleTopologyMessage(message maelstrom.Message) error {
	var topologyMessageBody topologyMessage
	if err := json.Unmarshal(message.Body, &topologyMessageBody); err != nil {
		return err
	}
	server.topologyLock.Lock()
	server.topology = topologyMessageBody.Topology
	server.topologyLock.Unlock()

	return server.node.Reply(message, map[string]any{
		"type": "topology_ok",
	})
}
