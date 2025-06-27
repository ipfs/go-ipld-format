package format

import (
	"context"
)

// NewNodePromise constructs a NodePromise with the given context. Canceling the
// context will immediately cancel the NodePromise.
func NewNodePromise(ctx context.Context) *NodePromise {
	return &NodePromise{
		done: make(chan struct{}),
		ctx:  ctx,
	}
}

// NodePromise provides a promise like interface for a dag Node
// the first call to Get will block until the Node is received
// from its internal channels, subsequent calls will return the
// cached node.
//
// Thread Safety: This is multiple-consumer/single-producer safe.
type NodePromise struct {
	value Node
	err   error
	done  chan struct{}

	ctx context.Context
}

// Fail fails this promise.
//
// Once a promise has been failed or fulfilled, further attempts to fail it will
// be silently dropped.
func (np *NodePromise) Fail(err error) {
	if np.err != nil || np.value != nil {
		// Already filled.
		return
	}
	np.err = err
	close(np.done)
}

// Send fulfills this promise.
//
// Once a promise has been fulfilled or failed, calling this function will
// panic.
func (np *NodePromise) Send(nd Node) {
	// if promise has a value, don't fail it
	if np.err != nil || np.value != nil {
		panic("already filled")
	}
	np.value = nd
	close(np.done)
}

// Poll returns the result of the promise if ready but doesn't block.
//
// Returns nil, nil if not ready.
func (np *NodePromise) Poll() (Node, error) {
	select {
	case <-np.done:
		return np.value, np.err
	case <-np.ctx.Done():
		return nil, np.ctx.Err()
	default:
		return nil, nil
	}
}

// Get the value of this promise.
//
// This function is safe to call concurrently from any number of goroutines.
func (np *NodePromise) Get(ctx context.Context) (Node, error) {
	select {
	case <-np.done:
		return np.value, np.err
	case <-np.ctx.Done():
		return nil, np.ctx.Err()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
