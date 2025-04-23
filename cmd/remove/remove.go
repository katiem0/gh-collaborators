package add

import (
	"encoding/csv"
	"fmt"
	"os"

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
	fileName string
	debug    bool
}

func NewCmdRemove() *cobra.Command {
	cmdFlags := cmdFlags{}
	var authToken string

	removeCmd := &cobra.Command{
		Use:   "remove [flags] <organization>",
		Short: "Remove repo access for repository collaborators.",
		Long:  "Remove repositories and permissions for repository collaborators.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(removeCmd *cobra.Command, args []string) error {
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

			return runCmdRemove(owner, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient))
		},
	}

	// Configure flags for command

	removeCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub Personal Access Token (default "gh auth token")`)
	removeCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	removeCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to remove access from (required)")
	removeCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	removeCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	err := removeCmd.MarkFlagRequired("from-file")
	if err != nil {
		zap.S().Errorf("Error marking flag 'from-file' as required: %v", err)
	}

	return removeCmd
}

func runCmdRemove(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var collabData [][]string
	var importRepoCollabList []data.ImportedRepoCollab

	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening repository collaborators csv file")
		}
		// remember to close the file at the end of the program
		defer func() {
			closeErr := f.Close()
			if closeErr != nil {
				zap.S().Warnf("Error closing file: %v", closeErr)
			}
		}()
		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		collabData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading collaborators to remove from csv file")
		}
		importRepoCollabList = g.DeleteRepoCollaboratorsList(collabData)
	} else {
		zap.S().Errorf("Error arose identifying users to add")
	}
	zap.S().Debugf("Determining users to remove")
	for _, importRepoCollab := range importRepoCollabList {
		zap.S().Debugf("Removing Repository Assignment for %s", importRepoCollab.Username)

		err := g.RemoveRepoCollaborator(owner, importRepoCollab.RepositoryName, importRepoCollab.Username)
		if err != nil {
			zap.S().Errorf("Error arose removing permission for user %s  and repo %s", importRepoCollab.Username, importRepoCollab.RepositoryName)
		}
	}

	fmt.Printf("Successfully removed repository assignments for repository collaborators in: %s.", owner)
	return nil
}
