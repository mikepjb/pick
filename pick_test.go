package main

import (
	"fmt"
	"testing"
)

func compareSlice(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var listOfFiles = []string{
	"thing",
	"other",
	"things",
	"bar",
	"foo",
	"broken",
}

// I want to be able to filter a list and only keep elements beginning with b.
func TestFindsSequentialSearches(t *testing.T) {
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

func TestFindsFuzzySearches(t *testing.T) {
	resultTS := Filter(listOfFiles, "ts")
	fmt.Println(resultTS)

	expectedResultForTS := []string{"things"}

	if !compareSlice(resultTS, expectedResultForTS) {
		t.Errorf("did not get the list of expected files back for filter ts, got: %v", resultTS)
	}
}
