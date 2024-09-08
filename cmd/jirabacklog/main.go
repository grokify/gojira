package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grokify/gojira/cmd"
	"github.com/grokify/gojira/jirarest"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/time/timeutil"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	opts := cmd.Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(opts)

	jrClient, err := opts.Client()
	logutil.FatalErr(errorsutil.Wrap(err, "Client"))

	// cfg := gojira.NewConfigDefault()
	// cfg.BaseURL = jrClient.Config.ServerURL

	svc := jirarest.NewBacklogService(jrClient)

	is, _, err := svc.GetBacklogIssuesSetAll(opts.BoardID,
		// "type in (Bug,Story) AND status in (Ready,\"Engineering Design\",\"Ready for Grooming\")",
		"type in (Bug,Story) AND status in (Ready)",
	)
	logutil.FatalErr(errorsutil.Wrap(err, "GetBacklogIssuesSetAll"))

	fmt.Printf("COUNT (%d)\n", len(is.IssuesMap))

	countsByStatus := is.Counts()
	fmtutil.MustPrintJSON(countsByStatus)

	timeStats := is.TimeStats()
	fmtutil.MustPrintJSON(timeStats)

	timeStatsDays, err := timeStats.SecondsToDays()
	logutil.FatalErr(errorsutil.Wrap(err, "SecondsToDays"))
	fmtutil.MustPrintJSON(timeStatsDays)

	dt := time.Now().UTC()
	dtfs := dt.Format(timeutil.RFC3339Dash)
	outfile := "backlog-" + dtfs + ".json"
	err = is.WriteFileJSON(outfile, "", "  ")
	logutil.FatalErr(errorsutil.Wrap(err, "WriteFileJSON"))

	tbl, err := is.TableDefault(nil, true, "", []string{})
	logutil.FatalErr(err)
	err = tbl.WriteXLSX(fmt.Sprintf("backlog-%s.xlsx", dtfs), "Backlog")
	logutil.FatalErr(errorsutil.Wrap(err, "WriteXLSX"))

	fmt.Println("DONE")
}
