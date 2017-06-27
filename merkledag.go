package format

import (
	"context"
	"fmt"

	cid "github.com/ipfs/go-cid"
)

var ErrNotFound = fmt.Errorf("merkledag: not found")

// Either a node or an error.
type NodeOption struct {
	Node Node
	Err  error
}

// The basic Node resolution service.
type NodeGetter interface {
	Get(context.Context, *cid.Cid) (Node, error)

	// TODO(ipfs/go-ipfs#4009): Remove this method after fixing.
	OfflineNodeGetter() NodeGetter
}

// NodeGetters can optionally implement this interface to make finding linked
// objects faster.
type LinkGetter interface {
	NodeGetter
	// TODO(ipfs/go-ipld-format#9): This should return []*cid.Cid
	GetLinks(ctx context.Context, nd *cid.Cid) ([]*Link, error)
}

func GetLinks(ctx context.Context, ng NodeGetter, c *cid.Cid) ([]*Link, error) {
	if c.Type() == cid.Raw {
		return nil, nil
	}
	if gl, ok := ng.(LinkGetter); ok {
		return gl.GetLinks(ctx, c)
	}
	node, err := ng.Get(ctx, c)
	if err != nil {
		return nil, err
	}
	return node.Links(), nil
}

// DAGService is an IPFS Merkle DAG service.
type DAGService interface {
	NodeGetter

	Add(Node) (*cid.Cid, error)

	// TODO(ipfs/go-ipfs#4010): Change this to take a CID.
	// This will require a fair amount of refactoring.
	Remove(Node) error

	// TODO: Consider using []NodePromise and providing helper functions
	// that take []NodePromise and return channels that yield nodes both
	// in-order and as-ready.
	GetMany(context.Context, []*cid.Cid) <-chan *NodeOption

	AddMany([]Node) ([]*cid.Cid, error)
}
