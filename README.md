# gh-collaborators

[![GitHub Release](https://img.shields.io/github/v/release/katiem0/gh-collaborators?style=flat&logo=github)](https://github.com/katiem0/gh-collaborators/releases)
[![PR Checks](https://github.com/katiem0/gh-collaborators/actions/workflows/main.yml/badge.svg)](https://github.com/katiem0/gh-collaborators/actions/workflows/main.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/katiem0/gh-collaborators)](https://goreportcard.com/report/github.com/katiem0/gh-collaborators)
[![Go Version](https://img.shields.io/github/go-mod/go-version/katiem0/gh-collaborators)](https://go.dev/)

A GitHub `gh` [CLI](https://cli.github.com/) extension to list and manage repository (outside)
collaborators in a given organization.

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation) instructions.

2. Install the extension:

   ```sh
   gh extension install katiem0/gh-collaborators
   ```

For more information: [`gh extension install`](https://cli.github.com/manual/gh_extension_install).

## Usage

The `gh-collaborators` extension supports `GitHub.com` and GitHub Enterprise Server, through
the use of `--hostname` and the following commands:

```sh
$ gh collaborators -h
List and maintain repository collaborators and their assigned repositories.

Usage:
  collaborators [command]

Available Commands:
  add         Add repo access for repository collaborators.
  list        Generate a report of repos that repository collaborators have access to.
  remove      Remove repo access for repository collaborators.

Flags:
  -h, --help   help for collaborators

Use "collaborators [command] --help" for more information about a command.
```

### List Collaborators

Repository permissions assigned to a Repository Collaborator can be listed and written to a `csv`
file for an organization or specific user.

```sh
$ gh collaborators list -h
Generate a report of repositories guest collaborators have access to.

Usage:
  collaborators list [flags] <organization>

Flags:
  -d, --debug                To debug logging
  -h, --help                 help for list
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV list to (default "RepoCollaboratorsReport-20231211162953.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")
  -u, --username string      Username of single repo collaborator to generate report for
```

The output `csv` file contains the following information:

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository where the data is extracted from. |
|`RepositoryID`| The `ID` associated with the Repository, for API usage. |
|`Visibility`| The visibility of the repository. |
|`Username`| The username of the repository collaborator. |
|`AccessLevel`| The repository access permissions granted to the repository collaborator. |

### Add Collaborators

Repository permissions can be assigned to a Repository Collaborator defined in a **required**
`csv` file for an organization.

```sh
$ gh collaborators add -h 
Add repositories and permissions for repository collaborators.

Usage:
  collaborators add [flags] <organization>

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to create access from (required)
  -h, --help               help for add
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub Personal Access Token (default "gh auth token")
```

The required  `csv` file should contain the following information:

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository that the user will be given access to. |
|`Username`| The username of the repository collaborator. |
|`AccessLevel`| The repository access permissions to grant the repository collaborator. |

### Remove Collaborators

Repository permissions can be removed for a Repository Collaborator defined in a **required**
`csv` file for an organization.

```sh
$ gh collaborators remove -h                                         
Remove repositories and permissions for repository collaborators.

Usage:
  collaborators remove [flags] <organization>

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to remove access from (required)
  -h, --help               help for remove
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub Personal Access Token (default "gh auth token")

```

The required  `csv` file should contain the following information:

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository that the user will be removed from. |
|`Username`| The username of the repository collaborator. |
