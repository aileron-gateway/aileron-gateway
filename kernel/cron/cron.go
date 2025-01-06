package cron

import (
	"sync/atomic"
	"time"
)

// Crontab is a cron scheduler.
// Crontab implements Scheduler interface.
type Crontab struct {
	loc    *time.Location
	timer  func() time.Time
	second uint64
	minute uint64
	hour   uint64
	day    uint64
	month  uint64
	week   uint64
}

func (c *Crontab) Now() time.Time {
	return c.timer()
}

func (c *Crontab) NewJob(job func()) *CronJob {
	return &CronJob{
		cron: c,
		job:  job,
	}
}

func (c *Crontab) WithLocation(loc *time.Location) {
	if loc != nil {
		c.loc = loc
	}
}

func (c *Crontab) WithTestTimer(timer func() time.Time) {
	if timer != nil {
		c.timer = timer
	}
}

func (c *Crontab) Next() time.Time {
	now := c.timer().In(c.loc)
	hour, min, sec := now.Clock()

	var year, day, week int
	var month time.Month

	for {
		year, month, day = now.Date()
		week = int(now.Weekday())

		if c.month&(1<<month) == 0 {
			month += 1
			if month > 12 {
				month = 1
				year += 1
			}
			now = time.Date(year, month, 1, 0, 0, 0, 0, c.loc)
			hour, min, sec = -1, -1, -1
			continue
		}
		if c.day&(1<<day) == 0 || c.week&(1<<week) == 0 {
			now = now.Add(24 * time.Hour)
			hour, min, sec = -1, -1, -1
			continue
		}

		s := nextTime(c.second, sec, 59)

		m := min
		if min == -1 || !(s > sec) || c.minute&(1<<min) == 0 {
			m = nextTime(c.minute, min, 59)
		}

		h := hour
		if hour == -1 || !(m > min || (m == min && s > sec)) || c.hour&(1<<hour) == 0 {
			h = nextTime(c.hour, hour, 23)
			if !(h > hour) {
				now = now.Add(24 * time.Hour)
				hour, min, sec = -1, -1, -1
				continue
			}
		}

		return time.Date(year, month, day, h, m, s, 0, c.loc)
	}
}

func (c *Crontab) valid() bool {
	now := c.timer().In(c.loc)
	var day, week int
	var month time.Month
	for i := 0; i < 10*366; i++ { // Check 10 years.
		now = now.Add(24 * time.Hour)
		month = now.Month()
		day = now.Day()
		week = int(now.Weekday())
		if c.day&(1<<day) > 0 && c.week&(1<<week) > 0 && c.month&(1<<month) > 0 {
			return true
		}
	}
	return false
}

func nextTime(targets uint64, now int, max int) int {
	for i := now + 1; i <= max; i++ {
		if targets&(1<<i) > 0 {
			return i
		}
	}
	for i := 0; i <= now; i++ {
		if targets&(1<<i) > 0 {
			return i
		}
	}
	return 0
}

// CronJob is the job that runs the job
// at every scheduled times.
// This job do not run the new job
// if there is already running job.
type CronJob struct {
	cron *Crontab
	job  func()
	stop chan struct{}

	running atomic.Bool
	waiter  func(time.Duration) time.Duration
}

// WithTestWaiter specify the wait duration
// instead of actual scheduled wait duration.
// THIS TESTING ONLY.
func (j *CronJob) WithTestWaiter(f func(scheduled time.Duration) (use time.Duration)) {
	j.waiter = f
}

func (j *CronJob) Start() {
	if j.stop != nil {
		return
	}
	j.stop = make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		target := j.cron.Next()
		wait := target.Sub(j.cron.Now())
		wait = max(wait, time.Nanosecond) // Duration must not be 0 or negative for ticker.
		calibrate := false
		if wait > time.Hour {
			calibrate = true
			wait = 95 * wait / 100
		}
		if j.waiter != nil {
			wait = j.waiter(wait) // Use alternative wait duration when testing.
		}
		ticker.Reset(wait)

		select {
		case <-ticker.C:
			if calibrate {
				continue
			}
		case <-j.stop:
			return
		}
		if j.running.Load() {
			continue // Already running.
		}
		j.running.Store(true)
		go func() { // Run the job.
			defer func() {
				// Cron discards all panics from job func.
				// It is job implementers responsibility
				// to handle all errors and panics in the job.
				// So, ignore the returned value.
				_ = recover()
			}()
			defer j.running.Store(false)
			j.job()
		}()
	}
}

func (j *CronJob) Stop() {
	if j.stop != nil {
		close(j.stop)
		j.stop = nil
	}
}

func max(a, b time.Duration) time.Duration {
	if a < b {
		return b
	}
	return a
}
