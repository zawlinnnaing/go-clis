package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"time"
)

type buttonSet struct {
	btnStart *button.Button
	btnPause *button.Button
}

func newButtonSet(ctx context.Context, config *pomodoro.IntervalConfig, widgets *widgets, summary *summary, redrawCh chan<- bool, errorCh chan<- error) (*buttonSet, error) {
	start := func(interval pomodoro.Interval) {
		var message string
		switch interval.Category {
		case pomodoro.CategoryPomodoro:
			message = "Focus on your task"
		case pomodoro.CategoryShortBreak:
			message = "Take a short break"
		case pomodoro.CategoryLongBreak:
			message = "Take a long break"
		}
		widgets.update([]int{}, interval.Category, message, "", redrawCh)
	}

	periodic := func(interval pomodoro.Interval) {
		widgets.update(
			[]int{int(interval.ActualDuration), int(interval.PlannedDuration)},
			"", "",
			fmt.Sprint(interval.PlannedDuration-interval.ActualDuration),
			redrawCh,
		)
	}

	var end func(interval pomodoro.Interval)

	end = func(interval pomodoro.Interval) {
		isBefore := time.Now().Before(config.RunUntil)
		if isBefore {
			go func() {
				nextInterval, err := pomodoro.GetInterval(config)
				errorCh <- err
				errorCh <- nextInterval.Start(ctx, config, start, periodic, end)
			}()
		} else {
			widgets.update([]int{}, "", "Nothing running...", "", redrawCh)
		}
		summary.update(redrawCh)
	}

	startInterval := func() {
		interval, err := pomodoro.GetInterval(config)
		errorCh <- err
		errorCh <- interval.Start(ctx, config, start, periodic, end)
	}

	pauseInterval := func() {
		interval, err := pomodoro.GetInterval(config)
		if err != nil {
			errorCh <- err
			return
		}
		if err = interval.Pause(config); err != nil {
			if errors.Is(err, pomodoro.ErrIntervalNotRunning) {
				return
			}
			errorCh <- err
			return
		}
		widgets.update([]int{}, "", "Paused. Press start to continue...", "", redrawCh)
	}

	btnStart, err := button.New("(s)tart", func() error {
		go startInterval()
		return nil
	},
		button.GlobalKey('s'),
		button.Height(2),
		button.WidthFor("(p)ause"),
	)
	if err != nil {
		return nil, err
	}
	btnPause, err := button.New("(p)ause", func() error {
		go pauseInterval()
		return nil
	},
		button.GlobalKey('p'),
		button.Height(2),
		button.FillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		return nil, err
	}
	return &buttonSet{btnPause: btnPause, btnStart: btnStart}, nil
}
