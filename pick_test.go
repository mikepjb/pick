package main

import (
	"testing"
)

func compareSlice(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	if len(a) != len(b) {
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
	"thing.ext",
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
		"thing.ext",
	}

	if !compareSlice(resultT, expectedResultForT) {
		t.Errorf("did not get the list of expected files back for filter t, got: %v", resultT)
	}
}

func TestFindsFuzzySearches(t *testing.T) {
	resultTS := Filter(listOfFiles, "ts")

	expectedResultForTS := []string{"things"}

	if !compareSlice(resultTS, expectedResultForTS) {
		t.Errorf("did not get the list of expected files back for filter ts, got: %v", resultTS)
	}
}

func TestReplacesFullStop(t *testing.T) {
	resultFS := Filter(listOfFiles, "tx")

	expectedResultForFS := []string{"thing.ext"}

	if !compareSlice(resultFS, expectedResultForFS) {
		t.Errorf("did not get the list of expected files back for filter t., got: %v", resultFS)
	}
}

func TestFiltersFromList(t *testing.T) {
	filterList := []string{"node_modules/"}

	projectList := []string{"node_modules/a_thing.txt", "this_thing.txt"}

	filteredList := filterByList(projectList, filterList)
	expectedList := []string{"this_thing.txt"}

	if !compareSlice(filteredList, expectedList) {
		t.Errorf("list not filtered correctly, got: %v\n", filteredList)
	}
}
