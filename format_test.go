package format

import (
	"context"
	"errors"
	"fmt"
	"testing"

	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

type EmptyNode struct{}

var ErrEmptyNode = errors.New("dummy node")

func (n *EmptyNode) Resolve([]string) (interface{}, []string, error) {
	return nil, nil, ErrEmptyNode
}

func (n *EmptyNode) Tree(string, int) []string {
	return nil
}

func (n *EmptyNode) ResolveLink([]string) (*Link, []string, error) {
	return nil, nil, ErrEmptyNode
}

func (n *EmptyNode) Copy() Node {
	return &EmptyNode{}
}

func (n *EmptyNode) Cid() *cid.Cid {
	id, err := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   mh.ID,
		MhLength: 0,
	}.Sum(nil)

	if err != nil {
		panic("failed to create an empty cid!")
	}
	return id
}

func (n *EmptyNode) Links() []*Link {
	return nil
}

func (n *EmptyNode) Loggable() map[string]interface{} {
	return nil
}

func (n *EmptyNode) String() string {
	return "[]"
}

func (n *EmptyNode) RawData() []byte {
	return nil
}

func (n *EmptyNode) Size() (uint64, error) {
	return 0, nil
}

func (n *EmptyNode) Stat() (*NodeStat, error) {
	return &NodeStat{}, nil
}

func TestNodeType(t *testing.T) {
	// Type assertion.
	var _ Node = &EmptyNode{}
}

func TestMakeLink(t *testing.T) {
	n := &EmptyNode{}
	l, err := MakeLink(n)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}
	expect := "z2yYDV"
	got := l.Cid.String()
	if expect != got {
		t.Errorf("cid mismatch. expected: '%s', got '%s'", expect, got)
	}
}

type SliceServ []Node

func (s SliceServ) Get(ctx context.Context, id *cid.Cid) (Node, error) {
	for _, n := range s {
		if n.Cid().Equals(id) {
			return n, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func TestLinkGetNode(t *testing.T) {
	n := &EmptyNode{}
	l, err := MakeLink(n)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	got, err := l.GetNode(context.Background(), SliceServ{n})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	if !got.Cid().Equals(n.Cid()) {
		t.Errorf("cid not equal. expected: %s, got: %s", n.Cid(), got.Cid())
		return
	}
}

func TestNodeStatString(t *testing.T) {
	ns := NodeStat{Hash: "foo", NumLinks: 1, BlockSize: 2, LinksSize: 3, DataSize: 4, CumulativeSize: 5}
	expect := "NodeStat{NumLinks: 1, BlockSize: 2, LinksSize: 3, DataSize: 4, CumulativeSize: 5}"
	got := ns.String()
	if expect != got {
		t.Errorf("string mismatch. expected: '%s'. got: '%s'", expect, got)
	}
}
