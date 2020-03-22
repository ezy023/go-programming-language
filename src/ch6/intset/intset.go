// An integer set implemented with a bit vector
package main

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word := x / 64
	bit := uint(x % 64)
	if len(s.words) < word {
		return false
	}
	return s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word := x / 64
	bit := uint(x % 64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Exercises
// func (s *IntSet) Len() int      {}
// func (s *IntSet) Remove(x int)  {}
// func (s *IntSet) Clear()        {}
// func (s *IntSet) Copy() *IntSet {}
// func (s *IntSet) AddAll(vals ...int) {}
// func (s *IntSet) IntersectWith(t *IntSet) {}
// func (s *IntSet) DifferenceWith(t *IntSet) {}
// func (s *IntSet) SymmetricDifference(t *IntSet) {}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
