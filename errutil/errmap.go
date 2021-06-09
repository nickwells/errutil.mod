package errutil

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/nickwells/twrap.mod/twrap"
)

const (
	categoryIndent = 6
	errorIndent    = 12
)

// ErrMap is a type that maps a string to a list of errors. It is useful for
// the case where you want to return more than one error and want to group
// them in some way for reporting. Each map entry represents a category for
// which some errors have been found.
type ErrMap map[string][]error

// NewErrMap returns a new instance of an ErrMap
func NewErrMap() *ErrMap {
	return &ErrMap{}
}

// AddError adds the error to the category slice
func (em *ErrMap) AddError(cat string, err error) {
	(*em)[cat] = append((*em)[cat], err)
}

// HasErrors returns true if there are any errors in the map, false
// otherwise.
func (em ErrMap) HasErrors() bool {
	return len(em) > 0
}

// CountErrors counts the total number of errors and the number of categories
// (in that order)
func (em ErrMap) CountErrors() (int, int) {
	totErrs := 0
	for _, errs := range em {
		totErrs += len(errs)
	}
	return totErrs, len(em)
}

// Summary returns a summary description of the ErrMap
func (em ErrMap) Summary() string {
	totErrs, categories := em.CountErrors()

	summary := "no errors were found"
	if totErrs == 1 {
		summary = "an error was found"
	} else if totErrs > 1 {
		summary = fmt.Sprintf("%d errors were found", totErrs)
		if categories > 1 {
			summary += fmt.Sprintf(" in %d categories", categories)
		}
	}
	return summary
}

// CategorySummary returns a summary description of the errors in the ErrMap
// for the given category
func (em ErrMap) CategorySummary(cat string) string {
	switch len(em[cat]) {
	case 0:
		return cat + " - no errors"
	case 1:
		return cat + ":"
	}
	return fmt.Sprintf("%s - %d errors:", cat, len(em[cat]))
}

// Keys returns the keys to the ErrMap (the categories)
func (em ErrMap) Keys() []string {
	cats := make([]string, 0, len(em))
	for cat := range em {
		cats = append(cats, cat)
	}
	return cats
}

// Report writes the error map out to the writer
func (em ErrMap) Report(w io.Writer, name string) {
	twc := twrap.NewTWConfOrPanic(twrap.SetWriter(w))

	if name == "" {
		twc.Wrap(em.Summary(), 0)
	} else {
		twc.WrapPrefixed(name+": ", em.Summary(), 0)
	}

	cats := em.Keys()
	sort.Strings(cats)

	for _, cat := range cats {
		em.reportErrors(twc, cat)
	}
}

// reportErrors reports the errors for the category
func (em ErrMap) reportErrors(twc *twrap.TWConf, cat string) {
	twc.Wrap(em.CategorySummary(cat), categoryIndent)

	errs := em[cat]
	digitCount := 1
	if len(errs) >= 10 {
		digitCount = 2
	}
	if len(errs) >= 100 {
		digitCount = 3
	}
	prefix := ""

	for i, e := range errs {
		if len(errs) > 1 {
			prefix = fmt.Sprintf("%*d : ", digitCount, i+1)
		}
		twc.WrapPrefixed(prefix, e.Error(), errorIndent)
	}
}

// Matches returns an error if the two ErrMaps differ, nil otherwise
func (em ErrMap) Matches(other ErrMap) error {
	differingCats := []string{}
	for k := range em {
		if _, ok := other[k]; !ok {
			differingCats = append(differingCats,
				fmt.Sprintf("%q in this, not other", k))
		}
	}
	for k := range other {
		if _, ok := em[k]; !ok {
			differingCats = append(differingCats,
				fmt.Sprintf("%q in other, not this", k))
		}
	}
	if len(differingCats) > 0 {
		sort.Strings(differingCats)
		return fmt.Errorf("the category names differ:\n\t%s",
			strings.Join(differingCats, "\n\t"))
	}

	errDiffs := []string{}
	for cat, errs := range em {
		errDiffs = append(errDiffs, errListDiffs(cat, errs, other[cat])...)
	}
	if len(errDiffs) > 0 {
		sort.Strings(errDiffs)
		return fmt.Errorf("the error details differ: %s",
			strings.Join(errDiffs, ", "))
	}
	return nil
}

// errListDiffs returns a list of differences between the two lists of errors
func errListDiffs(name string, list1, list2 []error) []string {
	pfx := ""
	if name != "" {
		pfx = fmt.Sprintf("%q: ", name)
	}
	errDiffs := []string{}
	if len(list1) != len(list2) {
		return append(errDiffs,
			fmt.Sprintf("%serror counts differ: %d != %d",
				pfx, len(list1), len(list2)))
	}
	for i, err := range list1 {
		if err.Error() != list2[i].Error() {
			errDiffs = append(errDiffs,
				fmt.Sprintf("%serror[%d]: %q != %q",
					pfx, i, err, list2[i]))
		}
	}
	return errDiffs
}
