package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

const (
	DefaultPomodoroDuration  = 25 * time.Minute
	DefaultShorBreakDuration = 5 * time.Minute
	DefaultLongBreakDuration = 15 * time.Minute
)

var (
	ErrNoIntervals        = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalCompleted  = errors.New("interval completed")
	ErrInvalidState       = errors.New("invalid state")
	ErrInvalidID          = errors.New("invalid ID")
)

type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

type IntervalConfig struct {
	repo                Repository
	PomodoroDuration    time.Duration
	ShortBreakDuration  time.Duration
	LongBreakDuration   time.Duration
	RunUntil            time.Time
	DisableNotification bool
}

type Callback func(interval Interval)

func NewIntervalConfig(repo Repository, pomodoro time.Duration, shortBreak time.Duration, longBreak time.Duration, total time.Duration, disableNoti bool) *IntervalConfig {
	intervalConfig := &IntervalConfig{
		repo:                repo,
		PomodoroDuration:    DefaultPomodoroDuration,
		ShortBreakDuration:  DefaultShorBreakDuration,
		LongBreakDuration:   DefaultLongBreakDuration,
		RunUntil:            time.Time{},
		DisableNotification: disableNoti,
	}
	if pomodoro > 0 {
		intervalConfig.PomodoroDuration = pomodoro
	}
	if shortBreak > 0 {
		intervalConfig.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		intervalConfig.LongBreakDuration = longBreak
	}
	if total > 0 {
		intervalConfig.RunUntil = time.Now().Add(total)
	}
	return intervalConfig
}

func nextCategory(repo Repository) (string, error) {
	lastInterval, err := repo.Last()
	if err != nil && errors.Is(err, ErrNoIntervals) {
		return CategoryPomodoro, nil
	}
	if err != nil {
		return "", err
	}
	if lastInterval.Category == CategoryShortBreak || lastInterval.Category == CategoryLongBreak {
		return CategoryPomodoro, nil
	}
	lastBreaks, err := repo.Breaks(3)
	if err != nil {
		return "", err
	}
	if len(lastBreaks) < 3 {
		return CategoryShortBreak, nil
	}

	for _, breakInterval := range lastBreaks {
		if breakInterval.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}
	return CategoryLongBreak, nil
}

func tick(ctx context.Context, id int64, config *IntervalConfig, start, periodic, end Callback) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	interval, err := config.repo.ByID(id)
	if err != nil {
		return err
	}
	expire := time.After(interval.PlannedDuration - interval.ActualDuration)

	start(interval)

	for {
		select {
		case <-ticker.C:
			interval, err = config.repo.ByID(id)
			if err != nil {
				return err
			}
			if interval.State == StatePaused {
				return nil
			}
			interval.ActualDuration += time.Second
			if err = config.repo.Update(interval); err != nil {
				return err
			}
			periodic(interval)
		case <-expire:
			interval, err = config.repo.ByID(id)
			if err != nil {
				return err
			}
			interval.State = StateDone
			if err = config.repo.Update(interval); err != nil {
				return err
			}
			ticker.Stop()
			end(interval)
			return nil
		case <-ctx.Done():
			interval, err = config.repo.ByID(id)
			if err != nil {
				return err
			}
			interval.State = StateCancelled
			return config.repo.Update(interval)
		}

	}
}

func newInterval(config *IntervalConfig) (Interval, error) {
	interval := &Interval{}
	category, err := nextCategory(config.repo)
	if err != nil {
		return *interval, err
	}
	interval.Category = category

	switch interval.Category {
	case CategoryPomodoro:
		interval.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		interval.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		interval.PlannedDuration = config.LongBreakDuration
	}

	if interval.ID, err = config.repo.Create(*interval); err != nil {
		return *interval, err
	}
	return *interval, nil
}

func GetInterval(config *IntervalConfig) (Interval, error) {
	interval := Interval{}
	var err error

	interval, err = config.repo.Last()
	if err != nil && !errors.Is(err, ErrNoIntervals) {
		return interval, err
	}
	if err == nil && (interval.State != StateDone && interval.State != StateCancelled) {
		return interval, nil
	}
	return newInterval(config)
}

func (interval *Interval) Start(ctx context.Context, config *IntervalConfig, start, periodic, end Callback) error {
	switch interval.State {
	case StateRunning:
		return nil
	case StateNotStarted:
		interval.StartTime = time.Now()
		fallthrough
	case StatePaused:
		interval.State = StateRunning
		if err := config.repo.Update(*interval); err != nil {
			return err
		}
		return tick(ctx, interval.ID, config, start, periodic, end)
	case StateCancelled, StateDone:
		return fmt.Errorf("%w: cannot start interval. State: %d", ErrIntervalCompleted, interval.State)
	default:
		return fmt.Errorf("%w: Invalid State, %d", ErrInvalidState, interval.State)
	}
}

func (interval *Interval) Pause(config *IntervalConfig) error {
	if interval.State != StateRunning {
		return ErrIntervalNotRunning
	}
	interval.State = StatePaused
	return config.repo.Update(*interval)
}
