package format

import (
	"context"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"testing"
)

func TestCopy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	from := newTestDag()
	if err := from.Add(ctx, new(EmptyNode)); err != nil {
		t.Fatal(err)
	}
	to := newTestDag()
	id, err := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   mh.ID,
		MhLength: 0,
	}.Sum(nil)
	err = Copy(ctx, from, to, id)
	if err != nil {
		t.Error(err)
	}
}
