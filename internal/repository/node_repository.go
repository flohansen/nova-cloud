package repository

import (
	"context"

	"github.com/flohansen/nova-cloud/internal/domain"
	"github.com/flohansen/nova-cloud/sql/generated/database"
)

type NodeRespository struct {
	q *database.Queries
}

func NewNodeRepository(db database.DBTX) *NodeRespository {
	return &NodeRespository{
		q: database.New(db),
	}
}

func (r *NodeRespository) FindAll(ctx context.Context) ([]domain.Node, error) {
	dbNodes, err := r.q.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make([]domain.Node, 0, len(dbNodes))
	for _, dbNode := range dbNodes {
		nodes = append(nodes, domain.Node{
			ID:      dbNode.ID,
			NodeID:  dbNode.NodeID,
			Ip:      dbNode.Ip,
			Port:    dbNode.Port,
			Cpus:    dbNode.Cpus,
			CpuArch: dbNode.CpuArch,
		})
	}

	return nodes, nil
}

func (r *NodeRespository) CreateOrUpdate(ctx context.Context, node domain.Node) error {
	return r.q.UpsertNode(ctx, database.UpsertNodeParams{
		NodeID:  node.NodeID,
		Ip:      node.Ip,
		Port:    node.Port,
		Cpus:    node.Cpus,
		CpuArch: node.CpuArch,
	})
}

func (r *NodeRespository) Delete(ctx context.Context, nodeID string) error {
	return r.q.DeleteNode(ctx, nodeID)
}
