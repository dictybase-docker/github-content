package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase-docker/github-content/internal/logger"
	"github.com/dictyBase-docker/github-content/internal/registry"
	"github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "github-content",
	Short: "cli to download modified files from github commit",
	Long: `A command line application that extract the list of modified files
	from the latest commit and then download them using github api.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		l, err := logger.NewLogger(cmd)
		if err != nil {
			return err
		}
		registry.SetLogger(l)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		d, _ := cmd.Flags().GetBool("doc")
		if d {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			docDir := filepath.Join(dir, "docs")
			if err := os.MkdirAll(docDir, 0700); err != nil {
				return err
			}
			if err := doc.GenMarkdownTree(cmd, docDir); err != nil {
				return err
			}
			fmt.Printf("created markdown docs in %s\n", docDir)
			return nil
		}
		l := registry.GetLogger()
		cp, _ := cmd.Flags().GetString("commit-payload")
		folder, _ := cmd.Flags().GetString("folder")
		p, _ := cmd.Flags().GetString("file-extension")
		owner, _ := cmd.Flags().GetString("owner")
		repo, _ := cmd.Flags().GetString("repository")
		var commits []*github.WebHookCommit
		if err := json.Unmarshal([]byte(cp), commits); err != nil {
			return fmt.Errorf("error in decoding commit payload %s", err)
		}
		client := github.NewClient(nil)
		for _, c := range commits {
			for _, m := range c.Modified {
				if !strings.HasSuffix(m, p) {
					l.Debugf("skipped file %s from downloading", m)
					continue
				}
				fcont, _, _, err := client.Repositories.GetContents(
					context.Background(), owner, repo,
					m, &github.RepositoryContentGetOptions{Ref: *c.ID},
				)
				str, err := fcont.GetContent()
				if err != nil {
					return fmt.Errorf("error in decoding github file content %s", err)
				}
				fname := filepath.Join(folder, filepath.Base(m))
				if err := ioutil.WriteFile(fname, []byte(str), 0644); err != nil {
					return fmt.Errorf("error in writing file %s %s", fname, err)
				}
				l.Infof("written file %s", fname)
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().Bool("doc", false, "generate markdown documentation")
	RootCmd.Flags().StringP(
		"commit-payload",
		"c",
		"",
		"commit data that is received from GitHub after triggered by a push event[required]",
	)
	RootCmd.Flags().StringP(
		"owner",
		"o",
		"",
		"github repository owner[required]",
	)
	RootCmd.Flags().StringP(
		"repository",
		"r",
		"",
		"github repository name[required]",
	)
	RootCmd.Flags().StringP(
		"folder",
		"f",
		"",
		"output folder[required]",
	)
	RootCmd.MarkFlagRequired("folder")
	RootCmd.MarkFlagRequired("commit-payload")
	RootCmd.MarkFlagRequired("repository")
	RootCmd.MarkFlagRequired("owner")
	RootCmd.Flags().StringP(
		"file-extension",
		"p",
		"obo",
		"file extension that will be screened in the commit payload",
	)
	RootCmd.Flags().StringP(
		"log-level",
		"",
		"error",
		"log level for the application",
	)
	RootCmd.Flags().StringP(
		"log-format",
		"",
		"json",
		"format of the logging out, either of json or text",
	)
	RootCmd.Flags().String(
		"log-file",
		"",
		"file for log output other than standard output, written to a temp folder by default",
	)
}
