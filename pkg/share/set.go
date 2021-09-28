// SPDX-License-Identifier: Apache-2.0

package share

// Set implemetation
type Set struct {
	m map[string]bool
}

// NewSet inits a set
func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]bool)

	return s
}

// Add add a item to the set
func (s *Set) Add(value string) {
	s.m[value] = true
}

// Remove removes a item from the set
func (s *Set) Remove(value string) {
	delete(s.m, value)
}

// Contains verifies whether set contains a key
func (s *Set) Contains(value string) bool {
	_, c := s.m[value]

	return c
}
