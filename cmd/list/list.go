package list

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

			if _, err := os.Stat(cmdFlags.listFile); errors.Is(err, os.ErrExist) {
				return err
			}

			reportWriter, err := os.OpenFile(cmdFlags.listFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			return runCmdList(owner, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient), reportWriter)
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

func runCmdList(owner string, cmdFlags *cmdFlags, g *utils.APIGetter, reportWriter io.Writer) error {
	var reposCursor *string

	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"RepositoryName",
		"RepositoryID",
		"Visibility",
		"Username",
		"AccessLevel",
	})

	if err != nil {
		zap.S().Error("Error raised in writing output", zap.Error(err))
	}

	zap.S().Debugf("Gathering repositories and access for %s", owner)
	repoCollabList, err := g.GetOrgGuestCollaborators(owner)
	if err != nil {
		zap.S().Error("Error raised in gathering users", zap.Error(err))
	}

	var repoCollaborators []data.RepoCollaborators
	err = json.Unmarshal(repoCollabList, &repoCollaborators)
	if err != nil {
		return err
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
						zap.S().Error("Error raised in gathering repositories and user permissions", zap.Error(err))
					}
					allRepoPerms = append(allRepoPerms, repoUserPermissions.Organization.Repositories.Nodes...)
					if !repoUserPermissions.Organization.Repositories.PageInfo.HasNextPage {
						break
					}
					reposCursor = &repoUserPermissions.Organization.Repositories.PageInfo.EndCursor
				}
				for _, repo := range allRepoPerms {
					if len(repo.Collaborators.Edges) > 0 {
						err = csvWriter.Write([]string{
							repo.Name,
							strconv.Itoa(repo.DatabaseId),
							repo.Visibility,
							cmdFlags.username,
							repo.Collaborators.Edges[0].Permission,
						})
						if err != nil {
							zap.S().Error("Error raised in writing output", zap.Error(err))
						}
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
					zap.S().Error("Error raised in gathering repositories and user permissions", zap.Error(err))
				}
				allRepoPerms = append(allRepoPerms, repoUserPermissions.Organization.Repositories.Nodes...)
				if !repoUserPermissions.Organization.Repositories.PageInfo.HasNextPage {
					break
				}
				reposCursor = &repoUserPermissions.Organization.Repositories.PageInfo.EndCursor
			}
			for _, repo := range allRepoPerms {
				if len(repo.Collaborators.Edges) > 0 {
					err = csvWriter.Write([]string{
						repo.Name,
						strconv.Itoa(repo.DatabaseId),
						repo.Visibility,
						repoCollab.Login,
						repo.Collaborators.Edges[0].Permission,
					})
					if err != nil {
						zap.S().Error("Error raised in writing output", zap.Error(err))
					}
				}
			}
		}
	}

	fmt.Printf("Successfully listed repository collaborator permissions for repositories in %s", owner)
	csvWriter.Flush()

	return nil
}
