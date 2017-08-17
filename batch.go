package format

import (
	cid "github.com/ipfs/go-cid"
)

// NewBatch returns a node buffer (Batch) that buffers nodes internally and
// commits them to the underlying DAGService in batches. Use this if you intend
// to add a lot of nodes all at once.
func NewBatch(ds DAGService) *Batch {
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

// Add a node to this batch of nodes, potentially committing the set of batched
// nodes to the underlying DAGService.
func (t *Batch) Add(nd Node) (*cid.Cid, error) {
	t.nodes = append(t.nodes, nd)
	t.size += len(nd.RawData())
	if t.size > t.MaxSize || len(t.nodes) > t.MaxBlocks {
		return nd.Cid(), t.Commit()
	}
	return nd.Cid(), nil
}

// Commit commits the buffered of nodes to the underlying DAGService.
// Make sure to call this after you're done adding nodes to the batch to ensure
// that they're actually added to the DAGService.
func (t *Batch) Commit() error {
	_, err := t.ds.AddMany(t.nodes)
	t.nodes = nil
	t.size = 0
	return err
}
