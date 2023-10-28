package cmd

import (
	"github.com/grokify/gojira/jirarest"
)

type Options struct {
	Authfile        string `short:"a" long:"Goauth authfile" description:"Goauth auth File"`
	Authkey         string `short:"k" long:"Goauth key" description:"Goauth credentials Key"`
	BoardID         uint   `short:"b" long:"boardid" description:"Jira Board ID"`
	IssueKey        string `short:"i" long:"key" description:"Jira Issue Key"`
	Customfield     string `short:"c" long:"customfield" description:"Custom field"`
	CustomfieldName string `short:"n" long:"customfield name" description:"Custom field name"` // 'Epic Link'
}

func (opts Options) Client() (*jirarest.Client, error) {
	return jirarest.ClientsBasicAuthFile(opts.Authfile, opts.Authkey)
}
