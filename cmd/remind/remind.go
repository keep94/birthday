package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/keep94/birthday"
	"github.com/keep94/birthday/cmd/remind/home"
	"github.com/keep94/birthday/cmd/remind/search"
	"github.com/keep94/consume2"
	"github.com/keep94/context"
	"github.com/keep94/toolbox/http_util"
	"github.com/keep94/toolbox/logging"
	"github.com/keep94/weblogs"
)

const (
	kMaxRows = 100
)

var (
	kFirstN = consume2.PSlice[birthday.Milestone](0, kMaxRows)
)

var (
	fFile      string
	fDaysAhead int
	fPort      string
)

func main() {
	flag.Parse()
	if fFile == "" {
		fmt.Println("Need to specify at least -file flag.")
		flag.Usage()
		os.Exit(1)
	}
	http.HandleFunc("/", rootRedirect)
	http.Handle(
		"/home",
		&home.Handler{File: fFile, DaysAhead: fDaysAhead, FirstN: kFirstN})
	http.Handle("/search", &search.Handler{File: fFile})
	defaultHandler := context.ClearHandler(
		weblogs.HandlerWithOptions(
			http.DefaultServeMux,
			&weblogs.Options{Logger: logging.ApacheCommonLoggerWithLatency()}))
	if err := http.ListenAndServe(fPort, defaultHandler); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http_util.Redirect(w, r, "/home")
	} else {
		http_util.Error(w, http.StatusNotFound)
	}
}

func init() {
	flag.StringVar(&fFile, "file", "", "Birthday file")
	flag.IntVar(&fDaysAhead, "days_ahead", 21, "Days ahead")
	flag.StringVar(&fPort, "http", ":8080", "Port to bind")
}
