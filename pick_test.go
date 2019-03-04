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

	result := Filter(listOfFiles, "b")

	expectedResult := []string{
		"bar",
		"broken",
	}

	if !compareSlice(result, expectedResult) {
		t.Errorf("did not get the list of expected files back")
	}
}
