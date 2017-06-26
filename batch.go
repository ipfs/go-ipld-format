package format

import (
	cid "github.com/ipfs/go-cid"
)

func Batching(ds DAGService) *Batch {
	return &Batch{
		ds:      ds,
		MaxSize: 8 << 20,

		// By default, only batch up to 128 nodes at a time.
		// The current implementation of flatfs opens this many file
		// descriptors at the same time for the optimized batch write.
		MaxBlocks: 128,
	}
}

type Batch struct {
	ds DAGService

	// TODO: try to re-use memory.
	nodes     []Node
	size      int
	MaxSize   int
	MaxBlocks int
}

func (t *Batch) Add(nd Node) (*cid.Cid, error) {
	t.nodes = append(t.nodes, nd)
	t.size += len(nd.RawData())
	if t.size > t.MaxSize || len(t.nodes) > t.MaxBlocks {
		return nd.Cid(), t.Commit()
	}
	return nd.Cid(), nil
}

func (t *Batch) Commit() error {
	_, err := t.ds.AddMany(t.nodes)
	t.nodes = nil
	t.size = 0
	return err
}
