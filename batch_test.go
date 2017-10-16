package format

import (
	"context"
	"sync"
	"testing"

	cid "github.com/ipfs/go-cid"
)

// Test dag
type testDag struct {
	mu    sync.Mutex
	nodes map[string]Node
}

func newTestDag() *testDag {
	return &testDag{nodes: make(map[string]Node)}
}

func (d *testDag) Get(ctx context.Context, cid *cid.Cid) (Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if n, ok := d.nodes[cid.KeyString()]; ok {
		return n, nil
	}
	return nil, ErrNotFound
}

func (d *testDag) GetMany(ctx context.Context, cids []*cid.Cid) <-chan *NodeOption {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make(chan *NodeOption, len(cids))
	for _, c := range cids {
		if n, ok := d.nodes[c.KeyString()]; ok {
			out <- &NodeOption{Node: n}
		} else {
			out <- &NodeOption{Err: ErrNotFound}
		}
	}
	return out
}

func (d *testDag) Add(node Node) (*cid.Cid, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	c := node.Cid()
	d.nodes[c.KeyString()] = node
	return c, nil
}

func (d *testDag) AddMany(nodes []Node) ([]*cid.Cid, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	cids := make([]*cid.Cid, len(nodes))
	for i, n := range nodes {
		c := n.Cid()
		d.nodes[c.KeyString()] = n
		cids[i] = c
	}
	return cids, nil
}

func (d *testDag) Remove(c *cid.Cid) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	key := c.KeyString()
	if _, exists := d.nodes[key]; !exists {
		return ErrNotFound
	}
	delete(d.nodes, key)
	return nil
}

var _ DAGService = new(testDag)

func TestBatch(t *testing.T) {
	d := newTestDag()
	b := NewBatch(d)
	for i := 0; i < 1000; i++ {
		// It would be great if we could use *many* different nodes here
		// but we can't add any dependencies and I don't feel like adding
		// any more testing code.
		if _, err := b.Add(new(EmptyNode)); err != nil {
			t.Fatal(err)
		}
	}
	if err := b.Commit(); err != nil {
		t.Fatal(err)
	}

	n, err := d.Get(context.Background(), new(EmptyNode).Cid())
	if err != nil {
		t.Fatal(err)
	}
	switch n.(type) {
	case *EmptyNode:
	default:
		t.Fatal("expected the node to exist in the dag")
	}

	if len(d.nodes) != 1 {
		t.Fatal("should have one node")
	}
}
