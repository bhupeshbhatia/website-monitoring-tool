package dashboard

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bhupeshbhatia/website-monitoring-tool/statsagent"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

var alerts []string

type View struct {
	UpdateInterval int   `json:"updateInterval"`
	TimeFrame      int64 `json:"timeFrame"`
}

// Run displays the statistics in terminal
func Run(urls []string, views []View, alertc chan string, done context.CancelFunc) {
	//Use go console user interface - gui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout(g, views))

	//There should be a goroutine for alerts message
	updateViews(views, g, urls)

	go func() {
		for alertMessage := range alertc {
			alerts = append(alerts, alertMessage)
			updateAlertView(g)
		}
	}()

	//need to set keybinds for gocui
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			done()
			return quit(g, v)
		}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {

		log.Panicln(err)
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {

		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func updateViews(views []View, g *gocui.Gui, urls []string) {

	//need to update views
	for index := range views {
		go func(currentIndex int, currentView *View) {
			ticker := time.NewTicker(time.Duration(currentView.UpdateInterval) * time.Second)
			for {
				select {
				case t := <-ticker.C:
					res := statsagent.GetStats(urls, t, currentView.TimeFrame)
					g.Update(func(g *gocui.Gui) error {
						v, err := g.View(strconv.Itoa(int(currentView.TimeFrame)))
						if err != nil {
							return err
						}
						v.Clear()

						header := color.New(color.FgYellow, color.Bold)
						header.Fprintln(v, fmt.Sprintf("%-30v %21v %21v %21v %21v %21v %25v\n", "Website", "Availability", "Avg Response Time", "Max Response Time", "Avg TTFB", "Max TTFB", "Status Codes"))

						for _, url := range urls {
							value := res[url]
							statusCodeSlice := make([]string, 0)
							for code, count := range value.StatusCodeCount {
								statusCodeSlice = append(statusCodeSlice, fmt.Sprintf("%v:%v", code, count))
							}
							statusCodeStr := fmt.Sprintf("[%v]", strings.Join(statusCodeSlice, " "))
							fmt.Fprintln(v, fmt.Sprintf("%-30v %20.2f%% %21v %21v %21v %21v %25v", url, 100*value.Availability, value.AvgResponseTime, value.MaxResponseTime, value.AvgTimeToFirstByte, value.MaxTimeToFirstByte, statusCodeStr))
						}
						return nil
					})
				}
			}
		}(index, &views[index])
	}

}

//UpdateAlertView function?
func updateAlertView(g *gocui.Gui) {

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("alerts")
		if err != nil {
			return err
		}
		v.Clear()

		for i := len(alerts) - 1; i >= 0; i-- {
			fmt.Fprintln(v, alerts[i])
		}
		return nil
	})
}

//layout function - gocui library
func layout(g *gocui.Gui, views []View) func(*gocui.Gui) error {
	maxX, maxY := g.Size()
	return func(g *gocui.Gui) error {
		// Set stats views
		for index, view := range views {
			v, err := g.SetView(strconv.Itoa(int(view.TimeFrame)), 0, index*(maxY/3), maxX, (index+1)*(maxY/3))
			v.FgColor = gocui.ColorCyan
			if err != nil {
				if err != gocui.ErrUnknownView {
					log.Panic("Error setting views")
				}

				loadingMessage := color.New(color.FgMagenta)
				loadingMessage.Fprintln(v, fmt.Sprintf("\n\n%v One moment, we're waiting for statistics for the last %vs...", "⌛ ", view.TimeFrame))
			}
			v.Title = fmt.Sprintf(" Statistics for the last %vs (updated every %vs) ", view.TimeFrame, view.UpdateInterval)
			v.Wrap = true
		}

		// Set alerts view
		v, err := g.SetView("alerts", 0, 2*(maxY/3), maxX, maxY)
		v.FgColor = gocui.ColorCyan
		if err != nil {
			if err != gocui.ErrUnknownView {
				log.Panic("Error setting views")
			}
		}
		v.Title = fmt.Sprintf(" Alerts ")
		v.Wrap = true
		return nil
	}
}

//Scroll view
func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
