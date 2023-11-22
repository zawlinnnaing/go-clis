package pomodoro

import (
	"fmt"
	"time"
)

type LineSeries struct {
	Name   string
	Labels map[int]string
	Values []float64
}

func DailySummary(day time.Time, config *IntervalConfig) ([]time.Duration, error) {
	var durations []time.Duration
	pomoRes, err := config.repo.CategorySummary(day, CategoryPomodoro)
	if err != nil {
		return nil, err
	}
	breaksRes, err := config.repo.CategorySummary(day, "%Break")
	if err != nil {
		return nil, err
	}
	durations = append(durations, pomoRes, breaksRes)

	return durations, nil
}

func RangeSummary(end time.Time, n int, config *IntervalConfig) ([]LineSeries, error) {
	pomoLine := LineSeries{
		Name:   "Pomodoro",
		Labels: make(map[int]string),
		Values: make([]float64, n),
	}

	breaksLine := LineSeries{
		Name:   "Breaks",
		Labels: make(map[int]string),
		Values: make([]float64, n),
	}
	// TODO: refactor this to use goroutines
	for i := 0; i < n; i++ {
		date := end.AddDate(0, 0, -i)
		dailyDuration, err := DailySummary(date, config)
		if err != nil {
			return nil, err
		}
		label := fmt.Sprintf("%02d/%s", date.Day(), date.Format("Jan"))
		pomoLine.Labels[i] = label
		pomoLine.Values[i] = dailyDuration[0].Seconds()

		breaksLine.Labels[i] = label
		breaksLine.Values[i] = dailyDuration[1].Seconds()
	}
	return []LineSeries{
		pomoLine, breaksLine,
	}, nil
}
