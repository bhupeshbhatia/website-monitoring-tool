package dashboard

import (
	"context"

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

	//There should be a goroutine for alerts message

	//need to set keybinds for gocui

}

func updateViews(views []View, g *gocui.Gui, urls []string) {

	//need to update views

}

//AlertView function?

//layout function - gocui library
// func layout(g *gocui.Gui, views []View) func(*gocui.Gui) error {
// }

//Scroll view
// func scrollView(v *gocui.View, dy int) error {
// }
