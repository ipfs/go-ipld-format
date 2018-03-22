package format

import (
	"fmt"
	"testing"

	errs "github.com/pkg/errors"
)

func TestIsNotFound(t *testing.T) {
	err1 := ErrNotFound{}
	if !IsNotFound(err1) {
		t.Errorf("IsNotFound not true when it should be")
	}
	err2 := errs.Wrap(err1, "wrapped")
	if !IsNotFound(err2) {
		t.Errorf("IsNotFound not true on wrapped error when it should be")
	}
	err3 := fmt.Errorf("CID not found")
	if IsNotFound(err3) {
		t.Errorf("IsNotFound true when it should not be")
	}
}
