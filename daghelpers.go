package format

import (
	"context"

	cid "github.com/ipfs/go-cid"
)

// GetDAG will fill out all of the links of the given Node.
// It returns an array of NodePromise with the linked nodes all in the proper
// order.
func GetDAG(ctx context.Context, ds NodeGetter, root Node) []*NodePromise {
	var cids []*cid.Cid
	for _, lnk := range root.Links() {
		cids = append(cids, lnk.Cid)
	}

	return GetNodes(ctx, ds, cids)
}

// GetNodes returns an array of 'FutureNode' promises, with each corresponding
// to the key with the same index as the passed in keys
func GetNodes(ctx context.Context, ds NodeGetter, keys []*cid.Cid) []*NodePromise {

	// Early out if no work to do
	if len(keys) == 0 {
		return nil
	}

	promises := make([]*NodePromise, len(keys))
	for i := range keys {
		promises[i] = NewNodePromise(ctx)
	}

	dedupedKeys := dedupeKeys(keys)
	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		nodechan := ds.GetMany(ctx, dedupedKeys)

		for count := 0; count < len(keys); {
			select {
			case opt, ok := <-nodechan:
				if !ok {
					for _, p := range promises {
						p.Fail(ErrNotFound)
					}
					return
				}

				if opt.Err != nil {
					for _, p := range promises {
						p.Fail(opt.Err)
					}
					return
				}

				nd := opt.Node
				c := nd.Cid()
				for i, lnk_c := range keys {
					if c.Equals(lnk_c) {
						count++
						promises[i].Send(nd)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return promises
}

// Remove duplicates from a list of keys
func dedupeKeys(cids []*cid.Cid) []*cid.Cid {
	set := cid.NewSet()
	for _, c := range cids {
		set.Add(c)
	}
	return set.Keys()
}
