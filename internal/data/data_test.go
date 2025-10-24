package data

import (
	"testing"
)

func TestRepoInfo(t *testing.T) {
	repo := RepoInfo{
		DatabaseId: 123,
		Name:       "test-repo",
		Visibility: "private",
		Collaborators: struct {
			Edges []Edge
		}{
			Edges: []Edge{
				{
					Permission: "admin",
					Node: struct {
						Login string
					}{Login: "testuser"},
				},
			},
		},
	}

	if repo.DatabaseId != 123 {
		t.Errorf("Expected DatabaseId 123, got %d", repo.DatabaseId)
	}

	if repo.Name != "test-repo" {
		t.Errorf("Expected Name 'test-repo', got %s", repo.Name)
	}

	if repo.Visibility != "private" {
		t.Errorf("Expected Visibility 'private', got %s", repo.Visibility)
	}

	if len(repo.Collaborators.Edges) != 1 {
		t.Errorf("Expected 1 collaborator edge, got %d", len(repo.Collaborators.Edges))
	}

	edge := repo.Collaborators.Edges[0]
	if edge.Permission != "admin" {
		t.Errorf("Expected permission 'admin', got %s", edge.Permission)
	}

	if edge.Node.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got %s", edge.Node.Login)
	}
}

func TestRepoCollaborators(t *testing.T) {
	collab := RepoCollaborators{
		Login: "testuser",
		Id:    123,
		Type:  "User",
	}

	if collab.Login != "testuser" {
		t.Errorf("Expected Login 'testuser', got %s", collab.Login)
	}

	if collab.Id != 123 {
		t.Errorf("Expected Id 123, got %d", collab.Id)
	}

	if collab.Type != "User" {
		t.Errorf("Expected Type 'User', got %s", collab.Type)
	}
}

func TestImportedRepoCollab(t *testing.T) {
	imported := ImportedRepoCollab{
		RepositoryName: "test-repo",
		Username:       "testuser",
		Permission:     "admin",
	}

	if imported.RepositoryName != "test-repo" {
		t.Errorf("Expected RepositoryName 'test-repo', got %s", imported.RepositoryName)
	}

	if imported.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got %s", imported.Username)
	}

	if imported.Permission != "admin" {
		t.Errorf("Expected Permission 'admin', got %s", imported.Permission)
	}
}

func TestPermission(t *testing.T) {
	perm := Permission{
		Permission: "write",
	}

	if perm.Permission != "write" {
		t.Errorf("Expected Permission 'write', got %s", perm.Permission)
	}
}

func TestOrganizationUserQuery(t *testing.T) {
	query := OrganizationUserQuery{
		Organization: struct {
			Repositories struct {
				Nodes    []RepoInfo
				PageInfo struct {
					EndCursor   string
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $endCursor)"`
		}{
			Repositories: struct {
				Nodes    []RepoInfo
				PageInfo struct {
					EndCursor   string
					HasNextPage bool
				}
			}{
				Nodes: []RepoInfo{
					{
						DatabaseId: 456,
						Name:       "another-repo",
						Visibility: "public",
					},
				},
				PageInfo: struct {
					EndCursor   string
					HasNextPage bool
				}{
					EndCursor:   "cursor456",
					HasNextPage: true,
				},
			},
		},
	}

	if len(query.Organization.Repositories.Nodes) != 1 {
		t.Errorf("Expected 1 repository node, got %d", len(query.Organization.Repositories.Nodes))
	}

	repo := query.Organization.Repositories.Nodes[0]
	if repo.DatabaseId != 456 {
		t.Errorf("Expected DatabaseId 456, got %d", repo.DatabaseId)
	}

	if repo.Name != "another-repo" {
		t.Errorf("Expected Name 'another-repo', got %s", repo.Name)
	}

	if query.Organization.Repositories.PageInfo.EndCursor != "cursor456" {
		t.Errorf("Expected EndCursor 'cursor456', got %s", query.Organization.Repositories.PageInfo.EndCursor)
	}

	if !query.Organization.Repositories.PageInfo.HasNextPage {
		t.Error("Expected HasNextPage to be true")
	}
}
