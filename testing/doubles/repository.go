package doubles

import (
	"context"

	"github.com/flohansen/nova-cloud/internal/domain"
)

type TestNodeRepository struct {
	Nodes map[string]domain.Node
}

func NewTestNodeRepository() *TestNodeRepository {
	return &TestNodeRepository{
		Nodes: make(map[string]domain.Node),
	}
}

func (t *TestNodeRepository) WithNodes(nodes ...domain.Node) *TestNodeRepository {
	for _, node := range nodes {
		t.CreateOrUpdate(nil, node)
	}
	return t
}

func (t *TestNodeRepository) CreateOrUpdate(ctx context.Context, node domain.Node) error {
	t.Nodes[node.NodeID] = node
	return nil
}

func (t *TestNodeRepository) Delete(ctx context.Context, nodeID string) error {
	delete(t.Nodes, nodeID)
	return nil
}

func (t *TestNodeRepository) FindAll(ctx context.Context) ([]domain.Node, error) {
	nodes := make([]domain.Node, 0, len(t.Nodes))
	for _, node := range t.Nodes {
		nodes = append(nodes, node)
	}
	return nodes, nil
}
