package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information. These can be set at build time using ldflags:
// go build -ldflags "-X main.Version=v1.0.0 -X main.Commit=abc123 -X main.BuildDate=2024-01-01"
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show version, commit, build date, and Go runtime information.`,
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

type versionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

func runVersion(cmd *cobra.Command, args []string) error {
	info := versionInfo{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	format := getOutputFormat()
	if format == OutputJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	// Table or TOON format - use simple text output
	fmt.Fprintf(cmd.OutOrStdout(), "gojira %s\n", info.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "  Commit:     %s\n", info.Commit)
	fmt.Fprintf(cmd.OutOrStdout(), "  Build Date: %s\n", info.BuildDate)
	fmt.Fprintf(cmd.OutOrStdout(), "  Go Version: %s\n", info.GoVersion)
	fmt.Fprintf(cmd.OutOrStdout(), "  OS/Arch:    %s/%s\n", info.OS, info.Arch)
	return nil
}
