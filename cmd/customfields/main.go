package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/gojira/cmd"
	"github.com/grokify/gojira/jirarest"
	"github.com/grokify/mogo/log/logutil"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	opts := cmd.Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	jrClient, err := opts.Client()
	logutil.FatalErr(err)

	cfs, err := jirarest.GetCustomFields(jrClient)
	logutil.FatalErr(err)
	cfs.WriteTable(os.Stdout)

	if opts.Customfield != "" {
		ids := strings.Split(opts.Customfield, ",")
		filtered := cfs.FilterByIDs(ids...)
		filtered.WriteTable(os.Stdout)
	}

	if opts.CustomfieldName != "" {
		names := strings.Split(opts.CustomfieldName, ",")
		filtered := cfs.FilterByNames(names...)
		filtered.WriteTable(os.Stdout)
	}

	// Get Epic Link Custom Field
	cfName, err := jirarest.GetCustomFieldEpicLink(jrClient)
	logutil.FatalErr(err)
	cfsName := jirarest.CustomFields{cfName}
	cfsName.WriteTable(os.Stdout)

	fmt.Println("DONE")
}
