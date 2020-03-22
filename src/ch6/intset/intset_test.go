package main

import (
	"testing"
)

func TestHas(t *testing.T) {
	var set IntSet
	set.Add(8)
	set.Add(144)

	if !set.Has(8) {
		t.Errorf("Set does not contain %d %v\n", 8, set.String())
	}
	if !set.Has(144) {
		t.Errorf("Set does not contain %d %v\n", 144, set.String())
	}
	if set.Has(142) {
		t.Errorf("Set contains %d %v but %d was never added\n", 142, set.String(), 142)
	}
}

func TestUnionWith(t *testing.T) {
	var s, p IntSet
	var sVals []int = []int{1, 3}
	var pVals []int = []int{2, 4, 65}
	for i := 0; i < len(sVals); i++ {
		s.Add(sVals[i])
	}
	for i := 0; i < len(pVals); i++ {
		p.Add(pVals[i])
	}

	s.UnionWith(&p)

	sVals = append(sVals, pVals...)
	for i := 0; i < len(sVals); i++ {
		if !s.Has(sVals[i]) {
			t.Errorf("Set does not contain %d %v\n", sVals[i], s.String())
		}
	}
}
