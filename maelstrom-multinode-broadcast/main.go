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
	node                  *maelstrom.Node
	messagesMap           map[int]interface{} // Hashset of golang version
	messageLock           sync.RWMutex
	nodeIDToKnownMessages map[string][]int // This server knows that for node ID n, it knows(seen) these messagesMap

	topology     map[string]interface{}
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
		return nil
	}
	server.messagesMap[messageValue] = struct{}{}
	server.messageLock.Unlock()

	// broadcast to neighbor node ID in the same cluster
	for _, knownNodeIDInSameCluster := range server.node.NodeIDs() {
		if knownNodeIDInSameCluster == message.Src || knownNodeIDInSameCluster == server.node.ID() {
			continue
		}

		// Since goroutine only created at ending of loop, so knownNodeIDInSameCluster will only have one value
		// that is why we need a new local value at every iteration
		targetNodeID := knownNodeIDInSameCluster
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
	server.messageLock.RLock()
	// To extract keys as array from hashset in server.messagesMap
	messages := make([]int, 0, len(server.messagesMap))
	for key := range server.messagesMap {
		messages = append(messages, key)
	}
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
