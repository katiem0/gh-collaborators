package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/katiem0/gh-collaborators/internal/data"
	"github.com/katiem0/gh-collaborators/internal/log"
	"github.com/katiem0/gh-collaborators/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	token    string
	hostname string
	listFile string
	username string
	debug    bool
}

func NewCmdList() *cobra.Command {
	cmdFlags := cmdFlags{}
	var authToken string

	listCmd := &cobra.Command{
		Use:   "list [flags] <organization>",
		Short: "Generate a report of repos that repository collaborators have access to.",
		Long:  "Generate a report of repos that repository collaborators have access to.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(listCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient *api.GraphQLClient
			var restClient *api.RESTClient

			// Reinitialize logging if debugging was enabled
			if cmdFlags.debug {
				logger, _ := log.NewLogger(cmdFlags.debug)
				defer logger.Sync() // nolint:errcheck
				zap.ReplaceGlobals(logger)
			}

			if cmdFlags.token != "" {
				authToken = cmdFlags.token
			} else {
				t, _ := auth.TokenForHost(cmdFlags.hostname)
				authToken = t
			}

			restClient, err = api.NewRESTClient(api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving rest client")
				return err
			}

			gqlClient, err = api.NewGraphQLClient(api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github.hawkgirl-preview+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving graphql client")
				return err
			}

			owner := args[0]

			// Check if file exists, but don't fail if it doesn't
			if _, err := os.Stat(cmdFlags.listFile); err == nil {
				return fmt.Errorf("output file %s already exists", cmdFlags.listFile)
			}

			// Create APIGetter
			apiGetter := utils.NewAPIGetter(gqlClient, restClient)

			// Collect all data first, don't create file yet
			if err := runCmdList(owner, &cmdFlags, apiGetter); err != nil {
				return err
			}
			return nil
		},
	}

	reportFileDefault := fmt.Sprintf("RepoCollaboratorsReport-%s.csv", time.Now().Format("20060102150405"))

	// Configure flags for command
	listCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub Personal Access Token (default "gh auth token")`)
	listCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	listCmd.Flags().StringVarP(&cmdFlags.listFile, "output-file", "o", reportFileDefault, "Name of file to write CSV list to")
	listCmd.PersistentFlags().StringVarP(&cmdFlags.username, "username", "u", "", "Username of single repo collaborator to generate report for")
	listCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")

	return listCmd
}

func runCmdList(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var reposCursor *string
	var csvData [][]string

	// Add header to data slice
	csvData = append(csvData, []string{
		"RepositoryName",
		"RepositoryID",
		"Visibility",
		"Username",
		"AccessLevel",
	})

	zap.S().Debugf("Gathering repositories and access for %s", owner)
	repoCollabList, err := g.GetOrgGuestCollaborators(owner)
	if err != nil {
		zap.S().Errorf("Failed to get organization collaborators for '%s'", owner)
		return err
	}

	var repoCollaborators []data.RepoCollaborators
	err = json.Unmarshal(repoCollabList, &repoCollaborators)
	if err != nil {
		zap.S().Errorf("Failed to parse collaborators data: %v", err)
		return fmt.Errorf("failed to parse collaborators data: %w", err)
	}

	if len(cmdFlags.username) > 0 {
		zap.S().Debugf("Checking if username %s is in list of repository collaborators", cmdFlags.username)
		for _, repoCollab := range repoCollaborators {
			if cmdFlags.username == repoCollab.Login {
				zap.S().Debugf("Gathering repositories for specified username %s", cmdFlags.username)
				var allRepoPerms []data.RepoInfo
				for {
					repoUserPermissions, err := g.GetOrgRepositoryPermissions(owner, cmdFlags.username, reposCursor)
					if err != nil {
						zap.S().Errorf("Failed to get repository permissions for user '%s' in organization '%s': %v", cmdFlags.username, owner, err)
						return fmt.Errorf("failed to get repository permissions for user %s: %w", cmdFlags.username, err)
					}

					allRepoPerms = append(allRepoPerms, repoUserPermissions.Organization.Repositories.Nodes...)
					if !repoUserPermissions.Organization.Repositories.PageInfo.HasNextPage {
						break
					}
					reposCursor = &repoUserPermissions.Organization.Repositories.PageInfo.EndCursor
				}
				for _, repo := range allRepoPerms {
					if len(repo.Collaborators.Edges) > 0 {
						csvData = append(csvData, []string{
							repo.Name,
							strconv.Itoa(repo.DatabaseId),
							repo.Visibility,
							cmdFlags.username,
							repo.Collaborators.Edges[0].Permission,
						})
					}
				}
			}
		}
	} else {
		for _, repoCollab := range repoCollaborators {
			zap.S().Debugf("Gathering repositories for username %s", repoCollab.Login)
			var allRepoPerms []data.RepoInfo
			for {
				repoUserPermissions, err := g.GetOrgRepositoryPermissions(owner, repoCollab.Login, reposCursor)
				if err != nil {
					zap.S().Errorf("Failed to get repository permissions for user '%s' in organization '%s': %v", repoCollab.Login, owner, err)
					return fmt.Errorf("failed to get repository permissions for user %s: %w", repoCollab.Login, err)
				}

				allRepoPerms = append(allRepoPerms, repoUserPermissions.Organization.Repositories.Nodes...)
				if !repoUserPermissions.Organization.Repositories.PageInfo.HasNextPage {
					break
				}
				reposCursor = &repoUserPermissions.Organization.Repositories.PageInfo.EndCursor
			}
			for _, repo := range allRepoPerms {
				if len(repo.Collaborators.Edges) > 0 {
					csvData = append(csvData, []string{
						repo.Name,
						strconv.Itoa(repo.DatabaseId),
						repo.Visibility,
						repoCollab.Login,
						repo.Collaborators.Edges[0].Permission,
					})
				}
			}
		}
	}

	// Only create and write to file after all data is successfully collected
	if len(csvData) <= 1 { // Only header, no actual data
		return fmt.Errorf("no collaborator data found for organization %s", owner)
	}

	zap.S().Debugf("Creating output file %s", cmdFlags.listFile)
	reportWriter, err := os.OpenFile(cmdFlags.listFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		closeErr := reportWriter.Close()
		if closeErr != nil {
			zap.S().Warnf("Error closing file: %v", closeErr)
		}
	}()

	csvWriter := csv.NewWriter(reportWriter)
	defer csvWriter.Flush()

	// Write all collected data to CSV
	for _, row := range csvData {
		err = csvWriter.Write(row)
		if err != nil {
			zap.S().Error("Error raised in writing output", zap.Error(err))
			return fmt.Errorf("failed to write CSV data: %w", err)
		}
	}

	fmt.Printf("Successfully listed repository collaborator permissions for repositories in %s\n", owner)
	fmt.Printf("Report saved to: %s\n", cmdFlags.listFile)

	return nil
}
