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
	// Get retrieves nodes by CID. Depending on the NodeGetter
	// implementation, this may involve fetching the Node from a remote
	// machine; consider setting a deadline in the context.
	Get(context.Context, *cid.Cid) (Node, error)

	// GetMany returns a channel of NodeOptions given a set of CIDs.
	GetMany(context.Context, []*cid.Cid) <-chan *NodeOption
}

// NodeGetters can optionally implement this interface to make finding linked
// objects faster.
type LinkGetter interface {
	NodeGetter

	// TODO(ipfs/go-ipld-format#9): This should return []*cid.Cid

	// GetLinks returns the children of the node refered to by the given
	// CID.
	GetLinks(ctx context.Context, nd *cid.Cid) ([]*Link, error)
}

// DAGService is an IPFS Merkle DAG service.
type DAGService interface {
	NodeGetter

	// Add adds a node to this DAG.
	Add(Node) (*cid.Cid, error)

	// Remove removes a node from this DAG.
	//
	// If the node is not in this DAG, Remove returns ErrNotFound.
	Remove(*cid.Cid) error

	// AddMany adds many nodes to this DAG.
	//
	// Consider using NewBatch instead of calling this directly if you need
	// to add an unbounded number of nodes to avoid buffering too much.
	AddMany([]Node) ([]*cid.Cid, error)
}
