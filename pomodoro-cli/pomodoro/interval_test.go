package pomodoro_test

import (
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
			config := pomodoro.NewIntervalConfig(repo, testCase.input[0], testCase.input[1], testCase.input[2])
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
