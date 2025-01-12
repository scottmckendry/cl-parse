package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"cl-parse/changelog"
)

const VERSION = "0.2.0" // x-release-please-version

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
			jsonData, err := json.MarshalIndent(entries[0], "", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
			return
		}

		if release != "" {
			found := false
			for _, entry := range entries {
				if entry.Version == release {
					jsonData, err := json.MarshalIndent(entry, "", "  ")
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println(string(jsonData))
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
		jsonData, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
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
}
