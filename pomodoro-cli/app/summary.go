package app

import (
	"context"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"math"
	"time"
)

type summary struct {
	bcDaily      *barchart.BarChart
	lcWeekly     *linechart.LineChart
	updateDaily  chan bool
	updateWeekly chan bool
}

func (summary *summary) update(redrawCh chan<- bool) {
	summary.updateDaily <- true
	summary.updateWeekly <- true
	redrawCh <- true
}

func newSummary(ctx context.Context, config *pomodoro.IntervalConfig, redrawCh chan<- bool, errorCh chan<- error) (*summary, error) {
	s := &summary{}

	var err error
	s.updateDaily = make(chan bool)
	s.updateWeekly = make(chan bool)

	s.bcDaily, err = newBarChart(ctx, config, s.updateDaily, errorCh)
	if err != nil {
		return nil, err
	}
	s.lcWeekly, err = newLineChar(ctx, config, s.updateWeekly, errorCh)
	if err != nil {
		return nil, err
	}

	return s, err
}

func newLineChar(ctx context.Context, config *pomodoro.IntervalConfig, updateCh chan bool, errorCh chan<- error) (*linechart.LineChart, error) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorBlue)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
		linechart.YAxisFormattedValues(
			linechart.ValueFormatterSingleUnitDuration(time.Second, 0)),
	)
	if err != nil {
		return nil, err
	}

	updateWidget := func() error {
		weeklySummary, err := pomodoro.RangeSummary(time.Now(), 7, config)
		if err != nil {
			return err
		}
		err = lc.Series(
			weeklySummary[0].Name,
			weeklySummary[0].Values,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlue)),
			linechart.SeriesXLabels(weeklySummary[0].Labels),
		)
		if err != nil {
			return err
		}
		return lc.Series(
			weeklySummary[1].Name,
			weeklySummary[1].Values,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorYellow)),
			linechart.SeriesXLabels(weeklySummary[1].Labels),
		)
	}

	go func() {
		for {
			select {
			case <-updateCh:
				errorCh <- updateWidget()
			case <-ctx.Done():
				return
			}
		}
	}()

	if err = updateWidget(); err != nil {
		return nil, err
	}

	return lc, err
}

func newBarChart(ctx context.Context, config *pomodoro.IntervalConfig, updateCh chan bool, errorCh chan<- error) (*barchart.BarChart, error) {
	bc, err := barchart.New(
		barchart.ShowValues(),
		barchart.BarColors([]cell.Color{
			cell.ColorBlue,
			cell.ColorYellow,
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorBlack,
			cell.ColorBlack,
		}),
		barchart.Labels([]string{
			"Pomodoro",
			"Breaks",
		}),
	)
	if err != nil {
		return nil, err
	}

	updateWidget := func() error {
		dailySummary, err := pomodoro.DailySummary(time.Now(), config)
		if err != nil {
			return err
		}
		return bc.Values(
			[]int{
				int(dailySummary[0].Minutes()),
				int(dailySummary[1].Minutes()),
			},
			int(math.Max(dailySummary[0].Minutes(), dailySummary[1].Minutes())*1.1)+1,
		)
	}

	go func() {
		for {
			select {
			case <-updateCh:
				errorCh <- updateWidget()
			case <-ctx.Done():
				return
			}
		}
	}()

	if err = updateWidget(); err != nil {
		return nil, err
	}

	return bc, nil
}
