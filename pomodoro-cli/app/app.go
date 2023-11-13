package app

import (
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"image"
	"time"
)

type App struct {
	ctx        context.Context
	controller *termdash.Controller
	redrawCh   chan bool
	errorCh    chan error
	terminal   *tcell.Terminal
	size       image.Point
}

func New(config *pomodoro.IntervalConfig) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())
	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}
	redrawCh := make(chan bool)
	errorCh := make(chan error)

	w, err := newWidgets(ctx, errorCh)
	if err != nil {
		return nil, err
	}
	btnSet, err := newButtonSet(ctx, config, w, redrawCh, errorCh)
	if err != nil {
		return nil, err
	}

	terminal, err := tcell.New()
	if err != nil {
		return nil, err
	}
	container, err := newGrid(btnSet, w, terminal)
	if err != nil {
		return nil, err
	}
	controller, err := termdash.NewController(terminal, container, termdash.KeyboardSubscriber(quitter))
	if err != nil {
		return nil, err
	}
	return &App{
		ctx:        ctx,
		controller: controller,
		redrawCh:   redrawCh,
		errorCh:    errorCh,
		terminal:   terminal,
	}, nil
}

func (a *App) resize() error {
	if !a.size.Eq(a.terminal.Size()) {
		return nil
	}
	a.size = a.terminal.Size()
	if err := a.terminal.Clear(); err != nil {
		return err
	}
	return a.controller.Redraw()
}

func (a *App) Run() error {
	defer a.controller.Close()
	defer a.terminal.Close()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-a.redrawCh:
			if err := a.controller.Redraw(); err != nil {
				return err
			}
		case err := <-a.errorCh:
			if err != nil {
				return err
			}
		case <-a.ctx.Done():
			return nil
		case <-ticker.C:
			if err := a.resize(); err != nil {
				return err
			}
		}
	}
}
