package main

import (
	"fmt"

	"github.com/grokify/gojira/cmd"
	"github.com/grokify/gojira/jirarest"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	opts := cmd.Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(opts)

	jrClient, err := opts.Client()
	logutil.FatalErr(errorsutil.Wrap(err, "Client"))

	iss, err := jrClient.Issue(opts.IssueKey)
	logutil.FatalErr(err)

	im := jirarest.NewIssueMore(iss)

	fmtutil.MustPrintJSON(im.Meta(""))
	fmtutil.MustPrintJSON(im.Meta(jrClient.Config.ServerURL))

	fmt.Println("DONE")
}
