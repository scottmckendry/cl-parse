package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"cl-parse/changelog"
	"cl-parse/git"
)

const VERSION = "0.5.1" // x-release-please-version

type options struct {
	version          bool
	latest           bool
	release          string
	last             int
	sinceDays        int
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
		parser.OriginToken = opts.token

		if (parser.IncludeBody || parser.FetchItemDetails) && !git.IsGitRepo(".") {
			fmt.Println("Cannot fetch commits: Not a git repository")
			os.Exit(1)
		}

		entries, err := parser.Parse(string(content))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := validateScopeOptions(opts); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		filtered := filterEntries(entries, opts.last, opts.sinceDays, time.Now().UTC())

		var outputErr error
		switch {
		case opts.latest:
			outputErr = handleLatest(filtered, opts.format)
		case opts.release != "":
			outputErr = handleRelease(filtered, opts.release, opts.format)
		default:
			outputErr = outputFormatted(filtered, opts.format)
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
	cmd.Flags().Int("last", 0, "limit output to the N most recent releases")
	cmd.Flags().Int("since-days", 0, "limit output to releases within the last N days (from today, UTC)")
}

func getOptions(cmd *cobra.Command) options {
	version, _ := cmd.Flags().GetBool("version")
	latest, _ := cmd.Flags().GetBool("latest")
	release, _ := cmd.Flags().GetString("release")
	last, _ := cmd.Flags().GetInt("last")
	sinceDays, _ := cmd.Flags().GetInt("since-days")
	includeBody, _ := cmd.Flags().GetBool("include-body")
	fetchItemDetails, _ := cmd.Flags().GetBool("fetch-item-details")
	token, _ := cmd.Flags().GetString("token")
	format, _ := cmd.Flags().GetString("format")

	return options{
		version:          version,
		latest:           latest,
		release:          release,
		last:             last,
		sinceDays:        sinceDays,
		includeBody:      includeBody,
		fetchItemDetails: fetchItemDetails,
		token:            token,
		format:           format,
	}
}

func marshalWithFormat(v any, format string) ([]byte, error) {
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

func outputFormatted(data any, format string) error {
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

func filterEntries(entries []changelog.ChangelogEntry, last, sinceDays int, now time.Time) []changelog.ChangelogEntry {
	filtered := entries
	if last > 0 && last < len(filtered) {
		filtered = filtered[:last]
	}

	if sinceDays <= 0 {
		return filtered
	}

	utcNow := now.UTC()
	midnight := time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day(), 0, 0, 0, 0, time.UTC)
	cutoff := midnight.AddDate(0, 0, -sinceDays)

	var filteredByCutoff []changelog.ChangelogEntry
	for _, entry := range filtered {
		if entry.Date.Before(cutoff) {
			continue
		}
		filteredByCutoff = append(filteredByCutoff, entry)
	}

	return filteredByCutoff
}

func validateScopeOptions(o options) error {
	if o.latest && (o.release != "" || o.last > 0 || o.sinceDays > 0) {
		return fmt.Errorf("--latest cannot be combined with --release, --last, or --since-days")
	}
	if o.release != "" && (o.last > 0 || o.sinceDays > 0) {
		return fmt.Errorf("--release cannot be combined with --last or --since-days")
	}
	if o.last > 0 && o.sinceDays > 0 {
		return fmt.Errorf("--last cannot be combined with --since-days")
	}
	if o.last < 0 || o.sinceDays < 0 {
		return fmt.Errorf("--last and --since-days must be positive integers")
	}
	return nil
}
