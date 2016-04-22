package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

type Devops interface {
	AvgCPUUsageDayByHour(*QueryBytes)
	AvgCPUUsageMonthByDay(*QueryBytes)
}

type InfluxDevops struct {
	DatabaseName   string
	DayIntervals   TimeIntervals
	WeekIntervals  TimeIntervals
	MonthIntervals TimeIntervals
}

type TimeInterval struct {
	Start, End time.Time
}

func (ti *TimeInterval) StartString() string {
	return ti.Start.Format(time.RFC3339)
}

func (ti *TimeInterval) EndString() string {
	return ti.End.Format(time.RFC3339)
}

type TimeIntervals []TimeInterval

func (tis TimeIntervals) RandChoice() *TimeInterval {
	return &tis[rand.Intn(len(tis))]
}

func NewTimeIntervals(start, end time.Time, window time.Duration) TimeIntervals {
	xs := TimeIntervals{}
	for start.Add(window).Before(end) || start.Add(window).Equal(end) {
		x := TimeInterval{
			Start: start,
			End:   start.Add(window),
		}
		xs = append(xs, x)

		start = start.Add(window)
	}

	return xs
}

func NewInfluxDevops(databaseName string, start, end time.Time) *InfluxDevops {
	if !start.Before(end) {
		panic("bad time order")
	}
	return &InfluxDevops{
		DatabaseName:   databaseName,
		DayIntervals:   NewTimeIntervals(start, end, 24*time.Hour),
		WeekIntervals:  NewTimeIntervals(start, end, 7*24*time.Hour),
		MonthIntervals: NewTimeIntervals(start, end, 31*24*time.Hour),
	}
}

func (d *InfluxDevops) AvgCPUUsageDayByHour(q *QueryBytes) {
	interval := d.DayIntervals.RandChoice()

	v := url.Values{}
	v.Set("db", d.DatabaseName)
	v.Set("q", fmt.Sprintf("SELECT mean(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h)", interval.StartString(), interval.EndString()))

	q.Method = []byte("GET")
	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
	q.Body = nil
}

func (d *InfluxDevops) AvgCPUUsageWeekByHour(q *QueryBytes) {
	interval := d.WeekIntervals.RandChoice()

	v := url.Values{}
	v.Set("db", d.DatabaseName)
	v.Set("q", fmt.Sprintf("SELECT mean(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h)", interval.StartString(), interval.EndString()))

	q.Method = []byte("GET")
	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
	q.Body = nil
}

func (d *InfluxDevops) AvgCPUUsageMonthByDay(q *QueryBytes) {
	interval := d.MonthIntervals.RandChoice()

	v := url.Values{}
	v.Set("db", d.DatabaseName)
	v.Set("q", fmt.Sprintf("SELECT mean(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1d)", interval.StartString(), interval.EndString()))

	q.Method = []byte("GET")
	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
	q.Body = nil
}