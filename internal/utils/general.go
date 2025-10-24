package utils

import (
	"fmt"
	"io"
	"log"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/katiem0/gh-collaborators/internal/data"
	"github.com/shurcooL/graphql"
	"go.uber.org/zap"
)

type Getter interface {
	AddRepoCollaborator(owner string, repo string, username string, data io.Reader) error
	CreateRepoCollaboratorsList(filedata [][]string) []data.ImportedRepoCollab
	CreateRepoPermData(permission string) *data.Permission
	GetOrgGuestCollaborators(owner string) ([]byte, error)
	GetOrgRepositoryPermissions(owner string, user string, endCursor *string) (*data.OrganizationUserQuery, error)
	RemoveRepoCollaborator(owner string, repo string, username string) error
}

type APIGetter struct {
	gqlClient  api.GraphQLClient
	restClient api.RESTClient
}

func NewAPIGetter(gqlClient *api.GraphQLClient, restClient *api.RESTClient) *APIGetter {
	getter := &APIGetter{}

	if gqlClient != nil {
		getter.gqlClient = *gqlClient
	}

	if restClient != nil {
		getter.restClient = *restClient
	}

	return getter
}

func (g *APIGetter) GetOrgGuestCollaborators(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/outside_collaborators", owner)
	zap.S().Debugf("Reading in repository collaborators from %v", url)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			zap.S().Warnf("Error closing response body: %v", closeErr)
		}
	}()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	return responseData, err
}

func (g *APIGetter) GetOrgRepositoryPermissions(owner string, user string, endCursor *string) (*data.OrganizationUserQuery, error) {
	query := new(data.OrganizationUserQuery)
	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(endCursor),
		"owner":     graphql.String(owner),
		"user":      graphql.String(user),
	}
	err := g.gqlClient.Query("getOrganizationRepoPermissions", &query, variables)

	return query, err
}

func (g *APIGetter) CreateRepoCollaboratorsList(filedata [][]string) []data.ImportedRepoCollab {
	//convert csv lines to array of structs
	var importRepoCollabs []data.ImportedRepoCollab
	var repoCollab data.ImportedRepoCollab
	for _, each := range filedata[1:] {
		if len(each) > 0 {
			repoCollab.RepositoryName = each[0]
		} else {
			repoCollab.RepositoryName = ""
		}

		if len(each) > 1 {
			repoCollab.Username = each[1]
		} else {
			repoCollab.Username = ""
		}

		if len(each) > 2 {
			repoCollab.Permission = each[2]
		} else {
			repoCollab.Permission = ""
		}

		importRepoCollabs = append(importRepoCollabs, repoCollab)
	}
	return importRepoCollabs
}

func (g *APIGetter) AddRepoCollaborator(owner string, repo string, username string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/collaborators/%s", owner, repo, username)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			zap.S().Warnf("Error closing response body: %v", closeErr)
		}
	}()
	return err
}

func (g *APIGetter) CreateRepoPermData(permission string) *data.Permission {
	s := data.Permission{
		Permission: permission,
	}
	return &s
}

func (g *APIGetter) DeleteRepoCollaboratorsList(filedata [][]string) []data.ImportedRepoCollab {
	//convert csv lines to array of structs
	var importRepoCollabs []data.ImportedRepoCollab
	var repoCollab data.ImportedRepoCollab
	for _, each := range filedata[1:] {
		if len(each) > 0 {
			repoCollab.RepositoryName = each[0]
		} else {
			repoCollab.RepositoryName = ""
		}

		if len(each) > 1 {
			repoCollab.Username = each[1]
		} else {
			repoCollab.Username = ""
		}

		importRepoCollabs = append(importRepoCollabs, repoCollab)
	}
	return importRepoCollabs
}

func (g *APIGetter) RemoveRepoCollaborator(owner string, repo string, username string) error {
	url := fmt.Sprintf("repos/%s/%s/collaborators/%s", owner, repo, username)

	resp, err := g.restClient.Request("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			zap.S().Warnf("Error closing response body: %v", closeErr)
		}
	}()
	return err
}
