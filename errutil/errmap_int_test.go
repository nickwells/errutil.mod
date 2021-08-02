package errutil

import (
	"errors"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestErrListDiffs(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		name     string
		list1    []error
		list2    []error
		expDiffs []string
	}{
		{
			ID: testhelper.MkID("empty lists - no diff"),
		},
		{
			ID:    testhelper.MkID("non-empty lists - no diff"),
			list1: []error{errors.New("err1")},
			list2: []error{errors.New("err1")},
		},
		{
			ID: testhelper.MkID(
				"non-empty lists, empty name - different lengths"),
			list1:    []error{errors.New("err1")},
			list2:    []error{errors.New("err1"), errors.New("err2")},
			expDiffs: []string{"error counts differ: 1 != 2"},
		},
		{
			ID: testhelper.MkID(
				"non-empty lists, non-empty name - different lengths"),
			name:     "xxx",
			list1:    []error{errors.New("err1")},
			list2:    []error{errors.New("err1"), errors.New("err2")},
			expDiffs: []string{`"xxx": error counts differ: 1 != 2`},
		},
		{
			ID: testhelper.MkID(
				"non-empty lists, empty name - different errors"),
			list1: []error{errors.New("err1"), errors.New("err2")},
			list2: []error{errors.New("err1"), errors.New("different")},
			expDiffs: []string{`error[1]:
	    "err2"
	 != "different"`},
		},
		{
			ID: testhelper.MkID(
				"non-empty lists, non-empty name - different errors"),
			name:  "xxx",
			list1: []error{errors.New("err1"), errors.New("err2")},
			list2: []error{errors.New("err1"), errors.New("different")},
			expDiffs: []string{`"xxx": error[1]:
	    "err2"
	 != "different"`},
		},
		{
			ID: testhelper.MkID(
				"non-empty lists, non-empty name - multiple different errors"),
			name: "xxx",
			list1: []error{
				errors.New("err1"),
				errors.New("err2"),
				errors.New("err3"),
			},
			list2: []error{
				errors.New("err1"),
				errors.New("different"),
				errors.New("different"),
			},
			expDiffs: []string{
				`"xxx": error[1]:
	    "err2"
	 != "different"`,
				`"xxx": error[2]:
	    "err3"
	 != "different"`,
			},
		},
		{
			ID: testhelper.MkID(
				"non-empty lists, empty name - multiple different errors"),
			list1: []error{
				errors.New("err1"),
				errors.New("err2"),
				errors.New("err3"),
			},
			list2: []error{
				errors.New("err1"),
				errors.New("different"),
				errors.New("different"),
			},
			expDiffs: []string{
				`error[1]:
	    "err2"
	 != "different"`,
				`error[2]:
	    "err3"
	 != "different"`,
			},
		},
	}

	for _, tc := range testCases {
		diffs := errListDiffs(tc.name, tc.list1, tc.list2)
		testhelper.DiffStringSlice(t, tc.IDStr(), "diffs", diffs, tc.expDiffs)
	}
}
