package usecases

import (
	"fmt"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (interactor GHInteractor) ShowRepos(username string) ([]domain.Repository, error) {
	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

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

func (interactor GHInteractor) CreateRepo(username, reponame, org string, private bool) (*domain.Repository, error) {
	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	repo := &github.Repository{
		Name:    github.String(reponame),
		Private: github.Bool(private),
	}
	repo, _, err = client.Repositories.Create(org, repo)

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

func (interactor GHInteractor) ShowRepo(username, repo string) (*domain.Repository, error) {

	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

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
