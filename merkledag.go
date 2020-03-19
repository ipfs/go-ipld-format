package format

import (
	"context"
	"fmt"

	cid "github.com/ipfs/go-cid"
)

// ErrNotFound is used to signal when a Node could not be found. The specific
// meaning will depend on the DAGService implementation, which may be trying
// to read nodes locally but also, trying to find them remotely.
var ErrNotFound = ErrNotFoundCid{}

// ErrNotFoundCid can be use to provide specific CID information in a NotFound
// error.
type ErrNotFoundCid struct {
	c cid.Cid
}

// Error implements the error interface and returns a human-readable
// message for this error.
func (e ErrNotFoundCid) Error() string {
	if e.c == cid.Undef {
		return "ipld: node not found"
	}

	return fmt.Sprintf("ipld: %s not found", e.c)
}

// Is allows to check whether any error is of this ErrNotFoundCid type.
// Do not use this directly, but rather errors.Is(yourError, ErrNotFound).
func (e ErrNotFoundCid) Is(err error) bool {
	switch err.(type) {
	case ErrNotFoundCid:
		return true
	default:
		return false
	}
}

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
	Get(context.Context, cid.Cid) (Node, error)

	// GetMany returns a channel of NodeOptions given a set of CIDs.
	GetMany(context.Context, []cid.Cid) <-chan *NodeOption
}

// NodeAdder adds nodes to a DAG.
type NodeAdder interface {
	// Add adds a node to this DAG.
	Add(context.Context, Node) error

	// AddMany adds many nodes to this DAG.
	//
	// Consider using the Batch NodeAdder (`NewBatch`) if you make
	// extensive use of this function.
	AddMany(context.Context, []Node) error
}

// NodeGetters can optionally implement this interface to make finding linked
// objects faster.
type LinkGetter interface {
	NodeGetter

	// TODO(ipfs/go-ipld-format#9): This should return []cid.Cid

	// GetLinks returns the children of the node refered to by the given
	// CID.
	GetLinks(ctx context.Context, nd cid.Cid) ([]*Link, error)
}

// DAGService is an IPFS Merkle DAG service.
type DAGService interface {
	NodeGetter
	NodeAdder

	// Remove removes a node from this DAG.
	//
	// Remove returns no error if the requested node is not present in this DAG.
	Remove(context.Context, cid.Cid) error

	// RemoveMany removes many nodes from this DAG.
	//
	// It returns success even if the nodes were not present in the DAG.
	RemoveMany(context.Context, []cid.Cid) error
}
