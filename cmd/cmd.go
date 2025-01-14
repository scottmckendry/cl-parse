package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"cl-parse/changelog"
	"cl-parse/git"
)

const VERSION = "0.3.0" // x-release-please-version

type options struct {
	version          bool
	latest           bool
	release          string
	includeBody      bool
	fetchItemDetails bool
	token            string
	format           string
}

var cmd = &cobra.Command{
	Use:  "cl-parse [flags] [path]",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		changelogPath := "./CHANGELOG.md"
		if len(args) > 0 {
			changelogPath = args[0]
		}

		opts := getOptions(cmd)

		if opts.version {
			fmt.Printf("cl-parse v%s\n", VERSION)
			os.Exit(0)
		}

		content, err := os.ReadFile(changelogPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parser := changelog.NewParser()
		parser.IncludeBody = opts.includeBody
		parser.FetchItemDetails = opts.fetchItemDetails
		if (parser.IncludeBody || parser.FetchItemDetails) && !git.IsGitRepo(".") {
			fmt.Println("Cannot fetch commits: Not a git repository")
			os.Exit(1)
		}

		entries, err := parser.Parse(string(content))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var outputErr error
		switch {
		case opts.latest:
			outputErr = handleLatest(entries, opts.format)
		case opts.release != "":
			outputErr = handleRelease(entries, opts.release, opts.format)
		default:
			outputErr = outputFormatted(entries, opts.format)
		}

		if outputErr != nil {
			fmt.Println(outputErr)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cmd.Flags().BoolP("version", "v", false, "display the current version of cl-parse")
	cmd.Flags().BoolP("latest", "l", false, "display the most recent version from the changelog")
	cmd.Flags().StringP("release", "r", "", "display the changelog entry for a specific release")
	cmd.Flags().Bool("include-body", false, "include the full commit body in changelog entry")
	cmd.Flags().
		Bool("fetch-item-details", false, "fetch details for related items (e.g. GitHub issues & PRs)")
	cmd.Flags().String("token", "", "token for fetching related items")
	cmd.Flags().StringP("format", "f", "json", "output format (json, yaml, or toml)")
}

func getOptions(cmd *cobra.Command) options {
	version, _ := cmd.Flags().GetBool("version")
	latest, _ := cmd.Flags().GetBool("latest")
	release, _ := cmd.Flags().GetString("release")
	includeBody, _ := cmd.Flags().GetBool("include-body")
	fetchItemDetails, _ := cmd.Flags().GetBool("fetch-item-details")
	token, _ := cmd.Flags().GetString("token")
	format, _ := cmd.Flags().GetString("format")

	return options{
		version:          version,
		latest:           latest,
		release:          release,
		includeBody:      includeBody,
		fetchItemDetails: fetchItemDetails,
		token:            token,
		format:           format,
	}
}

func marshalWithFormat(v interface{}, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(v, "", "  ")
	case "yaml":
		return yaml.Marshal(v)
	case "toml":
		return toml.Marshal(v)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func outputFormatted(data interface{}, format string) error {
	output, err := marshalWithFormat(data, format)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func handleLatest(entries []changelog.ChangelogEntry, format string) error {
	if len(entries) == 0 {
		return fmt.Errorf("no changelog entries found")
	}
	return outputFormatted(entries[0], format)
}

func handleRelease(entries []changelog.ChangelogEntry, release, format string) error {
	for _, entry := range entries {
		if entry.Version == release {
			return outputFormatted(entry, format)
		}
	}
	return fmt.Errorf("version %s not found in changelog", release)
}
