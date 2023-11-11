package app

import (
	"context"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	donTimer       *donut.Donut
	disType        *segmentdisplay.SegmentDisplay
	textInfo       *text.Text
	textTime       *text.Text
	updateDonTimer chan []int
	updateDisType  chan string
	updateTextInfo chan string
	updateTextTime chan string
}

func (widgets *widgets) update(time []int, disType, textInfo, textTime string, redraw chan<- bool) {
	if textInfo != "" {
		widgets.updateTextInfo <- textInfo
	}
	if disType != "" {
		widgets.updateDisType <- disType
	}
	if textTime != "" {
		widgets.updateTextTime <- textTime
	}
	if len(time) > 0 {
		widgets.updateDonTimer <- time
	}
	redraw <- true
}

func newWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error
	w.updateDonTimer = make(chan []int)
	w.updateTextTime = make(chan string)
	w.updateTextInfo = make(chan string)
	w.updateDisType = make(chan string)

	w.donTimer, err = newDonut(ctx, w.updateDonTimer, errorCh)
	if err != nil {
		return nil, err
	}
	w.disType, err = newSegmentDisplay(ctx, w.updateDisType, errorCh)
	if err != nil {
		return nil, err
	}
	w.textInfo, err = newText(ctx, w.updateTextInfo, errorCh)
	if err != nil {
		return nil, err
	}
	w.textTime, err = newText(ctx, w.updateTextTime, errorCh)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func newDonut(ctx context.Context, updateDonutCh <-chan []int, errorCh chan<- error) (*donut.Donut, error) {
	don, err := donut.New(
		donut.Clockwise(),
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case val := <-updateDonutCh:
				if val[0] <= val[1] {
					errorCh <- don.Absolute(val[0], val[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return don, nil
}

func newSegmentDisplay(ctx context.Context, updateDisTypeCh <-chan string, errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	segDisplay, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case val := <-updateDisTypeCh:
				if val == "" {
					val = " "
				}
				errorCh <- segDisplay.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(val),
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	return segDisplay, nil
}

func newText(ctx context.Context, updateTextCh <-chan string, errorCh chan<- error) (*text.Text, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case t := <-updateTextCh:
				txt.Reset()
				errorCh <- txt.Write(t)
			case <-ctx.Done():
				return
			}
		}
	}()
	return txt, nil
}
