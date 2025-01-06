package cron_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/cron"
	"github.com/aileron-gateway/aileron-gateway/kernel/testutil"
)

func ExampleCrontab() {
	ct, err := cron.Parse("59 59 23 31 * *") // Day 31 at 23:59:59
	if err != nil {
		panic(err)
	}

	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ct.WithTestTimer(func() time.Time { return now })
	ct.WithLocation(time.UTC)

	for i := 0; i < 10; i++ {
		now = now.Add(10 * 24 * time.Hour) // Move 10 days.
		fmt.Println(now.Format(time.DateTime), ct.Next().Format(time.DateTime))
	}

	// Output:
	// 2000-01-11 00:00:00 2000-01-31 23:59:59
	// 2000-01-21 00:00:00 2000-01-31 23:59:59
	// 2000-01-31 00:00:00 2000-01-31 23:59:59
	// 2000-02-10 00:00:00 2000-03-31 23:59:59
	// 2000-02-20 00:00:00 2000-03-31 23:59:59
	// 2000-03-01 00:00:00 2000-03-31 23:59:59
	// 2000-03-11 00:00:00 2000-03-31 23:59:59
	// 2000-03-21 00:00:00 2000-03-31 23:59:59
	// 2000-03-31 00:00:00 2000-03-31 23:59:59
	// 2000-04-10 00:00:00 2000-05-31 23:59:59
}

func ExampleCronJob() {
	ct, err := cron.Parse("* * * * * *") // Run every seconds.
	if err != nil {
		panic(err)
	}

	counter := 0
	job := ct.NewJob(func() {
		counter += 1
		fmt.Printf("%d ", counter)
	})
	go job.Start()
	defer job.Stop()

	time.Sleep(5 * time.Second)
	// Output:
	// 1 2 3 4 5
}

func ExampleCronJob_panicJob() {
	ct, err := cron.Parse("* * * * * *") // Run every seconds.
	if err != nil {
		panic(err)
	}

	job := ct.NewJob(func() {
		fmt.Println("I'm almost panic.")
		panic(errors.New("Now I'm in panic!!"))
	})
	go job.Start()
	defer job.Stop()

	time.Sleep(2 * time.Second)
	// Output:
	// I'm almost panic.
	// I'm almost panic.
}

func TestCrontab_Next(t *testing.T) {
	type condition struct {
		exp   string
		times map[time.Time]time.Time
	}

	type action struct {
	}

	tb := testutil.NewTableBuilder[*condition, *action]()
	tb.Name(t.Name())
	table := tb.Build()

	gen := testutil.NewCase[*condition, *action]
	testCases := []*testutil.Case[*condition, *action]{
		gen(
			"6 field, increment second by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC):    time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),  // increment second 0 > 1
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),  // increment second 0 > 1
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):      time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC),  // increment second 1 > 2
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):     time.Date(2000, 1, 1, 0, 0, 30, 0, time.UTC), // increment second 29 > 30
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):     time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC), // increment second 58 > 59
					time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC):     time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),  // increment minute 0 > 1
					time.Date(2000, 1, 1, 0, 29, 59, 0, time.UTC):    time.Date(2000, 1, 1, 0, 30, 0, 0, time.UTC), // increment minute 29 > 30
					time.Date(2000, 1, 1, 0, 58, 59, 0, time.UTC):    time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC), // increment minute 58 > 59
					time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC):    time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),  // increment hour 0 > 1
					time.Date(2000, 1, 1, 22, 59, 59, 0, time.UTC):   time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC), // increment hour 22 > 23
					time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 23, 59, 59, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 23, 59, 59, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 23, 59, 59, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"6 field, increment minute by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 * * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):    time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),  // increment minute 0 > 1
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):     time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),  // increment minute 0 > 1
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):    time.Date(2000, 1, 1, 0, 30, 0, 0, time.UTC), // increment minute 29 > 30
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):    time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC), // increment minute 58 > 59
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):    time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),  // increment hour 0 > 1
					time.Date(2000, 1, 1, 22, 59, 0, 0, time.UTC):   time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC), // increment hour 22 > 23
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 23, 59, 0, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 23, 59, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 23, 59, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"5 field, increment minute by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "* * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):    time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),  // increment minute 0 > 1
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):     time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),  // increment minute 0 > 1
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):    time.Date(2000, 1, 1, 0, 30, 0, 0, time.UTC), // increment minute 29 > 30
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):    time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC), // increment minute 58 > 59
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):    time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),  // increment hour 0 > 1
					time.Date(2000, 1, 1, 22, 59, 0, 0, time.UTC):   time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC), // increment hour 22 > 23
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 23, 59, 0, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 23, 59, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 23, 59, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"6 field, increment hour by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):    time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),  // increment hour 0 > 1
					time.Date(2000, 1, 1, 22, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC), // increment hour 22 > 23
					time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 23, 0, 0, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 23, 0, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 23, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 23, 0, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 23, 0, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"5 field, increment hour by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):    time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),  // increment hour 0 > 1
					time.Date(2000, 1, 1, 22, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC), // increment hour 22 > 23
					time.Date(2000, 1, 1, 23, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 23, 0, 0, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 23, 0, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 23, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 23, 0, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 23, 0, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"6 field, increment day by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 22, 0, 0, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"5 field, increment day by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 22, 0, 0, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),  // increment day 1 > 2
					time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC):  time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC), // increment day 15 > 16
					time.Date(2000, 1, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC), // increment day 30 > 31
					time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),  // increment month 1 > 2
					time.Date(2000, 4, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC),  // increment month 4 > 5
					time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),  // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"6 field, increment month by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 1 * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 22, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC), // increment month 1 > 2
					time.Date(2000, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC), // increment month 4 > 5
					time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"5 field, increment month by 1",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 1 * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 59, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 29, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 58, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 22, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC), // increment month 1 > 2
					time.Date(2000, 4, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC), // increment month 4 > 5
					time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), // increment year 2000 > 2001
				},
			},
			&action{},
		),
		gen(
			"every 2 seconds",
			[]string{},
			[]string{},
			&condition{
				exp: "*/2 * * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC):   time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 0, 30, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 57, 0, time.UTC):  time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 1, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC):   time.Date(2000, 1, 1, 0, 1, 4, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 30, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 57, 0, time.UTC):  time.Date(2000, 1, 1, 0, 1, 58, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 59, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 59, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC):   time.Date(2000, 1, 1, 0, 59, 4, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 59, 30, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 57, 0, time.UTC):  time.Date(2000, 1, 1, 0, 59, 58, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 58, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 23, 59, 0, 999, time.UTC): time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC):   time.Date(2000, 1, 1, 23, 59, 4, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 29, 0, time.UTC):  time.Date(2000, 1, 1, 23, 59, 30, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 57, 0, time.UTC):  time.Date(2000, 1, 1, 23, 59, 58, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 58, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 31, 23, 59, 0, 999, time.UTC): time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC):   time.Date(2000, 1, 31, 23, 59, 4, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 29, 0, time.UTC):  time.Date(2000, 1, 31, 23, 59, 30, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 57, 0, time.UTC):  time.Date(2000, 1, 31, 23, 59, 58, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 58, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 12, 31, 23, 59, 0, 999, time.UTC): time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC):   time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC):   time.Date(2000, 12, 31, 23, 59, 4, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 29, 0, time.UTC):  time.Date(2000, 12, 31, 23, 59, 30, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 57, 0, time.UTC):  time.Date(2000, 12, 31, 23, 59, 58, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 58, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"every 2 minutes",
			[]string{},
			[]string{},
			&condition{
				exp: "0 */2 * * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC):   time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 57, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 1, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC):   time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 57, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 59, 0, time.UTC):  time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 2, 0, 999, time.UTC): time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 2, 0, 0, time.UTC):   time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 2, 2, 0, time.UTC):   time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 2, 29, 0, time.UTC):  time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 2, 57, 0, time.UTC):  time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 2, 58, 0, time.UTC):  time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 2, 59, 0, time.UTC):  time.Date(2000, 1, 1, 0, 4, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 59, 0, 999, time.UTC): time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC):   time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 29, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 57, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 58, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC):  time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 23, 59, 0, 999, time.UTC): time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 29, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 57, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 58, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 31, 23, 59, 0, 999, time.UTC): time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 29, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 57, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 58, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 12, 31, 23, 59, 0, 999, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 29, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 57, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 58, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"every 2 hours",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 */2 * * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC):   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 57, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 1, 0, 999, time.UTC): time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC):   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC):   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 29, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 57, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 58, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 59, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 59, 0, 999, time.UTC): time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC):   time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 29, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 57, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 58, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC):  time.Date(2000, 1, 1, 2, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 23, 59, 0, 999, time.UTC): time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC):   time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 29, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 57, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 58, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC):  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 31, 23, 59, 0, 999, time.UTC): time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 29, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 57, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 58, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 12, 31, 23, 59, 0, 999, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 29, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 57, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 58, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"every 2 days",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 */2 * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 57, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 1, 0, 999, time.UTC): time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 29, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 57, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 58, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 59, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 59, 0, 999, time.UTC): time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 29, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 57, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 58, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 23, 59, 0, 999, time.UTC): time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC):   time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 29, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 57, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 58, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC):  time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 31, 23, 59, 0, 999, time.UTC): time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC):   time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 29, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 57, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 58, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC):  time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 12, 31, 23, 59, 0, 999, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 29, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 57, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 58, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"every 2 months",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 1 */2 *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 999, time.UTC): time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 29, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 57, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 58, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 1, 0, 999, time.UTC): time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 2, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 29, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 57, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 58, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 1, 59, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 0, 59, 0, 999, time.UTC): time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 0, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 2, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 29, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 57, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 58, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 23, 59, 0, 999, time.UTC): time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 0, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 2, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 29, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 57, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 58, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 31, 23, 59, 0, 999, time.UTC): time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 0, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 2, 0, time.UTC):   time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 29, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 57, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 58, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC):  time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 12, 31, 23, 59, 0, 999, time.UTC): time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 2, 0, time.UTC):   time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 29, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 57, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 58, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC):  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"day 31",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 31 * *",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 28, 0, 0, 0, 0, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 29, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 29, 0, 0, 0, 0, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 7, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 7, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 8, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 8, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 9, 1, 0, 0, 0, 0, time.UTC):   time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 9, 29, 0, 0, 0, 0, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 9, 30, 0, 0, 0, 0, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 10, 30, 0, 0, 0, 0, time.UTC): time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 11, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 11, 30, 0, 0, 0, 0, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 11, 31, 0, 0, 0, 0, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 30, 0, 0, 0, 0, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC): time.Date(2001, 1, 31, 0, 0, 0, 0, time.UTC),

					time.Date(2000, 1, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 999, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 28, 23, 59, 59, 999, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 29, 23, 59, 59, 999, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 31, 23, 59, 59, 999, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 29, 23, 59, 59, 999, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 5, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 31, 23, 59, 59, 999, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 29, 23, 59, 59, 999, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 7, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 7, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 7, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 7, 31, 23, 59, 59, 999, time.UTC):  time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 8, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 8, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 8, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 8, 31, 23, 59, 59, 999, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 9, 1, 23, 59, 59, 999, time.UTC):   time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 9, 29, 23, 59, 59, 999, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 9, 30, 23, 59, 59, 999, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 10, 1, 23, 59, 59, 999, time.UTC):  time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 10, 30, 23, 59, 59, 999, time.UTC): time.Date(2000, 10, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 10, 31, 23, 59, 59, 999, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 11, 1, 23, 59, 59, 999, time.UTC):  time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 11, 30, 23, 59, 59, 999, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 11, 31, 23, 59, 59, 999, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 1, 23, 59, 59, 999, time.UTC):  time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 30, 23, 59, 59, 999, time.UTC): time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 12, 31, 23, 59, 59, 999, time.UTC): time.Date(2001, 1, 31, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"Friday",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 * * FRI",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 6, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 8, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 9, 0, 0, 0, 0, time.UTC):       time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 11, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 12, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 13, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 16, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 17, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 18, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 19, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 20, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 22, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 23, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 24, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 25, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 26, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 27, 0, 0, 0, 0, time.UTC):      time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC):      time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 29, 0, 0, 0, 0, time.UTC):      time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 30, 0, 0, 0, 0, time.UTC):      time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC):      time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 1, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 2, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 3, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 4, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 5, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 6, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 7, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 8, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 9, 23, 59, 59, 999, time.UTC):  time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 10, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 11, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 12, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 13, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 14, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 14, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 15, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 16, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 17, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 18, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 19, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 20, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 21, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 22, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 23, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 24, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 25, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 26, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 27, 23, 59, 59, 999, time.UTC): time.Date(2000, 1, 28, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 28, 23, 59, 59, 999, time.UTC): time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 29, 23, 59, 59, 999, time.UTC): time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 30, 23, 59, 59, 999, time.UTC): time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 1, 31, 23, 59, 59, 999, time.UTC): time.Date(2000, 2, 4, 0, 0, 0, 0, time.UTC),
				},
			},
			&action{},
		),
		gen(
			"Friday, 30 JUN",
			[]string{},
			[]string{},
			&condition{
				exp: "0 0 0 30 JUN FRI",
				times: map[time.Time]time.Time{
					time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 4, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 5, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 1, 0, 0, 0, 0, time.UTC):  time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC),
					time.Date(2000, 6, 30, 0, 0, 0, 0, time.UTC): time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
					time.Date(2000, 7, 1, 0, 0, 0, 0, time.UTC):  time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
					time.Date(2000, 8, 1, 0, 0, 0, 0, time.UTC):  time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
					time.Date(2000, 9, 1, 0, 0, 0, 0, time.UTC):  time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
					time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC): time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
					time.Date(2000, 11, 1, 0, 0, 0, 0, time.UTC): time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
					time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC): time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC), // Next 2006
				},
			},
			&action{},
		),
	}

	testutil.Register(table, testCases...)

	for _, tt := range table.Entries() {
		tt := tt
		t.Run(tt.Name(), func(t *testing.T) {
			ct, err := cron.Parse(tt.C().exp)
			testutil.Diff(t, nil, err)

			now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			ct.WithTestTimer(func() time.Time { return now })
			ct.WithLocation(time.UTC)

			for k, v := range tt.C().times {
				now = k
				next := ct.Next()
				t.Log(k.Format(time.DateTime), v.Format(time.DateTime))
				testutil.Diff(t, v.Format(time.DateTime), next.Format(time.DateTime))
			}
		})
	}
}
