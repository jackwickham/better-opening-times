package main

import (
	"reflect"
	"testing"
)

func TestRangeSetSingleElement(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 1,
		EndTimestamp:   2,
	})
	if err != nil {
		t.Error("Error while adding: ", err)
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 1,
		EndTimestamp:   2,
	}}) {
		t.Errorf("expected [[1, 2]], was %v", rs.Ranges)
	}
}

func TestRangeSetMergesOverlappingElements(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 1,
		EndTimestamp:   4,
	})
	if err != nil {
		t.Error("Error while adding first element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 3,
		EndTimestamp:   6,
	})
	if err != nil {
		t.Error("Error while adding second element: ", err)
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 1,
		EndTimestamp:   6,
	}}) {
		t.Errorf("expected [[1, 6]], was %v", rs.Ranges)
	}
}

func TestRangeSetMergesAdjacentElements(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 1,
		EndTimestamp:   3,
	})
	if err != nil {
		t.Error("Error while adding first element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 3,
		EndTimestamp:   6,
	})
	if err != nil {
		t.Error("Error while adding second element: ", err)
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 1,
		EndTimestamp:   6,
	}}) {
		t.Errorf("expected [[1, 6]], was %v", rs.Ranges)
	}
}

func TestRangeSetAppendsNonAdjacentElements(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 1,
		EndTimestamp:   4,
	})
	if err != nil {
		t.Error("Error while adding first element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 5,
		EndTimestamp:   6,
	})
	if err != nil {
		t.Error("Error while adding second element: ", err)
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 1,
		EndTimestamp:   4,
	}, {
		StartTimestamp: 5,
		EndTimestamp:   6,
	}}) {
		t.Errorf("expected [[1, 4], [5, 6]], was %v", rs.Ranges)
	}
}

func TestRangeSetUpdatesLastElement(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 1,
		EndTimestamp:   3,
	})
	if err != nil {
		t.Error("Error while adding first element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 5,
		EndTimestamp:   7,
	})
	if err != nil {
		t.Error("Error while adding second element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 6,
		EndTimestamp:   9,
	})
	if err != nil {
		t.Error("Error while adding third element: ", err)
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 1,
		EndTimestamp:   3,
	}, {
		StartTimestamp: 5,
		EndTimestamp:   9,
	}}) {
		t.Errorf("expected [[1, 3], [5, 9]], was %v", rs.Ranges)
	}
}

func TestRangeSetMergesContainedElements(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 1,
		EndTimestamp:   6,
	})
	if err != nil {
		t.Error("Error while adding first element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 3,
		EndTimestamp:   5,
	})
	if err != nil {
		t.Error("Error while adding second element: ", err)
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 1,
		EndTimestamp:   6,
	}}) {
		t.Errorf("expected [[1, 6]], was %v", rs.Ranges)
	}
}

func TestRangeSetErrorsIfNotSorted(t *testing.T) {
	rs := NewRangeSet()
	err := rs.Add(Range{
		StartTimestamp: 4,
		EndTimestamp:   6,
	})
	if err != nil {
		t.Error("Error while adding first element: ", err)
	}
	err = rs.Add(Range{
		StartTimestamp: 3,
		EndTimestamp:   5,
	})
	if err == nil {
		t.Error("Did not receive error when adding unordered element")
	}
	if !reflect.DeepEqual(rs.Ranges, []Range{{
		StartTimestamp: 4,
		EndTimestamp:   6,
	}}) {
		t.Errorf("expected [[4, 6]], was %v", rs.Ranges)
	}
}
