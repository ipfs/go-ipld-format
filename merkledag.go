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

// TODO: This name kind of sucks.
// NodeResolver?
// NodeService?
// Just Resolver?
type NodeGetter interface {
	Get(context.Context, *cid.Cid) (Node, error)
}

// DAGService is an IPFS Merkle DAG service.
type DAGService interface {
	NodeGetter

	Add(Node) (*cid.Cid, error)
	Remove(Node) error

	// TODO: This is returning them in-order?? Why not just use []NodePromise?
	// Maybe add a couple of helpers for getting them in-order and as-available?
	// GetDAG returns, in order, all the single leve child
	// nodes of the passed in node.
	GetMany(context.Context, []*cid.Cid) <-chan *NodeOption

	AddMany([]Node) ([]*cid.Cid, error)

	LinkService
}

// TODO: Replace this? I'm really not convinced this interface pulls its weight.
//
// Instead, we could add an `Offline()` function to `NodeGetter` that returns an
// offline `NodeGetter` and then define the following function:
//
// ```
// func GetLinks(ctx context.Context, ng NodeGetter, c *cid.Cid) ([]*Link, error) {
// 	if c.Type() == cid.Raw {
// 		return nil, nil
// 	}
// 	node, err := ng.Get(ctx, c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return node.Links(), nil
// }
// ```
//
// Why *not* do this? We might decide to store a light-weight DAG of links
// without actually storing the data. I don't really find that to be a
// convincing argument.
type LinkService interface {
	// GetLinks return all links for a node.  The complete node does not
	// necessarily have to exist locally, or at all.  For example, raw
	// leaves cannot possibly have links so there is no need to look
	// at the node.
	// TODO: These *really* should be Cids, not Links
	GetLinks(context.Context, *cid.Cid) ([]*Link, error)

	GetOfflineLinkService() LinkService
}
