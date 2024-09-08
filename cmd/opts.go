package cmd

import (
	"github.com/grokify/goauth"
	"github.com/grokify/gojira/jirarest"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Options goauth.Options
	// Authfile        string `short:"a" long:"Goauth authfile" description:"Goauth auth File"`
	// Authkey         string `short:"k" long:"Goauth key" description:"Goauth credentials Key"`
	BoardID         uint   `short:"b" long:"boardid" description:"Jira Board ID"`
	IssueKey        string `short:"i" long:"key" description:"Jira Issue Key"`
	JQL             string `short:"j" long:"jql" description:"Jira Query Language"`
	Customfield     string `short:"c" long:"customfield" description:"Custom field"`
	CustomfieldName string `short:"n" long:"customfield name" description:"Custom field name"` // 'Epic Link'
}

func (opts Options) Client() (*jirarest.Client, error) {
	return jirarest.NewClientGoauthBasicAuthFile(opts.Options.CredsPath, opts.Options.Account)
}

func NewOptions() (Options, error) {
	opts := Options{}
	_, err := flags.Parse(&opts)
	return opts, err
}
