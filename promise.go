package format

import (
	"context"
	"sync"
)

// TODO: I renamed this to NodePromise because:
// 1. NodeGetter is a naming conflict.
// 2. It's a promise...

// TODO: Should this even be an interface? It seems like a simple struct would
// suffice.

// NodePromise provides a promise like interface for a dag Node
// the first call to Get will block until the Node is received
// from its internal channels, subsequent calls will return the
// cached node.
type NodePromise interface {
	Get(context.Context) (Node, error)
	Fail(err error)
	Send(Node)
}

func newNodePromise(ctx context.Context) NodePromise {
	return &nodePromise{
		recv: make(chan Node, 1),
		ctx:  ctx,
		err:  make(chan error, 1),
	}
}

type nodePromise struct {
	cache Node
	clk   sync.Mutex
	recv  chan Node
	ctx   context.Context
	err   chan error
}

func (np *nodePromise) Fail(err error) {
	np.clk.Lock()
	v := np.cache
	np.clk.Unlock()

	// if promise has a value, don't fail it
	if v != nil {
		return
	}

	np.err <- err
}

func (np *nodePromise) Send(nd Node) {
	var already bool
	np.clk.Lock()
	if np.cache != nil {
		already = true
	}
	np.cache = nd
	np.clk.Unlock()

	if already {
		panic("sending twice to the same promise is an error!")
	}

	np.recv <- nd
}

func (np *nodePromise) Get(ctx context.Context) (Node, error) {
	np.clk.Lock()
	c := np.cache
	np.clk.Unlock()
	if c != nil {
		return c, nil
	}

	select {
	case nd := <-np.recv:
		return nd, nil
	case <-np.ctx.Done():
		return nil, np.ctx.Err()
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-np.err:
		return nil, err
	}
}
