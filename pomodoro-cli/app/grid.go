package app

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func newGrid(btnSet *buttonSet, w *widgets, summary *summary, api terminalapi.Terminal) (*container.Container, error) {
	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(30,
			grid.ColWidthPercWithOpts(30, []container.Option{
				container.Border(linestyle.Light),
				container.BorderTitle("Press Q to Quit"),
			},
				grid.RowHeightPerc(80, grid.Widget(w.donTimer)),
				grid.RowHeightPercWithOpts(20, []container.Option{
					container.AlignHorizontal(align.HorizontalCenter),
				}, grid.Widget(w.textTime,
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignVertical(align.VerticalMiddle),
				)),
			),
			grid.ColWidthPerc(70,
				grid.RowHeightPerc(80, grid.Widget(w.disType, container.Border(linestyle.Light))),
				grid.RowHeightPerc(20, grid.Widget(w.textInfo, container.Border(linestyle.Light))),
			),
		),
	)
	builder.Add(
		grid.RowHeightPerc(10,
			grid.ColWidthPerc(50, grid.Widget(btnSet.btnStart)),
			grid.ColWidthPerc(50, grid.Widget(btnSet.btnPause)),
		))
	builder.Add(
		grid.RowHeightPerc(60,
			grid.ColWidthPerc(30,
				grid.Widget(
					summary.bcDaily,
					container.Border(linestyle.Light),
					container.BorderTitle("Daily Summary (Minutes)"),
				),
			),
			grid.ColWidthPerc(70,
				grid.Widget(
					summary.lcWeekly,
					container.Border(linestyle.Light),
					container.BorderTitle("Weekly Summary"),
				),
			),
		))

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	c, err := container.New(api, gridOpts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
