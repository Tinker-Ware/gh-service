package usecases

import (
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
