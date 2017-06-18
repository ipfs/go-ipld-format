package format

import (
	"fmt"

	blocks "github.com/ipfs/go-block-format"
)

// DecodeBlock functions decode blocks into nodes.
type DecodeBlockFunc func(block blocks.Block) (Node, error)

// Map from codec types to decoder functions
type BlockDecoder map[uint64]DecodeBlockFunc

// A default set of block decoders.
//
// You SHOULD populate this map from `init` functions in packages that support
// decoding various IPLD formats. You MUST NOT modify this map once `main` has
// been called.
var DefaultBlockDecoder BlockDecoder = map[uint64]DecodeBlockFunc{}

func (b BlockDecoder) Decode(block blocks.Block) (Node, error) {
	// Short-circuit by cast if we already have a Node.
	if node, ok := block.(Node); ok {
		return node, nil
	}

	ty := block.Cid().Type()
	if decoder, ok := b[ty]; ok {
		return decoder(block)
	} else {
		// TODO: get the *long* name for this format
		return nil, fmt.Errorf("unrecognized object type: %d", ty)
	}
}

// Decode the given block using the default block decoder.
func Decode(block blocks.Block) (Node, error) {
	return DefaultBlockDecoder.Decode(block)
}
