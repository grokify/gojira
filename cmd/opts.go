package cmd

import (
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira/jirarest"
)

type Options struct {
	Authfile        string `short:"a" long:"Goauth authfile" description:"Goauth auth File"`
	Authkey         string `short:"k" long:"Goauth key" description:"Goauth credentials Key"`
	Customfield     string `short:"c" long:"customfield" description:"Custom field"`
	CustomfieldName string `short:"n" long:"customfield name" description:"Custom field name"` // 'Epic Link'
}

func (opts Options) Clients() (*http.Client, *jira.Client, string, error) {
	return jirarest.ClientsBasicAuthFile(opts.Authfile, opts.Authkey)
}
