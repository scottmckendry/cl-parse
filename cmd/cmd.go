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

var cmd = &cobra.Command{
	Use:  "cl-parse [flags] [path]",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		changelogPath := "./CHANGELOG.md"
		if len(args) > 0 {
			changelogPath = args[0]
		}

		ver, _ := cmd.Flags().GetBool("version")
		latest, _ := cmd.Flags().GetBool("latest")
		release, _ := cmd.Flags().GetString("release")
		includeBody, _ := cmd.Flags().GetBool("include-body")
		format, _ := cmd.Flags().GetString("format")

		if ver {
			fmt.Printf("cl-parse v%s\n", VERSION)
			os.Exit(0)
		}

		content, err := os.ReadFile(changelogPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parser := changelog.NewParser()
		parser.IncludeBody = includeBody
		if parser.IncludeBody {
			if !git.IsGitRepo(".") {
				fmt.Println("Cannot fetch commits: Not a git repository")
				os.Exit(1)
			}
		}

		entries, err := parser.Parse(string(content))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if latest {
			if len(entries) == 0 {
				fmt.Println("No changelog entries found")
				os.Exit(1)
			}
			outputData, err := marshalWithFormat(entries[0], format)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(outputData))
			return
		}

		if release != "" {
			found := false
			for _, entry := range entries {
				if entry.Version == release {
					outputData, err := marshalWithFormat(entry, format)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println(string(outputData))
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Version %s not found in changelog\n", release)
				os.Exit(1)
			}
			return
		}

		// default to printing all entries
		outputData, err := marshalWithFormat(entries, format)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(outputData))
	},
}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cmd.Flags().BoolP("version", "v", false, "display the current version of cl-parse")
	cmd.Flags().BoolP("latest", "l", false, "display the most recent version from the changelog")
	cmd.Flags().StringP("release", "r", "", "display the changelog entry for a specific release")
	cmd.Flags().Bool("include-body", false, "include the full commit body in changelog entry")
	cmd.Flags().StringP("format", "f", "json", "output format (json, yaml, or toml)")
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
