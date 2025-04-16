// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package resilience_test

import (
	"fmt"

	"github.com/aileron-gateway/aileron-gateway/util/resilience"
)

type ExampleEntry struct {
	resilience.Entry
	id     string
	weight int
	active bool
}

func (e *ExampleEntry) ID() string {
	return e.id
}

func (e *ExampleEntry) Weight() int {
	return e.weight
}

func (e *ExampleEntry) Active() bool {
	return e.active
}

func ExampleRoundRobinLB() {
	e1 := &ExampleEntry{id: "e1", weight: 1, active: true}
	e2 := &ExampleEntry{id: "e2", weight: 1, active: true}
	e3 := &ExampleEntry{id: "e3", weight: 1, active: true}
	lb := &resilience.RoundRobinLB[*ExampleEntry]{}
	lb.Add(e1, e2, e3)

	count := map[string]int{}
	for i := 0; i < 300; i++ {
		entry := lb.Get(-1)
		count[entry.ID()] += 1
	}

	fmt.Println(count)
	// Output:
	// map[e1:100 e2:100 e3:100]
}

func ExampleRandomLB() {
	e1 := &ExampleEntry{id: "e1", weight: 1, active: true}
	e2 := &ExampleEntry{id: "e2", weight: 1, active: true}
	e3 := &ExampleEntry{id: "e3", weight: 1, active: true}
	lb := &resilience.RandomLB[*ExampleEntry]{}
	lb.Add(e1, e2, e3)

	count := map[string]int{}
	for i := 0; i < 300; i++ {
		entry := lb.Get(-1)
		count[entry.ID()] += 1
	}

	fmt.Println(count)
	// Output:
	// map[e1:100 e2:100 e3:100]
}

func ExampleDirectHashLB() {
	e1 := &ExampleEntry{id: "e1", weight: 1, active: true}
	e2 := &ExampleEntry{id: "e2", weight: 1, active: true}
	e3 := &ExampleEntry{id: "e3", weight: 1, active: true}
	lb := &resilience.DirectHashLB[*ExampleEntry]{}
	lb.Add(e1, e2, e3)

	count := map[string]int{}
	for i := 0; i < 300; i++ {
		entry := lb.Get(i) // Input should be a int generated from hash.
		count[entry.ID()] += 1
	}

	fmt.Println(count)
	// Output:
	// map[e1:100 e2:100 e3:100]
}

func ExampleRingHashLB() {
	e1 := &ExampleEntry{id: "e1", weight: 1, active: true}
	e2 := &ExampleEntry{id: "e2", weight: 1, active: true}
	e3 := &ExampleEntry{id: "e3", weight: 1, active: true}
	lb := &resilience.RingHashLB[*ExampleEntry]{}
	lb.Add(e1, e2, e3)

	count := map[string]int{}
	for i := 0; i < 300; i++ {
		entry := lb.Get(i) // Input should be a int generated from hash.
		count[entry.ID()] += 1
	}

	fmt.Println(count)
	// Output:
	// map[e1:113 e2:90 e3:97]
}

func ExampleMaglevLB() {
	e1 := &ExampleEntry{id: "e1", weight: 1, active: true}
	e2 := &ExampleEntry{id: "e2", weight: 1, active: true}
	e3 := &ExampleEntry{id: "e3", weight: 1, active: true}
	lb := &resilience.MaglevLB[*ExampleEntry]{}
	lb.Add(e1, e2, e3)

	count := map[string]int{}
	for i := 0; i < 300; i++ {
		entry := lb.Get(i) // Input should be a int generated from hash.
		count[entry.ID()] += 1
	}

	fmt.Println(count)
	// Output:
	// map[e1:103 e2:104 e3:93]
}
