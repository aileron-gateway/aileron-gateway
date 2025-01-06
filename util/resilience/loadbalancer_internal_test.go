package resilience

import (
	"testing"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	_ LoadBalancer[Entry] = &RoundRobinLB[Entry]{}
	_ LoadBalancer[Entry] = &RandomLB[Entry]{}
	_ LoadBalancer[Entry] = &DirectHashLB[Entry]{}
	_ LoadBalancer[Entry] = &RingHashLB[Entry]{}
	_ LoadBalancer[Entry] = &MaglevLB[Entry]{}
)

type testEntry struct {
	Entry
	id     string
	weight int
	active bool
}

func (e *testEntry) ID() string {
	return e.id
}

func (e *testEntry) Weight() int {
	return e.weight
}

func (e *testEntry) Active() bool {
	return e.active
}

func TestRoundRobinLB(t *testing.T) {
	type condition struct {
		entries []*testEntry
	}

	type action struct {
		weights  []int
		entries  []*testEntry
		returned []*testEntry
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e1w2 := &testEntry{id: "e1w2", weight: 2, active: true}
	e2w0 := &testEntry{id: "e2w0", weight: 0, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}
	inactive := &testEntry{id: "inactive", weight: 1, active: false}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no entries",
			[]string{},
			[]string{},
			&condition{
				entries: nil,
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e1w1},
				returned: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2},
			},
			&action{
				weights:  []int{2},
				entries:  []*testEntry{e1w2},
				returned: []*testEntry{e1w2, e1w2, e1w2}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/contains weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e2w1},
				returned: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, e2w1},
				returned: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/contains weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2, e2w1},
			},
			&action{
				weights:  []int{2, 1},
				entries:  []*testEntry{e1w2, e2w1},
				returned: []*testEntry{e1w2, e1w2, e2w1, e1w2, e1w2, e2w1}, // Check 6 times.
			},
		),
		gen(
			"single entry/inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{inactive},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{inactive},
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
			},
			&action{
				weights:  []int{1, 1, 1},
				entries:  []*testEntry{e1w1, e2w1, e3w1},
				returned: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/contains inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, inactive},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, inactive},
				returned: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/mixed",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e1w2, e1w0, e2w0, inactive},
			},
			&action{
				weights:  []int{1, 2, 1},
				entries:  []*testEntry{e1w1, e1w2, inactive},
				returned: []*testEntry{e1w1, e1w2, e1w2, e1w1, e1w2, e1w2}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RoundRobinLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			testutil.Diff(t, tt.A().weights, lb.weights)
			testutil.Diff(t, tt.A().entries, lb.entries, cmp.AllowUnexported(testEntry{}))

			for _, r := range tt.A().returned {
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestRoundRobinLB_Get(t *testing.T) {
	type condition struct {
		entries    []*testEntry
		inactivate []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before inactivated.
		returnedAft []*testEntry // Returned entries after inactivated.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w0, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/inactivate first",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/inactivate last",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/inactivate multiple",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1, e3w1},
				inactivate: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RoundRobinLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().inactivate {
				e.(*testEntry).active = false
			}
			defer func() { // Reset active to true.
				for _, e := range tt.C().inactivate {
					e.(*testEntry).active = true
				}
			}()

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestRoundRobinLB_Remove(t *testing.T) {
	type condition struct {
		entries []*testEntry
		removes []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before remove.
		returnedAft []*testEntry // Returned entries after remove.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/remove first",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove last",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove multiple",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
				removes: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RoundRobinLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().removes {
				lb.Remove(e.ID())
			}

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestRandomLB(t *testing.T) {
	type condition struct {
		entries []*testEntry
	}

	type action struct {
		weights    []int
		entries    []*testEntry
		returned   []*testEntry
		checkRatio map[string]int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e1w2 := &testEntry{id: "e1w2", weight: 2, active: true}
	e2w0 := &testEntry{id: "e2w0", weight: 0, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}
	inactive := &testEntry{id: "inactive", weight: 1, active: false}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no entries",
			[]string{},
			[]string{},
			&condition{
				entries: nil,
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e1w1},
				returned: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2},
			},
			&action{
				weights:  []int{2},
				entries:  []*testEntry{e1w2},
				returned: []*testEntry{e1w2, e1w2, e1w2}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/contains weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e2w1},
				returned: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
			},
			&action{
				weights:    []int{1, 1},
				entries:    []*testEntry{e1w1, e2w1},
				checkRatio: map[string]int{"e1w1": 10000 * 1 / 2, "e2w1": 10000 * 1 / 2}, // Total 5000*len(entries)
			},
		),
		gen(
			"multiple entries/contains weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2, e2w1},
			},
			&action{
				weights:    []int{2, 1},
				entries:    []*testEntry{e1w2, e2w1},
				checkRatio: map[string]int{"e1w2": 10000 * 2 / 3, "e2w1": 10000 * 1 / 3}, // Total 5000*len(entries)
			},
		),
		gen(
			"single entry/inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{inactive},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{inactive},
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
			},
			&action{
				weights:    []int{1, 1, 1},
				entries:    []*testEntry{e1w1, e2w1, e3w1},
				checkRatio: map[string]int{"e1w1": 15000 * 1 / 3, "e2w1": 15000 * 1 / 3, "e3w1": 15000 * 1 / 3}, // Total 5000*len(entries)
			},
		),
		gen(
			"multiple entries/contains inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, inactive},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, inactive},
				returned: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/mixed",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e1w2, e1w0, e2w0, inactive},
			},
			&action{
				weights:    []int{1, 2, 1},
				entries:    []*testEntry{e1w1, e1w2, inactive},
				checkRatio: map[string]int{"e1w1": 25000 * 1 / 3, "e1w2": 25000 * 2 / 3}, // Total 5000*len(entries)
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RandomLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			testutil.Diff(t, tt.A().weights, lb.weights)
			testutil.Diff(t, tt.A().entries, lb.entries, cmp.AllowUnexported(testEntry{}))

			for _, r := range tt.A().returned {
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}

			// We use random algorithm.
			// So check the resulting ration instead.
			if len(tt.A().checkRatio) > 0 {
				count := map[string]int{}
				for i := 0; i < 5000*len(tt.C().entries); i++ {
					got := lb.Get(-1)
					count[got.ID()] += 1
				}
				for k, v := range tt.A().checkRatio {
					t.Log(k)
					testutil.Diff(t, float64(v), float64(count[k]), cmpopts.EquateApprox(0, 200))
				}
			}
		})
	}
}

func TestRandomLB_Remove(t *testing.T) {
	type condition struct {
		entries []*testEntry
		removes []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before remove.
		returnedAft []*testEntry // Returned entries after remove.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/remove first",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{},                       // Do not check before.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove last",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{},                       // Do not check before.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove multiple",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
				removes: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{},                                   // Do not check before.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RandomLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().removes {
				lb.Remove(e.ID())
			}

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(-1)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestDirectHashLB(t *testing.T) {
	type condition struct {
		entries []*testEntry
	}

	type action struct {
		weights  []int
		entries  []*testEntry
		returned []*testEntry
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e1w2 := &testEntry{id: "e1w2", weight: 2, active: true}
	e2w0 := &testEntry{id: "e2w0", weight: 0, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}
	inactive := &testEntry{id: "inactive", weight: 1, active: false}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no entries",
			[]string{},
			[]string{},
			&condition{
				entries: nil,
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e1w1},
				returned: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2},
			},
			&action{
				weights:  []int{2},
				entries:  []*testEntry{e1w2},
				returned: []*testEntry{e1w2, e1w2, e1w2}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/contains weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e2w1},
				returned: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, e2w1},
				returned: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/contains weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2, e2w1},
			},
			&action{
				weights:  []int{2, 1},
				entries:  []*testEntry{e1w2, e2w1},
				returned: []*testEntry{e1w2, e1w2, e2w1, e1w2, e1w2, e2w1}, // Check 6 times.
			},
		),
		gen(
			"single entry/inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{inactive},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{inactive},
				returned: []*testEntry{inactive, inactive, inactive}, // Check 3 times.
			},
		),
		gen(
			"multiple entries",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
			},
			&action{
				weights:  []int{1, 1, 1},
				entries:  []*testEntry{e1w1, e2w1, e3w1},
				returned: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/contains inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, inactive},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, inactive},
				returned: []*testEntry{e1w1, inactive, e1w1, inactive}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/mixed",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e1w2, e1w0, e2w0, inactive},
			},
			&action{
				weights:  []int{2, 1, 1}, // Entries are sorted in this algorithm.
				entries:  []*testEntry{e1w2, e1w1, inactive},
				returned: []*testEntry{e1w2, e1w2, e1w1, inactive, e1w2, e1w2, e1w1, inactive}, // Check 8 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &DirectHashLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			testutil.Diff(t, tt.A().weights, lb.weights)
			testutil.Diff(t, tt.A().entries, lb.entries, cmp.AllowUnexported(testEntry{}))

			for i, r := range tt.A().returned {
				got := lb.Get(-i) // Negative value is converted to positive.
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestDirectHashLB_Get(t *testing.T) {
	type condition struct {
		entries    []*testEntry
		inactivate []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before inactivated.
		returnedAft []*testEntry // Returned entries after inactivated.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w0, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/inactivate first",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/inactivate last",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/inactivate multiple",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1, e3w1},
				inactivate: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
				returnedAft: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &DirectHashLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().inactivate {
				e.(*testEntry).active = false
			}
			defer func() { // Reset active to true.
				for _, e := range tt.C().inactivate {
					e.(*testEntry).active = true
				}
			}()

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestDirectHashLB_Remove(t *testing.T) {
	type condition struct {
		entries []*testEntry
		removes []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before remove.
		returnedAft []*testEntry // Returned entries after remove.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/remove first",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove last",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e1w1, e2w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove multiple",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
				removes: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e3w1, e1w1, e2w1, e3w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &DirectHashLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().removes {
				lb.Remove(e.ID())
			}

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestRingHashLB(t *testing.T) {
	type condition struct {
		entries []*testEntry
	}

	type action struct {
		weights  []int
		entries  []*testEntry
		returned []*testEntry
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e1w2 := &testEntry{id: "e1w2", weight: 2, active: true}
	e2w0 := &testEntry{id: "e2w0", weight: 0, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}
	inactive := &testEntry{id: "inactive", weight: 1, active: false}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no entries",
			[]string{},
			[]string{},
			&condition{
				entries: nil,
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e1w1},
				returned: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2},
			},
			&action{
				weights:  []int{2},
				entries:  []*testEntry{e1w2},
				returned: []*testEntry{e1w2, e1w2, e1w2}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/contains weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e2w1},
				returned: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, e2w1},
				returned: []*testEntry{e1w1, e1w1, e2w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/contains weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2, e2w1},
			},
			&action{
				weights:  []int{2, 1},
				entries:  []*testEntry{e1w2, e2w1},
				returned: []*testEntry{e1w2, e1w2, e1w2, e1w2, e1w2, e2w1}, // Check 6 times.
			},
		),
		gen(
			"single entry/inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{inactive},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{inactive},
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
			},
			&action{
				weights:  []int{1, 1, 1},
				entries:  []*testEntry{e1w1, e2w1, e3w1},
				returned: []*testEntry{e1w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/contains inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, inactive},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, inactive},
				returned: []*testEntry{e1w1, e1w1, e1w1, e1w1, e1w1, e1w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/mixed",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e1w2, e1w0, e2w0, inactive},
			},
			&action{
				weights:  []int{1, 2, 1},
				entries:  []*testEntry{e1w1, e1w2, inactive},
				returned: []*testEntry{e1w2, e1w2, e1w2, e1w2, e1w2, e1w2, e1w2, e1w2}, // Check 8 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RingHashLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			testutil.Diff(t, tt.A().weights, lb.weights)
			testutil.Diff(t, tt.A().entries, lb.entries, cmp.AllowUnexported(testEntry{}))

			for i, r := range tt.A().returned {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestRingHashLB_Get(t *testing.T) {
	type condition struct {
		entries    []*testEntry
		inactivate []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before inactivated.
		returnedAft []*testEntry // Returned entries after inactivated.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w0, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/inactivate first",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e2w1, e1w1, e1w1, e2w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/inactivate last",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e2w1, e1w1, e1w1, e2w1}, // Check 6 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1, e1w1, e1w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/inactivate multiple",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1, e3w1},
				inactivate: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1},             // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RingHashLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().inactivate {
				e.(*testEntry).active = false
			}
			defer func() { // Reset active to true.
				for _, e := range tt.C().inactivate {
					e.(*testEntry).active = true
				}
			}()

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestRingHashLB_Remove(t *testing.T) {
	type condition struct {
		entries []*testEntry
		removes []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before remove.
		returnedAft []*testEntry // Returned entries after remove.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/remove first",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e2w1, e1w1}, // Check 4 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove last",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e2w1, e1w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove multiple",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
				removes: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 8 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 8 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &RingHashLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().removes {
				lb.Remove(e.ID())
			}

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestMaglevLB(t *testing.T) {
	type condition struct {
		entries []*testEntry
	}

	type action struct {
		weights  []int
		entries  []*testEntry
		returned []*testEntry
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e1w2 := &testEntry{id: "e1w2", weight: 2, active: true}
	e2w0 := &testEntry{id: "e2w0", weight: 0, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}
	e4 := &testEntry{id: "e4", weight: 100000, active: true}
	inactive := &testEntry{id: "inactive", weight: 1, active: false}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"no entries",
			[]string{},
			[]string{},
			&condition{
				entries: nil,
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e1w1},
				returned: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
			},
		),
		gen(
			"single entry/weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2},
			},
			&action{
				weights:  []int{2},
				entries:  []*testEntry{e1w2},
				returned: []*testEntry{e1w2, e1w2, e1w2}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w0},
			},
			&action{
				weights:  nil,
				entries:  nil,
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/contains weight 0",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{e2w1},
				returned: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
			},
		),
		gen(
			"multiple entries/weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, e2w1},
				returned: []*testEntry{e2w1, e2w1, e2w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/contains weight 2",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w2, e2w1},
			},
			&action{
				weights:  []int{2, 1},
				entries:  []*testEntry{e1w2, e2w1},
				returned: []*testEntry{e1w2, e2w1, e2w1, e1w2, e1w2, e1w2}, // Check 6 times.
			},
		),
		gen(
			"single entry/inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{inactive},
			},
			&action{
				weights:  []int{1},
				entries:  []*testEntry{inactive},
				returned: []*testEntry{nil, nil, nil}, // Check 3 times.
			},
		),
		gen(
			"multiple entries",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
			},
			&action{
				weights:  []int{1, 1, 1},
				entries:  []*testEntry{e1w1, e2w1, e3w1},
				returned: []*testEntry{e2w1, e2w1, e2w1, e1w1, e3w1, e3w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/contains inactive",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, inactive},
			},
			&action{
				weights:  []int{1, 1},
				entries:  []*testEntry{e1w1, inactive},
				returned: []*testEntry{e1w1, e1w1, e1w1, e1w1, e1w1, e1w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/mixed",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e1w2, e1w0, e2w0, inactive},
			},
			&action{
				weights:  []int{1, 2, 1},
				entries:  []*testEntry{e1w1, e1w2, inactive},
				returned: []*testEntry{e1w2, e1w1, e1w2, e1w1, e1w2, e1w2, e1w2, e1w2}, // Check 8 times.
			},
		),
		gen(
			"multiple entries/big weight",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e4, e2w1},
			},
			&action{
				weights:  []int{20, 1},
				entries:  []*testEntry{e4, e2w1},
				returned: []*testEntry{e4, e4, e4, e4}, // Check 4 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &MaglevLB[*testEntry]{}
			if len(tt.C().entries) > 0 {
				lb.Add(tt.C().entries...)
			}

			testutil.Diff(t, tt.A().weights, lb.weights)
			testutil.Diff(t, tt.A().entries, lb.entries, cmp.AllowUnexported(testEntry{}))

			for i, r := range tt.A().returned {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestMaglevLB_Get(t *testing.T) {
	type condition struct {
		entries    []*testEntry
		inactivate []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before inactivated.
		returnedAft []*testEntry // Returned entries after inactivated.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w0, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/inactivate first",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1, e1w1, e1w1, e1w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/inactivate last",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1},
				inactivate: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1, e1w1, e1w1, e1w1}, // Check 6 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1, e1w1, e1w1}, // Check 6 times.
			},
		),
		gen(
			"multiple entries/inactivate multiple",
			[]string{},
			[]string{},
			&condition{
				entries:    []*testEntry{e1w1, e2w1, e3w1},
				inactivate: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1, e1w1, e3w1, e3w1, e3w1, e2w1}, // Check 6 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 6 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &MaglevLB[*testEntry]{}
			if len(tt.C().entries) > 0 {
				lb.Add(tt.C().entries...)
			}

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().inactivate {
				e.(*testEntry).active = false
			}
			defer func() { // Reset active to true.
				for _, e := range tt.C().inactivate {
					e.(*testEntry).active = true
				}
			}()

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestMaglevLB_Remove(t *testing.T) {
	type condition struct {
		entries []*testEntry
		removes []Entry // Inactivate entry.
	}

	type action struct {
		returnedBef []*testEntry // Returned entries before remove.
		returnedAft []*testEntry // Returned entries after remove.
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	e1w0 := &testEntry{id: "e1w1", weight: 0, active: true}
	e1w1 := &testEntry{id: "e1w1", weight: 1, active: true}
	e2w1 := &testEntry{id: "e2w1", weight: 1, active: true}
	e3w1 := &testEntry{id: "e3w1", weight: 1, active: true}

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"single entry",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e1w1, e1w1, e1w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/one weight 1",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w0, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1}, // Check 3 times.
				returnedAft: []*testEntry{nil, nil, nil},    // Check 3 times.
			},
		),
		gen(
			"multiple entries/remove first",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e1w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1, e1w1}, // Check 4 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove last",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1},
				removes: []Entry{e2w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1, e1w1}, // Check 4 times.
				returnedAft: []*testEntry{e1w1, e1w1, e1w1, e1w1}, // Check 4 times.
			},
		),
		gen(
			"multiple entries/remove multiple",
			[]string{},
			[]string{},
			&condition{
				entries: []*testEntry{e1w1, e2w1, e3w1},
				removes: []Entry{e1w1, e3w1},
			},
			&action{
				returnedBef: []*testEntry{e2w1, e2w1, e2w1, e1w1, e3w1, e3w1, e3w1, e2w1}, // Check 8 times.
				returnedAft: []*testEntry{e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1, e2w1}, // Check 8 times.
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			lb := &MaglevLB[*testEntry]{}
			lb.Add(tt.C().entries...)

			for i, r := range tt.A().returnedBef {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
			for _, e := range tt.C().removes {
				lb.Remove(e.ID())
			}

			for i, r := range tt.A().returnedAft {
				t.Log(i)
				got := lb.Get(i)
				testutil.Diff(t, r, got, cmp.AllowUnexported(testEntry{}))
			}
		})
	}
}

func TestGenPrimeEuler(t *testing.T) {
	type condition struct {
		min int
	}

	type action struct {
		value int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"min 0",
			[]string{},
			[]string{},
			&condition{
				min: 0,
			},
			&action{
				value: 41,
			},
		),
		gen(
			"min 42",
			[]string{},
			[]string{},
			&condition{
				min: 42,
			},
			&action{
				value: 43,
			},
		),
		gen(
			"min 1681",
			[]string{},
			[]string{},
			&condition{
				min: 1680,
			},
			&action{
				value: 1847,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			value := genPrimeEuler(tt.C().min)
			testutil.Diff(t, tt.A().value, value)
		})
	}
}

func TestIsPrime(t *testing.T) {
	type condition struct {
		num []int
	}

	type action struct {
		prime bool
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"less than 2",
			[]string{},
			[]string{},
			&condition{
				num: []int{-3, -2, -1, 0, 1},
			},
			&action{
				prime: false,
			},
		),
		gen(
			"2",
			[]string{},
			[]string{},
			&condition{
				num: []int{2},
			},
			&action{
				prime: true,
			},
		),
		gen(
			"even",
			[]string{},
			[]string{},
			&condition{
				num: []int{4, 50, 100, 200, 500, 1000},
			},
			&action{
				prime: false,
			},
		),
		gen(
			"odd not prime",
			[]string{},
			[]string{},
			&condition{
				num: []int{3 * 123, 5 * 23, 7 * 99},
			},
			&action{
				prime: false,
			},
		),
		gen(
			"prime",
			[]string{},
			[]string{},
			&condition{
				num: []int{
					8803, 8807, 8819, 8821, 8831,
					8837, 8839, 8849, 8861, 8863,
					8867, 8887, 8893, 8923, 8929,
					8933, 8941, 8951, 8963, 8969,
					8971, 8999, 9001, 9007, 9011,
					9013, 9029, 9041, 9043, 9049,
					9059, 9067, 9091, 9103, 9109,
					9127, 9133, 9137, 9151, 9157,
					9161, 9173, 9181, 9187, 9199,
					9203, 9209, 9221, 9227, 9239,
					9241, 9257, 9277, 9281, 9283,
					9293, 9311, 9319, 9323, 9337,
					9341, 9343, 9349, 9371, 9377,
					9391, 9397, 9403, 9413, 9419,
					9421, 9431, 9433, 9437, 9439,
					9461, 9463, 9467, 9473, 9479,
					9491, 9497, 9511, 9521, 9533,
					9539, 9547, 9551, 9587, 9601,
					103231, 103237, 103289, 103291, 103307,
					103319, 103333, 103349, 103357, 103387,
					103391, 103393, 103399, 103409, 103421,
					103423, 103451, 103457, 103471, 103483,
					103511, 103529, 103549, 103553, 103561,
					103567, 103573, 103577, 103583, 103591,
					103613, 103619, 103643, 103651, 103657,
					103669, 103681, 103687, 103699, 103703,
					103723, 103769, 103787, 103801, 103811,
					103813, 103837, 103841, 103843, 103867,
					103889, 103903, 103913, 103919, 103951,
					103963, 103967, 103969, 103979, 103981,
					103991, 103993, 103997, 104003, 104009,
					104021, 104033, 104047, 104053, 104059,
					104087, 104089, 104107, 104113, 104119,
					104123, 104147, 104149, 104161, 104173,
					104179, 104183, 104207, 104231, 104233,
					104239, 104243, 104281, 104287, 104297,
					104309, 104311, 104323, 104327, 104347,
					104369, 104381, 104383, 104393, 104399,
					104417, 104459, 104471, 104473, 104479,
					104491, 104513, 104527, 104537, 104543,
					104549, 104551, 104561, 104579, 104593,
					104597, 104623, 104639, 104651, 104659,
					104677, 104681, 104683, 104693, 104701,
				},
			},
			&action{
				prime: true,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			for _, n := range tt.C().num {
				prime := isPrime(n)
				t.Log(n, " should prime", tt.A().prime, " actual", prime)
				testutil.Diff(t, tt.A().prime, prime)
			}
		})
	}
}
