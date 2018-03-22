package format

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	errs "github.com/pkg/errors"
)

// ErrNotFound indicates that a CID was not found
type ErrNotFound struct {
	Cid *cid.Cid
}

func (e ErrNotFound) Error() string {
	if e.Cid == nil {
		return fmt.Sprintf("CID not found")
	}
	return fmt.Sprintf("CID not found: %s", e.Cid)
}

// IsNotFound returns true if the cause of the error is ErrNotFound
func IsNotFound(e error) bool {
	_, ok := errs.Cause(e).(ErrNotFound)
	return ok
}
