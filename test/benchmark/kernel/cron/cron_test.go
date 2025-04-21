// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package cron_test

import (
	"testing"
	"time"

	"github.com/aileron-gateway/aileron-gateway/kernel/cron"
)

func loopTimer(d time.Duration) func() time.Time {
	num := 10000
	times := make([]time.Time, num)
	now := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < num; i++ {
		times[i] = now.Add(d)
	}
	count := 0
	return func() time.Time {
		if count > num {
			count = 0
		}
		return times[count]
	}
}

func BenchmarkMinutely(b *testing.B) {
	ct, _ := cron.Parse("* * * * *")
	ct.WithTestTimer(loopTimer((60 + 30) * time.Second))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Next()
	}
}

func BenchmarkHourly(b *testing.B) {
	ct, _ := cron.Parse("@hourly")
	ct.WithTestTimer(loopTimer((60 + 30) * time.Minute))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Next()
	}
}

func BenchmarkDaily(b *testing.B) {
	ct, _ := cron.Parse("@daily")
	ct.WithTestTimer(loopTimer((24 + 12) * time.Hour))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Next()
	}
}

func BenchmarkWeekly(b *testing.B) {
	ct, _ := cron.Parse("@weekly")
	ct.WithTestTimer(loopTimer((7 + 3) * 24 * time.Hour))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Next()
	}
}

func BenchmarkMonthly(b *testing.B) {
	ct, _ := cron.Parse("@monthly")
	ct.WithTestTimer(loopTimer((15 + 6) * 24 * time.Hour))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Next()
	}
}

func BenchmarkYearly(b *testing.B) {
	ct, _ := cron.Parse("@yearly")
	ct.WithTestTimer(loopTimer((365 + 182) * 24 * time.Hour))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct.Next()
	}
}
