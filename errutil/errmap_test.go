package errutil_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestErrMap(t *testing.T) {
	type catErr struct {
		cat string
		err error
	}
	type catSummary struct {
		cat     string
		summary string
	}
	testCases := []struct {
		testhelper.ID
		catErrs    []catErr
		expTotErrs int
		expTotCats int
		expSummary string
		expCS      []catSummary
		expReport  string
	}{
		{
			ID:         testhelper.MkID("no errors"),
			expSummary: "no errors were found",
			expCS: []catSummary{
				{cat: "nonesuch", summary: "nonesuch - no errors"},
			},
			expReport: `testname: no errors were found
`,
		},
		{
			ID: testhelper.MkID("one error"),
			catErrs: []catErr{
				{cat: "tiger", err: errors.New("cat too scary")},
			},
			expTotErrs: 1,
			expTotCats: 1,
			expSummary: "an error was found",
			expCS: []catSummary{
				{cat: "nonesuch", summary: "nonesuch - no errors"},
				{cat: "tiger", summary: "tiger:"},
			},
			expReport: `testname: an error was found
      tiger:
            cat too scary
`,
		},
		{
			ID: testhelper.MkID("two errors, one category"),
			catErrs: []catErr{
				{cat: "tiger", err: errors.New("cat too scary")},
				{cat: "tiger", err: errors.New("still too scary")},
			},
			expTotErrs: 2,
			expTotCats: 1,
			expSummary: "2 errors were found",
			expCS: []catSummary{
				{cat: "nonesuch", summary: "nonesuch - no errors"},
				{cat: "tiger", summary: "tiger - 2 errors:"},
			},
			expReport: `testname: 2 errors were found
      tiger - 2 errors:
            1 : cat too scary
            2 : still too scary
`,
		},
		{
			ID: testhelper.MkID("3 errors, 2 categories"),
			catErrs: []catErr{
				{cat: "tiger", err: errors.New("cat too scary")},
				{cat: "tiger", err: errors.New("still too scary")},
				{cat: "tigger", err: errors.New("not a real cat")},
			},
			expTotErrs: 3,
			expTotCats: 2,
			expSummary: "3 errors were found in 2 categories",
			expCS: []catSummary{
				{cat: "nonesuch", summary: "nonesuch - no errors"},
				{cat: "tiger", summary: "tiger - 2 errors:"},
				{cat: "tigger", summary: "tigger:"},
			},
			expReport: `testname: 3 errors were found in 2 categories
      tiger - 2 errors:
            1 : cat too scary
            2 : still too scary
      tigger:
            not a real cat
`,
		},
		{
			ID: testhelper.MkID("11 errors, 2 categories"),
			catErrs: []catErr{
				{cat: "tiger", err: errors.New("cat too scary")},
				{cat: "tiger", err: errors.New("still too scary")},
				{cat: "tiger", err: errors.New("very scary")},
				{cat: "tiger", err: errors.New("really scary")},
				{cat: "tiger", err: errors.New("frightening")},
				{cat: "tiger", err: errors.New("terrifying")},
				{cat: "tiger", err: errors.New("scary")},
				{cat: "tiger", err: errors.New("scary")},
				{cat: "tiger", err: errors.New("scary")},
				{cat: "tiger", err: errors.New("scary")},
				{cat: "tigger", err: errors.New("not a real cat")},
			},
			expTotErrs: 11,
			expTotCats: 2,
			expSummary: "11 errors were found in 2 categories",
			expCS: []catSummary{
				{cat: "nonesuch", summary: "nonesuch - no errors"},
				{cat: "tiger", summary: "tiger - 10 errors:"},
				{cat: "tigger", summary: "tigger:"},
			},
			expReport: `testname: 11 errors were found in 2 categories
      tiger - 10 errors:
             1 : cat too scary
             2 : still too scary
             3 : very scary
             4 : really scary
             5 : frightening
             6 : terrifying
             7 : scary
             8 : scary
             9 : scary
            10 : scary
      tigger:
            not a real cat
`,
		},
	}

	for _, tc := range testCases {
		em := errutil.NewErrMap()
		for _, ce := range tc.catErrs {
			em.AddError(ce.cat, ce.err)
		}
		totErrs, totCats := em.CountErrors()
		testhelper.DiffInt(t, tc.IDStr(), "tot errs", totErrs, tc.expTotErrs)
		testhelper.DiffInt(t, tc.IDStr(), "tot cats", totCats, tc.expTotCats)
		s := em.Summary()
		testhelper.DiffString(t, tc.IDStr(), "summary", s, tc.expSummary)
		for i, ecs := range tc.expCS {
			s := em.CategorySummary(ecs.cat)
			testhelper.DiffString(t, tc.IDStr(),
				fmt.Sprintf("category summary: %d (%q)", i, ecs.cat),
				s, ecs.summary)
		}
		var b bytes.Buffer
		em.Report(&b, "testname")
		testhelper.DiffString(t, tc.IDStr(), "report", b.String(), tc.expReport)
	}
}

func TestMatches(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		em1 errutil.ErrMap
		em2 errutil.ErrMap
	}{
		{
			ID: testhelper.MkID("matches - empty maps"),
		},
		{
			ID:  testhelper.MkID("matches - one cat, one error"),
			em1: errutil.ErrMap{"cat1": []error{errors.New("an error")}},
			em2: errutil.ErrMap{
				"cat1": []error{errors.New("an error")},
			},
		},
		{
			ID: testhelper.MkID("matches fail - different numbers of cats"),
			ExpErr: testhelper.MkExpErr(
				`the category names differ:
	"cat2" in other, not this`),
			em1: errutil.ErrMap{
				"cat1": []error{errors.New("an error")},
			},
			em2: errutil.ErrMap{
				"cat1": []error{errors.New("an error")},
				"cat2": []error{errors.New("an error")},
			},
		},
		{
			ID: testhelper.MkID("matches fail - different cats"),
			ExpErr: testhelper.MkExpErr(
				`the category names differ:
	"cat1" in other, not this
	"cat3" in this, not other`),
			em1: errutil.ErrMap{
				"cat2": []error{errors.New("an error")},
				"cat3": []error{errors.New("an error")},
			},
			em2: errutil.ErrMap{
				"cat1": []error{errors.New("an error")},
				"cat2": []error{errors.New("an error")},
			},
		},
	}

	for _, tc := range testCases {
		err := tc.em1.Matches(tc.em2)
		testhelper.CheckExpErr(t, err, tc)
	}
}
