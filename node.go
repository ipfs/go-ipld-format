package node

import (
	"context"
	"fmt"

	cid "github.com/ipfs/go-cid"
)

type Node interface {
	Resolve(path []string) (*Link, []string, error)
	Tree() []string
	Cid() *cid.Cid

	Links() []*Link

	//
	Stat() (*NodeStat, error)
	Size() (uint64, error)

	// RawData marshals the node and returns the marshaled bytes
	RawData() []byte

	String() string
	Loggable() map[string]interface{}
}

type NodeGetter interface {
	Get(context.Context, *cid.Cid) (Node, error)
}

// Link represents an IPFS Merkle DAG Link between Nodes.
type Link struct {
	// utf string name. should be unique per object
	Name string // utf8

	// cumulative size of target object
	Size uint64

	// multihash of the target object
	Cid *cid.Cid
}

// NodeStat is a statistics object for a Node. Mostly sizes.
type NodeStat struct {
	Hash           string
	NumLinks       int // number of links in link table
	BlockSize      int // size of the raw, encoded data
	LinksSize      int // size of the links segment
	DataSize       int // size of the data segment
	CumulativeSize int // cumulative size of object and its references
}

func (ns NodeStat) String() string {
	f := "NodeStat{NumLinks: %d, BlockSize: %d, LinksSize: %d, DataSize: %d, CumulativeSize: %d}"
	return fmt.Sprintf(f, ns.NumLinks, ns.BlockSize, ns.LinksSize, ns.DataSize, ns.CumulativeSize)
}

// MakeLink creates a link to the given node
func MakeLink(n Node) (*Link, error) {
	s, err := n.Size()
	if err != nil {
		return nil, err
	}

	return &Link{
		Size: s,
		Cid:  n.Cid(),
	}, nil
}

// GetNode returns the MDAG Node that this link points to
func (l *Link) GetNode(ctx context.Context, serv NodeGetter) (Node, error) {
	return serv.Get(ctx, l.Cid)
}
