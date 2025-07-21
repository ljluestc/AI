// Package knowledgegraph provides a knowledge graph for code analysis
package knowledgegraph

import (
	"context"
	"errors"
)

// Node represents a node in the knowledge graph
type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// Edge represents an edge in the knowledge graph
type Edge struct {
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// Graph represents a knowledge graph
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Service provides operations for working with the knowledge graph
type Service interface {
	GetNode(ctx context.Context, id string) (*Node, error)
	CreateNode(ctx context.Context, node *Node) error
	CreateEdge(ctx context.Context, edge *Edge) error
	Query(ctx context.Context, query string, params map[string]interface{}) (*Graph, error)
}

// memoryService is an in-memory implementation of the Service interface
type memoryService struct {
	nodes map[string]Node
	edges []Edge
}

// NewMemoryService creates a new in-memory knowledge graph service
func NewMemoryService() Service {
	return &memoryService{
		nodes: make(map[string]Node),
		edges: make([]Edge, 0),
	}
}

func (s *memoryService) GetNode(ctx context.Context, id string) (*Node, error) {
	node, ok := s.nodes[id]
	if !ok {
		return nil, errors.New("node not found")
	}
	return &node, nil
}

func (s *memoryService) CreateNode(ctx context.Context, node *Node) error {
	if node.ID == "" {
		return errors.New("node ID cannot be empty")
	}
	s.nodes[node.ID] = *node
	return nil
}

func (s *memoryService) CreateEdge(ctx context.Context, edge *Edge) error {
	if _, ok := s.nodes[edge.Source]; !ok {
		return errors.New("source node not found")
	}
	if _, ok := s.nodes[edge.Target]; !ok {
		return errors.New("target node not found")
	}
	s.edges = append(s.edges, *edge)
	return nil
}

func (s *memoryService) Query(ctx context.Context, query string, params map[string]interface{}) (*Graph, error) {
	// This is a simplified implementation
	// In a real system, you would parse and execute the query

	// Just return all nodes and edges for now
	nodes := make([]Node, 0, len(s.nodes))
	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}

	return &Graph{
		Nodes: nodes,
		Edges: s.edges,
	}, nil
}
