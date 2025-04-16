// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cron

import (
	"strconv"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func TestNextTime(t *testing.T) {
	type condition struct {
		target uint64
		now    int
		max    int
	}

	type action struct {
		out int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"invalid target 0",
			[]string{},
			[]string{},
			&condition{
				target: 0,
				now:    0,
				max:    9,
			},
			&action{
				out: 0,
			},
		),
		gen(
			"now 0, next 1",
			[]string{},
			[]string{},
			&condition{
				target: 0b_00011,
				now:    0,
				max:    9,
			},
			&action{
				out: 1,
			},
		),
		gen(
			"now 0, next 2",
			[]string{},
			[]string{},
			&condition{
				target: 0b_00101,
				now:    0,
				max:    9,
			},
			&action{
				out: 2,
			},
		),
		gen(
			"now 0, next 9",
			[]string{},
			[]string{},
			&condition{
				target: 0b_10000_00001,
				now:    0,
				max:    9,
			},
			&action{
				out: 9,
			},
		),
		gen(
			"now 0, next 0",
			[]string{},
			[]string{},
			&condition{
				target: 0b_00001_00000_00001,
				now:    0,
				max:    9,
			},
			&action{
				out: 0,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := nextTime(tt.C().target, tt.C().now, tt.C().max)
			testutil.Diff(t, tt.A().out, out)
		})
	}
}

func TestCronJob(t *testing.T) {
	type condition struct {
		exp         string
		jobDuration time.Duration
		timer       func() time.Time
	}

	type action struct {
		count []int
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"run once",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * * *",
			},
			&action{
				count: []int{1},
			},
		),
		gen(
			"run 3 times",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * * *",
			},
			&action{
				count: []int{1, 2, 3},
			},
		),
		gen(
			"no duplicate run",
			[]string{},
			[]string{},
			&condition{
				exp:         "* * * * * *",
				jobDuration: 5 * time.Second,
			},
			&action{
				count: []int{1, 1, 1},
			},
		),
		gen(
			"should calibrate",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 " + strconv.Itoa(time.Now().Local().Day()) + " * *",
			},
			&action{
				count: []int{0, 0, 0},
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ct, err := Parse(tt.C().exp)
			testutil.Diff(t, nil, err)
			if tt.C().timer != nil {
				ct.WithTestTimer(tt.C().timer)
			}

			counter := 0
			job := ct.NewJob(func() {
				counter += 1
				time.Sleep(tt.C().jobDuration)
			})

			notice := make(chan struct{})
			job.WithTestWaiter(func(_ time.Duration) time.Duration {
				notice <- struct{}{}
				return time.Nanosecond
			})

			go job.Start()
			go job.Start() // Duplicate run is ignored.

			for i := 0; i < len(tt.A().count); i++ {
				<-notice
				time.Sleep(100 * time.Millisecond) // Wait the job stated if necessary.
				testutil.Diff(t, tt.A().count[i], counter)
			}

			testutil.Diff(t, false, job.stop == nil)
			job.Stop()
			testutil.Diff(t, true, job.stop == nil)
		})
	}
}

func TestMax(t *testing.T) {
	type condition struct {
		a time.Duration
		b time.Duration
	}

	type action struct {
		out time.Duration
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"0<a<b",
			[]string{},
			[]string{},
			&condition{
				a: time.Second,
				b: time.Hour,
			},
			&action{
				out: time.Hour,
			},
		),
		gen(
			"a<0<b",
			[]string{},
			[]string{},
			&condition{
				a: -time.Second,
				b: time.Hour,
			},
			&action{
				out: time.Hour,
			},
		),
		gen(
			"a<b<0",
			[]string{},
			[]string{},
			&condition{
				a: -time.Second,
				b: -time.Hour,
			},
			&action{
				out: -time.Second,
			},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			out := max(tt.C().a, tt.C().b)
			testutil.Diff(t, tt.A().out, out)
		})
	}
}
