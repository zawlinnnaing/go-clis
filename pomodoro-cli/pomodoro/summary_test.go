package pomodoro_test

import (
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"testing"
	"time"
)

func initData(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	repo, cleanup := getRepo(t)

	intervals := []pomodoro.Interval{
		{
			ID:              0,
			StartTime:       time.Now(),
			PlannedDuration: time.Minute * 20,
			ActualDuration:  time.Minute * 20,
			Category:        pomodoro.CategoryPomodoro,
			State:           pomodoro.StateDone,
		},
		{
			ID:              1,
			StartTime:       time.Now(),
			PlannedDuration: time.Minute * 5,
			ActualDuration:  time.Minute * 5,
			Category:        pomodoro.CategoryShortBreak,
			State:           pomodoro.StateDone,
		},
		{
			ID:              2,
			StartTime:       time.Now(),
			PlannedDuration: time.Minute * 20,
			ActualDuration:  time.Minute * 20,
			Category:        pomodoro.CategoryPomodoro,
			State:           pomodoro.StateDone,
		},
		{
			ID:              3,
			StartTime:       time.Now().AddDate(0, 0, -1),
			PlannedDuration: time.Minute * 20,
			ActualDuration:  time.Minute * 20,
			Category:        pomodoro.CategoryPomodoro,
			State:           pomodoro.StateDone,
		},
	}

	for _, interval := range intervals {
		if _, err := repo.Create(interval); err != nil {
			t.Fatal(err)
		}
	}

	return repo, func() {
		cleanup()
	}
}

func TestDailySummary(t *testing.T) {
	repo, cleanup := initData(t)
	defer cleanup()
	expectDailySummary := []time.Duration{
		40 * time.Minute,
		5 * time.Minute,
	}
	intervalConfig := pomodoro.NewIntervalConfig(repo, 20*time.Minute, 5*time.Minute, 15*time.Minute, 0, false)
	actualSummary, err := pomodoro.DailySummary(time.Now(), intervalConfig)
	if err != nil {
		t.Fatal(err)
	}
	if expectDailySummary[0] != actualSummary[0] || expectDailySummary[1] != actualSummary[1] {
		t.Errorf("Expect daily summary to be: %v, received: %v", expectDailySummary, actualSummary)
	}
}

func TestWeeklySummary(t *testing.T) {
	repo, cleanup := initData(t)
	defer cleanup()
	intervalConfig := pomodoro.NewIntervalConfig(repo, 20*time.Minute, 5*time.Minute, 15*time.Minute, 0, false)
	weeklySummaries, err := pomodoro.RangeSummary(time.Now(), 3, intervalConfig)
	if err != nil {
		t.Fatal(err)
	}

	expSummaries := [][]float64{
		{
			(40 * time.Minute).Seconds(),
			(20 * time.Minute).Seconds(),
			0,
		},
		{
			(5 * time.Minute).Seconds(),
			0,
			0,
		},
	}
	for idx, summary := range weeklySummaries {
		expSummary := expSummaries[idx]
		for valIdx, value := range summary.Values {
			expVal := expSummary[valIdx]
			if value != expVal {
				t.Errorf("Expect summary to be: %v, received: %v, at value index: %v, summary index: %v",
					expVal, value, valIdx, idx,
				)
			}
		}
	}
}
