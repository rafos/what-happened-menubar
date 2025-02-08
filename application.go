package main

import (
	_ "embed"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"log/slog"
	"time"
)

const GITHUB_URL = "https://github.com/rafos/what-happened-menubar/"

//go:embed icon.png
var iconByte []byte

var (
	title            = "What Happened"
	fullTitle        = title + " Today"
	eventsItems      []*systray.MenuItem
	birthsItems      []*systray.MenuItem
	deathsItems      []*systray.MenuItem
	errorMessageItem *systray.MenuItem
	refreshItem      *systray.MenuItem
	aboutItem        *systray.MenuItem
	quitItem         *systray.MenuItem
)

type application struct {
}

func newApplication() *application {
	return &application{}
}

func (a *application) start() {
	slog.Info("Starting application " + fullTitle)
	systray.Run(onReady, nil)
}

func onReady() {
	today := time.Now()

	generateHeader(today)
	generateItems(today)
	generateFooter()

	for {
		select {
		case <-refreshItem.ClickedCh:
			slog.Info("Refreshing")
			go func() {
				for _, eventsItem := range eventsItems {
					eventsItem.Hide()
				}
			}()

			go func() {
				for _, birthsItem := range birthsItems {
					birthsItem.Hide()
				}
			}()

			go func() {
				for _, deathsItem := range deathsItems {
					deathsItem.Hide()
				}
			}()

			if errorMessageItem != nil {
				errorMessageItem.Hide()
			}
			refreshItem.Hide()
			aboutItem.Hide()
			quitItem.Hide()

			today = time.Now()

			generateItems(today)
			generateFooter()

		case <-aboutItem.ClickedCh:
			if err := open.Run(GITHUB_URL); err != nil {
				continue
			}

		case <-quitItem.ClickedCh:
			slog.Info("Quitting")
			systray.Quit()
			return
		}
	}
}

func generateItems(today time.Time) {
	events, err := getAllEventsFrom(today)
	if err != nil {
		errorMessageItem = systray.AddMenuItem("Something went wrong. Please try again later.", "")
		return
	}

	if len(events.Data.Events) > 0 {
		eventsItem := systray.AddMenuItem("Events", "Events")

		for _, event := range events.Data.Events {
			itemTitle := fmt.Sprintf("%s: %s", event.Year, event.Text)
			eventsItem.AddSubMenuItem(itemTitle, "")
		}
		eventsItems = append(eventsItems, eventsItem)
	}

	if len(events.Data.Births) > 0 {
		birthsItem := systray.AddMenuItem("Births", "Births")

		for _, event := range events.Data.Births {
			itemTitle := fmt.Sprintf("%s: %s", event.Year, event.Text)
			birthsItem.AddSubMenuItem(itemTitle, "")
		}
		birthsItems = append(birthsItems, birthsItem)
	}

	if len(events.Data.Deaths) > 0 {
		deathsItem := systray.AddMenuItem("Deaths", "Deaths")

		for _, event := range events.Data.Deaths {
			itemTitle := fmt.Sprintf("%s: %s", event.Year, event.Text)
			deathsItem.AddSubMenuItem(itemTitle, "")
		}
		deathsItems = append(deathsItems, deathsItem)
	}
}

func generateHeader(today time.Time) {
	systray.SetIcon(iconByte)
	systray.SetTooltip(fullTitle)
	systray.AddMenuItem(fmt.Sprintf("%s on %s", title, today.Format("2 January")), "").Disable()
	systray.AddSeparator()
}

func generateFooter() {
	systray.AddSeparator()
	refreshItem = systray.AddMenuItem("Refresh", fmt.Sprintf("Refresh %s", fullTitle))
	aboutItem = systray.AddMenuItem("About", fmt.Sprintf("About %s", fullTitle))
	quitItem = systray.AddMenuItem("Quit", fmt.Sprintf("Quit %s", fullTitle))
}
