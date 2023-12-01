package pomodoro_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"testing"
	"time"
)

func TestNewIntervalConfig(t *testing.T) {
	testCases := []struct {
		name   string
		input  [3]time.Duration
		expect pomodoro.IntervalConfig
	}{
		{
			name: "Default",
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   pomodoro.DefaultPomodoroDuration,
				ShortBreakDuration: pomodoro.DefaultShorBreakDuration,
				LongBreakDuration:  pomodoro.DefaultLongBreakDuration,
			},
		},
		{
			name:  "SingleInput",
			input: [3]time.Duration{30 * time.Minute},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   30 * time.Minute,
				ShortBreakDuration: pomodoro.DefaultShorBreakDuration,
				LongBreakDuration:  pomodoro.DefaultLongBreakDuration,
			},
		},
		{
			name:  "MultiInput",
			input: [3]time.Duration{30 * time.Minute, 7 * time.Minute, 12 * time.Minute},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   30 * time.Minute,
				ShortBreakDuration: 7 * time.Minute,
				LongBreakDuration:  12 * time.Minute,
			},
		},
	}
	repo, cleanUp := getRepo(t)
	defer cleanUp()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := pomodoro.NewIntervalConfig(repo, testCase.input[0], testCase.input[1], testCase.input[2], 0)
			if config.PomodoroDuration != testCase.expect.PomodoroDuration {
				t.Errorf("Expected pomo duration: %d, recevied: %d", testCase.expect.PomodoroDuration, config.PomodoroDuration)
			}
			if config.ShortBreakDuration != testCase.expect.ShortBreakDuration {
				t.Errorf("Expected short duration: %d, recevied: %d", testCase.expect.ShortBreakDuration, config.ShortBreakDuration)
			}
			if config.LongBreakDuration != testCase.expect.LongBreakDuration {
				t.Errorf("Expected long duration: %d, received: %d", testCase.expect.LongBreakDuration, config.LongBreakDuration)
			}
		})
	}
}

func TestGetInterval(t *testing.T) {
	repo, cleanUp := getRepo(t)
	defer cleanUp()
	const shortBreakDuration = time.Millisecond
	const pomodoroDuration = 3 * time.Millisecond
	const longBreakDuration = 2 * time.Millisecond
	config := pomodoro.NewIntervalConfig(repo, pomodoroDuration, shortBreakDuration, longBreakDuration, 0)
	for i := 1; i <= 16; i++ {
		var (
			expCategory string
			expDuration time.Duration
		)
		switch {
		case i%2 != 0:
			expCategory = pomodoro.CategoryPomodoro
			expDuration = pomodoroDuration
		case i%8 == 0:
			expCategory = pomodoro.CategoryLongBreak
			expDuration = longBreakDuration
		default:
			expCategory = pomodoro.CategoryShortBreak
			expDuration = shortBreakDuration
		}
		testName := fmt.Sprintf("%d. %s", i, expCategory)
		t.Run(testName, func(t *testing.T) {
			res, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			noop := func(i pomodoro.Interval) {}
			err = res.Start(context.Background(), config, noop, noop, noop)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if res.Category != expCategory {
				t.Errorf("Expected category: %s, received: %s", expCategory, res.Category)
			}
			if res.PlannedDuration != expDuration {
				t.Errorf("Expected duration: %s, recevied: %s", expDuration, res.PlannedDuration)
			}
			if res.State != pomodoro.StateRunning {
				t.Errorf("Expected state to be running, recevied: %d", res.State)
			}
			interval, err := repo.ByID(res.ID)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if interval.State != pomodoro.StateDone {
				t.Errorf("Expected state to be done, recevied: %d", interval.State)
			}
		})
	}
}

func TestPause(t *testing.T) {
	const duration = 2 * time.Second
	repo, cleanup := getRepo(t)
	defer cleanup()

	config := pomodoro.NewIntervalConfig(repo, duration, duration, duration, 0)
	testCases := []struct {
		name        string
		start       bool
		expState    int
		expDuration time.Duration
	}{
		{
			name:        "NotStarted",
			start:       false,
			expState:    pomodoro.StateNotStarted,
			expDuration: 0,
		},
		{
			name:        "Pause",
			start:       true,
			expState:    pomodoro.StatePaused,
			expDuration: time.Second,
		},
	}

	expErr := pomodoro.ErrIntervalNotRunning

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			interval, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}
			start := func(i pomodoro.Interval) {}
			end := func(i pomodoro.Interval) {
				t.Errorf("Should not reach end of the interval")
			}
			periodic := func(i pomodoro.Interval) {
				if err = i.Pause(config); err != nil {
					t.Fatal(err)
				}
			}
			if testCase.start {
				if err = interval.Start(ctx, config, start, periodic, end); err != nil {
					t.Fatal(err)
				}
			}
			interval, err = pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}
			err = interval.Pause(config)
			if !errors.Is(err, expErr) {
				t.Errorf("Expected error: %v, received: %v", expErr, err)
			}
			interval, err = repo.ByID(interval.ID)
			if err != nil {
				t.Fatal(err)
			}
			if interval.State != testCase.expState {
				t.Errorf("Expected state: %d, recevied: %d", testCase.expState, interval.State)
			}
			if interval.ActualDuration != testCase.expDuration {
				t.Errorf("Expected duration: %d, recevied: %d", testCase.expDuration, interval.ActualDuration)
			}
			cancel()
		})
	}
}

func TestStart(t *testing.T) {
	const duration = 2 * time.Second
	repo, cleanup := getRepo(t)
	defer cleanup()
	config := pomodoro.NewIntervalConfig(repo, duration, duration, duration, 0)
	testCases := []struct {
		name        string
		cancel      bool
		expState    int
		expDuration time.Duration
	}{
		{name: "Canceled", cancel: true, expState: pomodoro.StateCancelled, expDuration: time.Second},
		{name: "Finished", cancel: false, expState: pomodoro.StateDone, expDuration: duration},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			interval, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}
			start := func(i pomodoro.Interval) {
				if i.State != pomodoro.StateRunning {
					t.Errorf("Expected state to be: %d, recevied: %d", pomodoro.StateRunning, i.State)
				}
				if i.PlannedDuration <= i.ActualDuration {
					t.Errorf("Expected planned duration: %d to be less than actual: %d", i.PlannedDuration, i.ActualDuration)
				}
			}
			periodic := func(i pomodoro.Interval) {
				if i.State != pomodoro.StateRunning {
					t.Errorf("Expected state to be: %d, recevied: %d", pomodoro.StateRunning, i.State)
				}
				if testCase.cancel {
					cancel()
				}
			}
			end := func(i pomodoro.Interval) {
				if i.State != testCase.expState {
					t.Errorf("Expected state: %d, recevied state: %d", testCase.expState, i.State)
				}
				if testCase.cancel {
					t.Errorf("End callback should not be reached")
				}
			}
			if err = interval.Start(ctx, config, start, periodic, end); err != nil {
				t.Fatal(err)
			}
			interval, err = repo.ByID(interval.ID)
			if err != nil {
				t.Fatal(err)
			}
			if interval.State != testCase.expState {
				t.Errorf("Expected state: %d, recevied state: %d", testCase.expState, interval.State)
			}
			if interval.ActualDuration != testCase.expDuration {
				t.Errorf("Expected duration: %d, actual duration: %d", testCase.expDuration, interval.ActualDuration)
			}
		})
	}
}
