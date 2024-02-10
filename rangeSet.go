package main

import "errors"

type Range struct {
	StartTimestamp int64
	EndTimestamp   int64
}

type RangeSet struct {
	Ranges []Range
}

func (rangeSet *RangeSet) Add(r Range) error {
	if len(rangeSet.Ranges) > 0 {
		existing := &rangeSet.Ranges[len(rangeSet.Ranges)-1]
		if r.StartTimestamp < existing.StartTimestamp {
			return errors.New("Tried to insert into range set out of order")
		} else if r.StartTimestamp <= existing.EndTimestamp {
			existing.EndTimestamp = max(existing.EndTimestamp, r.EndTimestamp)
			return nil
		}
	}
	rangeSet.Ranges = append(rangeSet.Ranges, r)
	return nil
}

func NewRangeSet() RangeSet {
	return RangeSet{}
}
