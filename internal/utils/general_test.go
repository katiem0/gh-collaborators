package utils

import (
	"fmt"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/katiem0/gh-collaborators/internal/data"
)

func TestCreateRepoCollaboratorsList(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username", "Permission"}, // header
		{"repo1", "user1", "read"},
		{"repo2", "user2", "write"},
		{"repo3", "user3", "admin"},
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 3 {
		t.Errorf("Expected 3 collaborators, got %d", len(result))
	}

	expected := []data.ImportedRepoCollab{
		{RepositoryName: "repo1", Username: "user1", Permission: "read"},
		{RepositoryName: "repo2", Username: "user2", Permission: "write"},
		{RepositoryName: "repo3", Username: "user3", Permission: "admin"},
	}

	for i, collab := range result {
		if collab.RepositoryName != expected[i].RepositoryName {
			t.Errorf("Expected repository name %s, got %s", expected[i].RepositoryName, collab.RepositoryName)
		}
		if collab.Username != expected[i].Username {
			t.Errorf("Expected username %s, got %s", expected[i].Username, collab.Username)
		}
		if collab.Permission != expected[i].Permission {
			t.Errorf("Expected permission %s, got %s", expected[i].Permission, collab.Permission)
		}
	}
}

func TestDeleteRepoCollaboratorsList(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username"}, // header
		{"repo1", "user1"},
		{"repo2", "user2"},
	}

	result := getter.DeleteRepoCollaboratorsList(testData)

	if len(result) != 2 {
		t.Errorf("Expected 2 collaborators, got %d", len(result))
	}

	expected := []data.ImportedRepoCollab{
		{RepositoryName: "repo1", Username: "user1"},
		{RepositoryName: "repo2", Username: "user2"},
	}

	for i, collab := range result {
		if collab.RepositoryName != expected[i].RepositoryName {
			t.Errorf("Expected repository name %s, got %s", expected[i].RepositoryName, collab.RepositoryName)
		}
		if collab.Username != expected[i].Username {
			t.Errorf("Expected username %s, got %s", expected[i].Username, collab.Username)
		}
	}
}

func TestCreateRepoPermData(t *testing.T) {
	getter := &APIGetter{}
	permission := "admin"
	result := getter.CreateRepoPermData(permission)

	if result == nil {
		t.Fatal("Expected Permission struct, got nil")
	}

	if result.Permission != permission {
		t.Errorf("Expected permission %s, got %s", permission, result.Permission)
	}
}

func TestCreateRepoCollaboratorsListEmptyData(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username", "Permission"}, // header only
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 0 {
		t.Errorf("Expected 0 collaborators, got %d", len(result))
	}
}

func TestDeleteRepoCollaboratorsListEmptyData(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username"}, // header only
	}

	result := getter.DeleteRepoCollaboratorsList(testData)

	if len(result) != 0 {
		t.Errorf("Expected 0 collaborators, got %d", len(result))
	}
}

func TestCreateRepoCollaboratorsListMultipleEntries(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username", "Permission"},
		{"repo1", "user1", "read"},
		{"repo2", "user2", "write"},
		{"repo3", "user3", "admin"},
		{"repo4", "user4", "maintain"},
		{"repo5", "user5", "triage"},
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 5 {
		t.Errorf("Expected 5 collaborators, got %d", len(result))
	}

	// Verify last entry
	lastEntry := result[4]
	if lastEntry.RepositoryName != "repo5" {
		t.Errorf("Expected last repo to be 'repo5', got %s", lastEntry.RepositoryName)
	}
	if lastEntry.Username != "user5" {
		t.Errorf("Expected last username to be 'user5', got %s", lastEntry.Username)
	}
	if lastEntry.Permission != "triage" {
		t.Errorf("Expected last permission to be 'triage', got %s", lastEntry.Permission)
	}
}

func TestCreateRepoPermDataDifferentPermissions(t *testing.T) {
	getter := &APIGetter{}
	permissions := []string{"read", "write", "admin", "maintain", "triage"}

	for _, perm := range permissions {
		result := getter.CreateRepoPermData(perm)

		if result == nil {
			t.Fatalf("Expected Permission struct for %s, got nil", perm)
		}

		if result.Permission != perm {
			t.Errorf("Expected permission %s, got %s", perm, result.Permission)
		}
	}
}

func TestAPIGetterStructure(t *testing.T) {
	// Test that APIGetter struct exists and can be instantiated
	getter := &APIGetter{}

	if getter == nil {
		t.Fatal("Expected APIGetter to be created, got nil")
	}
}

func TestNewAPIGetter(t *testing.T) {
	// Create mock clients - these can be nil for this test since we're just testing the constructor
	var gqlClient *api.GraphQLClient
	var restClient *api.RESTClient

	getter := NewAPIGetter(gqlClient, restClient)

	if getter == nil {
		t.Fatal("Expected APIGetter to be created, got nil")
	}
}

func TestCreateRepoCollaboratorsListInvalidData(t *testing.T) {
	getter := &APIGetter{}

	// Test with incomplete rows
	testData := [][]string{
		{"RepositoryName", "Username", "Permission"},
		{"repo1", "user1"}, // Missing permission
		{"repo2"},          // Missing username and permission
		{},                 // Empty row
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	// Should handle invalid data gracefully
	if len(result) != 3 {
		t.Errorf("Expected 3 collaborators (including invalid ones), got %d", len(result))
	}
}

func TestDeleteRepoCollaboratorsListInvalidData(t *testing.T) {
	getter := &APIGetter{}

	// Test with incomplete rows
	testData := [][]string{
		{"RepositoryName", "Username"},
		{"repo1"},          // Missing username
		{},                 // Empty row
		{"repo2", "user2"}, // Valid row
	}

	result := getter.DeleteRepoCollaboratorsList(testData)

	// Should handle invalid data gracefully
	if len(result) != 3 {
		t.Errorf("Expected 3 collaborators (including invalid ones), got %d", len(result))
	}
}

func TestCreateRepoCollaboratorsListSingleEntry(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username", "Permission"},
		{"single-repo", "single-user", "write"},
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 1 {
		t.Errorf("Expected 1 collaborator, got %d", len(result))
	}

	expected := data.ImportedRepoCollab{
		RepositoryName: "single-repo",
		Username:       "single-user",
		Permission:     "write",
	}

	if result[0].RepositoryName != expected.RepositoryName {
		t.Errorf("Expected repository name %s, got %s", expected.RepositoryName, result[0].RepositoryName)
	}
	if result[0].Username != expected.Username {
		t.Errorf("Expected username %s, got %s", expected.Username, result[0].Username)
	}
	if result[0].Permission != expected.Permission {
		t.Errorf("Expected permission %s, got %s", expected.Permission, result[0].Permission)
	}
}

func TestDeleteRepoCollaboratorsListSingleEntry(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username"},
		{"single-repo", "single-user"},
	}

	result := getter.DeleteRepoCollaboratorsList(testData)

	if len(result) != 1 {
		t.Errorf("Expected 1 collaborator, got %d", len(result))
	}

	expected := data.ImportedRepoCollab{
		RepositoryName: "single-repo",
		Username:       "single-user",
	}

	if result[0].RepositoryName != expected.RepositoryName {
		t.Errorf("Expected repository name %s, got %s", expected.RepositoryName, result[0].RepositoryName)
	}
	if result[0].Username != expected.Username {
		t.Errorf("Expected username %s, got %s", expected.Username, result[0].Username)
	}
}

func TestCreateRepoPermDataEmptyPermission(t *testing.T) {
	getter := &APIGetter{}
	result := getter.CreateRepoPermData("")

	if result == nil {
		t.Fatal("Expected Permission struct, got nil")
	}

	if result.Permission != "" {
		t.Errorf("Expected empty permission, got %s", result.Permission)
	}
}

func TestCreateRepoPermDataSpecialCharacters(t *testing.T) {
	getter := &APIGetter{}
	specialPermissions := []string{
		"permission-with-dashes",
		"permission_with_underscores",
		"permission with spaces",
		"UPPERCASE",
		"123numeric",
		"special@#$%chars",
	}

	for _, perm := range specialPermissions {
		result := getter.CreateRepoPermData(perm)

		if result == nil {
			t.Fatalf("Expected Permission struct for %s, got nil", perm)
		}

		if result.Permission != perm {
			t.Errorf("Expected permission %s, got %s", perm, result.Permission)
		}
	}
}

func TestCreateRepoCollaboratorsListDuplicateEntries(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username", "Permission"},
		{"repo1", "user1", "read"},
		{"repo1", "user1", "read"},  // Duplicate
		{"repo1", "user1", "write"}, // Same repo/user, different permission
		{"repo2", "user2", "admin"},
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 4 {
		t.Errorf("Expected 4 collaborators (including duplicates), got %d", len(result))
	}

	// Verify that duplicates are preserved
	if result[0].RepositoryName != "repo1" || result[0].Username != "user1" || result[0].Permission != "read" {
		t.Error("First entry incorrect")
	}
	if result[1].RepositoryName != "repo1" || result[1].Username != "user1" || result[1].Permission != "read" {
		t.Error("Second entry (duplicate) incorrect")
	}
	if result[2].RepositoryName != "repo1" || result[2].Username != "user1" || result[2].Permission != "write" {
		t.Error("Third entry (different permission) incorrect")
	}
}

func TestDeleteRepoCollaboratorsListDuplicateEntries(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username"},
		{"repo1", "user1"},
		{"repo1", "user1"}, // Duplicate
		{"repo2", "user2"},
	}

	result := getter.DeleteRepoCollaboratorsList(testData)

	if len(result) != 3 {
		t.Errorf("Expected 3 collaborators (including duplicates), got %d", len(result))
	}
}

func TestCreateRepoCollaboratorsListHeaderVariations(t *testing.T) {
	getter := &APIGetter{}

	// Test with different header case/format
	testData := [][]string{
		{"Repository", "User", "Role"}, // Different header names
		{"repo1", "user1", "read"},
		{"repo2", "user2", "write"},
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	// Should still work (function uses index-based access)
	if len(result) != 2 {
		t.Errorf("Expected 2 collaborators, got %d", len(result))
	}
}

func TestCreateRepoCollaboratorsListLargeDataset(t *testing.T) {
	getter := &APIGetter{}

	// Create a larger dataset to test performance/handling
	testData := [][]string{
		{"RepositoryName", "Username", "Permission"},
	}

	// Add 100 entries
	for i := 0; i < 100; i++ {
		testData = append(testData, []string{
			fmt.Sprintf("repo%d", i),
			fmt.Sprintf("user%d", i),
			"read",
		})
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 100 {
		t.Errorf("Expected 100 collaborators, got %d", len(result))
	}

	// Verify first and last entries
	if result[0].RepositoryName != "repo0" {
		t.Errorf("Expected first repo to be 'repo0', got %s", result[0].RepositoryName)
	}
	if result[99].RepositoryName != "repo99" {
		t.Errorf("Expected last repo to be 'repo99', got %s", result[99].RepositoryName)
	}
}

func TestAPIGetterInterface(t *testing.T) {
	// Test that APIGetter implements the Getter interface
	var _ Getter = &APIGetter{}
}

func TestCreateRepoCollaboratorsListPermissionCaseSensitivity(t *testing.T) {
	getter := &APIGetter{}

	testData := [][]string{
		{"RepositoryName", "Username", "Permission"},
		{"repo1", "user1", "READ"},     // Uppercase
		{"repo2", "user2", "Write"},    // Mixed case
		{"repo3", "user3", "ADMIN"},    // Uppercase
		{"repo4", "user4", "maintain"}, // Lowercase
	}

	result := getter.CreateRepoCollaboratorsList(testData)

	if len(result) != 4 {
		t.Errorf("Expected 4 collaborators, got %d", len(result))
	}

	// Verify permissions are preserved as-is
	expectedPermissions := []string{"READ", "Write", "ADMIN", "maintain"}
	for i, expected := range expectedPermissions {
		if result[i].Permission != expected {
			t.Errorf("Expected permission %s, got %s", expected, result[i].Permission)
		}
	}
}
