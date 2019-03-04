package main

import (
	"testing"
)

func compareSlice(a, b []string) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// I want to be able to filter a list and only keep elements beginning with b.
func TestFindsAllStringsContainingB(t *testing.T) {
	listOfFiles := []string{
		"thing",
		"other",
		"things",
		"bar",
		"foo",
		"broken",
	}

	resultB := Filter(listOfFiles, "b")

	expectedResultForB := []string{
		"bar",
		"broken",
	}

	if !compareSlice(resultB, expectedResultForB) {
		t.Errorf("did not get the list of expected files back for filter b, got %v", resultB)
	}

	resultT := Filter(listOfFiles, "t")

	expectedResultForT := []string{
		"thing",
		"other",
		"things",
	}

	if !compareSlice(resultT, expectedResultForT) {
		t.Errorf("did not get the list of expected files back for filter t, got: %v", resultT)
	}
}
