package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grokify/gojira/cmd"
	"github.com/grokify/gojira/jirarest"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	opts := cmd.Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	jrClient, err := opts.Client()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	cfs, err := jrClient.CustomFieldAPI.GetCustomFields()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(3)
	}

	err = cfs.WriteTable(os.Stdout)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(4)
	}

	if opts.Customfield != "" {
		ids := strings.Split(opts.Customfield, ",")
		filtered := cfs.FilterByIDs(ids...)
		if err := filtered.WriteTable(os.Stdout); err != nil {
			slog.Error(err.Error())
			os.Exit(5)
		}
	}

	if opts.CustomfieldName != "" {
		names := strings.Split(opts.CustomfieldName, ",")
		filtered := cfs.FilterByNames(names...)
		if err := filtered.WriteTable(os.Stdout); err != nil {
			slog.Error(err.Error())
			os.Exit(6)
		}
	}

	// Get Epic Link Custom Field
	cfName, err := jrClient.CustomFieldAPI.GetCustomFieldEpicLink()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(7)
	}

	cfsName := jirarest.CustomFields{cfName}
	if err := cfsName.WriteTable(os.Stdout); err != nil {
		slog.Error(err.Error())
		os.Exit(8)
	}

	slog.Info("DONE")
	os.Exit(0)
}
