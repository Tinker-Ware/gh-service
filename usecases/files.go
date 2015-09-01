package usecases

import (
	"strings"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
)

func (interactor GHInteractor) CreateFile(file domain.File, author domain.Author, username, repo, token string) error {
	client := getClient(token)

	opt := &github.RepositoryContentFileOptions{
		Message: github.String(author.Message),
		Content: file.Content,
		Branch:  github.String(author.Branch),
		Committer: &github.CommitAuthor{
			Name:  github.String(author.Author),
			Email: github.String(author.Email),
		},
	}

	_, _, err := client.Repositories.CreateFile(username, repo, file.Path, opt)
	if err != nil {
		return err
	}
	return nil
}

var README = domain.File{
	Path:    "README.md",
	Content: []byte("# This is a README"),
}

func (interactor GHInteractor) AddFiles(files []domain.File, author domain.Author, username, repo, token string) error {
	client := getClient(token)
	tree := []github.TreeEntry{}
	emptyRepo := "409 Git Repository is empty"
	lastCommit := ""

	// Get the reference sha
	// TODO remove hardcoded branch name
	ghTree, _, err := client.Git.GetRef(username, repo, "heads/master")
	if err != nil {

		// if the repo is empty create a README, else unexpected error
		if strings.Contains(err.Error(), emptyRepo) {

			// Create a copy of author for the initial commit
			author2 := author
			author2.Message = "Initial Commit"

			err = interactor.CreateFile(README, author2, username, repo, token)
			if err != nil {
				return err
			}

			// Repository should not be empty, otherwise, unexpected error
			ghTree, _, err = client.Git.GetRef(username, repo, "heads/master")
			if err != nil {
				return err
			}

		} else {
			return err

		}

	}
	lastCommit = *ghTree.Ref

	// Create a new tree
	// TODO Check for existing files and paths

	for _, file := range files {
		t := github.TreeEntry{
			Path:    github.String(file.Path),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String(string(file.Content)),
		}

		tree = append(tree, t)
	}

	// Get the current tree
	currentTree, _, err := client.Git.GetTree(username, repo, lastCommit, false)
	if err != nil {
		return err
	}

	// Add the existing files to the tree
	for _, file := range currentTree.Entries {
		t := github.TreeEntry{
			Path: file.Path,
			Mode: file.Mode,
			Type: file.Type,
			SHA:  file.SHA,
		}

		tree = append(tree, t)
	}

	newTree, _, err := client.Git.CreateTree(username, repo, *currentTree.SHA, tree)
	if err != nil {
		return err
	}

	commit := github.Commit{
		Message: github.String(author.Message),
		Parents: []github.Commit{{SHA: currentTree.SHA}},
		Tree:    &github.Tree{SHA: newTree.SHA},
	}

	newCommit, _, err := client.Git.CreateCommit(username, repo, &commit)
	if err != nil {
		return err
	}

	reference := github.Reference{
		Ref: github.String("refs/heads/master"),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	}

	_, _, err = client.Git.UpdateRef(username, repo, &reference, false)
	if err != nil {
		return err
	}

	return nil
}
