package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

const format = "2006-01-02_15-04-05"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var username, password string
	destination := filepath.Join(wd, fmt.Sprintf("membaz-export-%s.csv", time.Now().Format(format)))

	fs := flag.NewFlagSet("member-sync", flag.ExitOnError)
	fs.StringVar(&username, "username", "", "Membaz login username.")
	fs.StringVar(&password, "password", "", "Membaz login password.")
	fs.StringVar(&destination, "destination", destination, "Membaz file destination.")
	fs.Parse(os.Args[1:])

	if username == "" || username == "NOT_SET" {
		log.Fatal("Membaz username is not set.")
	}
	if password == "" || password == "NOT_SET" {
		log.Fatal("Membaz password is not set.")
	}
	log.SetOutput(os.Stderr)
	err = exportMembazCSV(username, password, destination)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote Membaz export to %s.", destination)
}

func exportMembazCSV(username, password, destination string) error {
	ctx, cancel := chromedp.NewContext(context.Background(),
		chromedp.WithLogf(log.Printf),
		chromedp.WithDebugf(log.Printf),
		chromedp.WithErrorf(log.Printf),
	)
	defer cancel()

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// set up a channel, so we can block later while we monitor the download
	// progress
	done := make(chan string, 1)
	// set up a listener to watch the download events and close the channel
	// when complete this could be expanded to handle multiple downloads
	// through creating a guid map, monitor download urls via
	// EventDownloadWillBegin, etc
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			completed := "(unknown)"
			if ev.TotalBytes != 0 {
				completed = fmt.Sprintf("%0.2f%%", ev.ReceivedBytes/ev.TotalBytes*100.0)
			}
			log.Printf("state: %s, completed: %s\n", ev.State.String(), completed)
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}
		}
	})

	// get working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(username)
	// to get a screenshot, add the following:
	// var res []byte
	// chromedp.CaptureScreenshot(&res),
	// os.WriteFile("screenshot.png", res, 0644)
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 768),
		chromedp.Navigate("https://www.membaz.com/login"),
		chromedp.WaitReady("body"),
		chromedp.SetValue(`lgnLogin_UserName`, username, chromedp.ByID),
		chromedp.SetValue(`lgnLogin_Password`, password, chromedp.ByID),
		chromedp.Click(`lgnLogin_LoginButton`, chromedp.ByID),
		chromedp.Sleep(time.Second*5),
		chromedp.Click(`a[title="View Reports"]`, chromedp.ByQuery),
		chromedp.Sleep(time.Second*2),
		chromedp.Click(`a[href="/App/User/CustomReports.aspx"]`, chromedp.ByQuery),
		chromedp.Sleep(time.Second*2),
		// configure headless browser downloads. note that
		// SetDownloadBehaviorBehaviorAllowAndName is preferred here over
		// SetDownloadBehaviorBehaviorAllow so that the file will be named as
		// the GUID. please note that it only works with 92.0.4498.0 or later
		// due to issue 1204880, see https://bugs.chromium.org/p/chromium/issues/detail?id=1204880
		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(wd).
			WithEventsEnabled(true),
		chromedp.Click(`ctl00_cpMainContent_btnExportToCSV`, chromedp.ByID),
	)
	if err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Note: Ignoring the net::ERR_ABORTED page error is essential here
		// since downloads will cause this error to be emitted, although the
		// download will still succeed.
		return err
	}

	// This will block until the chromedp listener closes the channel
	guid := <-done

	// We can predict the exact file location and name here because of how we
	// configured SetDownloadBehavior and WithDownloadPath
	err = os.Rename(filepath.Join(wd, guid), destination)
	return err
}
