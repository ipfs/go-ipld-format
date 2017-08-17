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

	// TODO(ipfs/go-ipfs#4009): Remove this method after fixing.

	// OfflineNodeGetter returns an version of this NodeGetter that will
	// make no network requests.
	OfflineNodeGetter() NodeGetter
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

// GetLinks returns the CIDs of the children of the given node. Prefer this
// method over looking up the node itself and calling `Links()` on it as this
// method may be able to use a link cache.
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
