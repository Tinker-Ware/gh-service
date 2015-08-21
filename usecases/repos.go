package usecases

import (
	"fmt"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
)

func (interactor GHInteractor) ShowRepos(username, token string) ([]domain.Repository, error) {

	client := getClient(token)

	// list all repositories for the authenticated user
	ghrepos, _, err := client.Repositories.List("", nil)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	repos := []domain.Repository{}

	for _, repo := range ghrepos {

		r := domain.Repository{
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Private:     repo.Private,
			HTMLURL:     repo.HTMLURL,
			CloneURL:    repo.CloneURL,
			SSHURL:      repo.SSHURL,
		}
		repos = append(repos, r)
	}

	return repos, nil
}

func (interactor GHInteractor) CreateRepo(username, token, reponame, org string, private bool) (*domain.Repository, error) {
	client := getClient(token)

	repo := &github.Repository{
		Name:    github.String(reponame),
		Private: github.Bool(private),
	}
	repo, _, err := client.Repositories.Create(org, repo)

	if err != nil {
		return nil, err
	}

	r := &domain.Repository{
		Name:        repo.Name,
		FullName:    repo.FullName,
		Description: repo.Description,
		Private:     repo.Private,
		HTMLURL:     repo.HTMLURL,
		CloneURL:    repo.CloneURL,
		SSHURL:      repo.SSHURL,
	}

	return r, nil
}

func (interactor GHInteractor) ShowRepo(username, token, repo string) (*domain.Repository, error) {

	client := getClient(token)

	rp, _, err := client.Repositories.Get(username, repo)
	if err != nil {
		return nil, err
	}

	r := &domain.Repository{
		Name:        rp.Name,
		FullName:    rp.FullName,
		Description: rp.Description,
		Private:     rp.Private,
		HTMLURL:     rp.HTMLURL,
		CloneURL:    rp.CloneURL,
		SSHURL:      rp.SSHURL,
	}

	return r, nil
}
