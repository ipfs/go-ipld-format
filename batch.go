package format

import (
	"runtime"

	cid "github.com/ipfs/go-cid"
)

// ParallelBatchCommits is the number of batch commits that can be in-flight before blocking.
// TODO(ipfs/go-ipfs#4299): Experiment with multiple datastores, storage
// devices, and CPUs to find the right value/formula.
var ParallelBatchCommits = runtime.NumCPU() * 2

// NewBatch returns a node buffer (Batch) that buffers nodes internally and
// commits them to the underlying DAGService in batches. Use this if you intend
// to add a lot of nodes all at once.
func NewBatch(ds DAGService) *Batch {
	return &Batch{
		ds:            ds,
		commitResults: make(chan error, ParallelBatchCommits),
		MaxSize:       8 << 20,

		// By default, only batch up to 128 nodes at a time.
		// The current implementation of flatfs opens this many file
		// descriptors at the same time for the optimized batch write.
		MaxNodes: 128,
	}
}

// Batch is a buffer for batching adds to a dag.
type Batch struct {
	ds DAGService

	activeCommits int
	commitError   error
	commitResults chan error

	nodes []Node
	size  int

	MaxSize  int
	MaxNodes int
}

func (t *Batch) processResults() {
	for t.activeCommits > 0 && t.commitError == nil {
		select {
		case err := <-t.commitResults:
			t.activeCommits--
			if err != nil {
				t.commitError = err
			}
		default:
			return
		}
	}
}

func (t *Batch) asyncCommit() {
	numBlocks := len(t.nodes)
	if numBlocks == 0 || t.commitError != nil {
		return
	}
	if t.activeCommits >= ParallelBatchCommits {
		err := <-t.commitResults
		t.activeCommits--

		if err != nil {
			t.commitError = err
			return
		}
	}
	go func(b []Node) {
		_, err := t.ds.AddMany(b)
		t.commitResults <- err
	}(t.nodes)

	t.activeCommits++
	t.nodes = make([]Node, 0, numBlocks)
	t.size = 0

	return
}

// Add adds a node to the batch and commits the batch if necessary.
func (t *Batch) Add(nd Node) (*cid.Cid, error) {
	// Not strictly necessary but allows us to catch errors early.
	t.processResults()
	if t.commitError != nil {
		return nil, t.commitError
	}

	t.nodes = append(t.nodes, nd)
	t.size += len(nd.RawData())
	if t.size > t.MaxSize || len(t.nodes) > t.MaxNodes {
		t.asyncCommit()
	}
	return nd.Cid(), t.commitError
}

// Commit commits batched nodes.
func (t *Batch) Commit() error {
	t.asyncCommit()
	for t.activeCommits > 0 && t.commitError == nil {
		err := <-t.commitResults
		t.activeCommits--
		if err != nil {
			t.commitError = err
		}
	}

	return t.commitError
}
