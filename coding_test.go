package format

import (
	"errors"
	"testing"

	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

func TestDecode(t *testing.T) {
	decoder := func(b blocks.Block) (Node, error) {
		node := &EmptyNode{}
		if b.RawData() != nil || !b.Cid().Equals(node.Cid()) {
			return nil, errors.New("can only decode empty blocks")
		}
		return node, nil
	}

	id, err := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   mh.IDENTITY,
		MhLength: 0,
	}.Sum(nil)

	if err != nil {
		t.Fatalf("failed to create cid: %s", err)
	}

	block, err := blocks.NewBlockWithCid(nil, id)
	if err != nil {
		t.Fatalf("failed to create empty block: %s", err)
	}
	node, err := Decode(block, decoder)
	if err != nil {
		t.Fatalf("failed to decode empty node: %s", err)
	}
	if !node.Cid().Equals(id) {
		t.Fatalf("empty node doesn't have the right cid")
	}

	if _, ok := node.(*EmptyNode); !ok {
		t.Fatalf("empty node doesn't have the right type")
	}

}

func TestRegistryDecode(t *testing.T) {
	decoder := func(b blocks.Block) (Node, error) {
		node := &EmptyNode{}
		if b.RawData() != nil || !b.Cid().Equals(node.Cid()) {
			return nil, errors.New("can only decode empty blocks")
		}
		return node, nil
	}

	id, err := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   mh.IDENTITY,
		MhLength: 0,
	}.Sum(nil)

	if err != nil {
		t.Fatalf("failed to create cid: %s", err)
	}

	block, err := blocks.NewBlockWithCid(nil, id)
	if err != nil {
		t.Fatalf("failed to create empty block: %s", err)
	}

	reg := Registry{}
	_, err = reg.Decode(block)
	if err == nil || err.Error() != "unrecognized object type: 85" {
		t.Fatalf("expected error, got %v", err)
	}
	reg.Register(cid.Raw, decoder)
	node, err := reg.Decode(block)
	if err != nil {
		t.Fatalf("failed to decode empty node: %s", err)
	}

	if !node.Cid().Equals(id) {
		t.Fatalf("empty node doesn't have the right cid")
	}

	if _, ok := node.(*EmptyNode); !ok {
		t.Fatalf("empty node doesn't have the right type")
	}

}
