package usecases

import (
	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
)

func (interactor GHInteractor) CreateFile(file domain.File, username, repo, token string) error {
	client := getClient(token)

	opt := &github.RepositoryContentFileOptions{
		Message: github.String(file.Message),
		Content: file.Content,
		Branch:  github.String(file.Branch),
		Committer: &github.CommitAuthor{
			Name:  github.String(file.Author),
			Email: github.String(file.Email),
		},
	}

	_, _, err := client.Repositories.CreateFile(username, repo, file.Path, opt)
	if err != nil {
		return err
	}
	return nil
}
