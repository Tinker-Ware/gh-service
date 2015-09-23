package interfaces

import (
	"fmt"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubRepository struct {
	client *github.Client
}

func (repo *GithubRepository) SetToken(token string) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	repo.client = client

}

func (repo GithubRepository) GetAllRepos(username string) ([]domain.Repository, error) {

	// list all repositories for the authenticated user
	ghrepos, _, err := repo.client.Repositories.List("", nil)

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

func (repo GithubRepository) GetRepo(username, token, reponame string) (*domain.Repository, error) {

	rp, _, err := repo.client.Repositories.Get(username, reponame)
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

func (repo GithubRepository) CreateRepo(username, reponame, org string, private bool) (*domain.Repository, error) {
	rp := &github.Repository{
		Name:    github.String(reponame),
		Private: github.Bool(private),
	}
	rp, _, err := repo.client.Repositories.Create(org, rp)

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
