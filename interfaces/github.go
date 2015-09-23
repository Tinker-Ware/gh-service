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

func (repo GithubRepository) GetKey(username, token string, id int) (*domain.Key, error) {
	ghkey, _, err := repo.client.Users.GetKey(id)
	if err != nil {
		return nil, err
	}

	key := &domain.Key{
		ID:    ghkey.ID,
		Title: ghkey.Title,
		Key:   ghkey.Key,
		URL:   ghkey.URL,
	}

	return key, nil
}

func (repo GithubRepository) ShowKeys(username, token string) ([]domain.Key, error) {
	ghKeys, _, err := repo.client.Users.ListKeys(username, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	keys := []domain.Key{}

	for _, k := range ghKeys {
		key := domain.Key{
			ID:    k.ID,
			Key:   k.Key,
			Title: k.Title,
			URL:   k.URL,
		}
		keys = append(keys, key)
	}

	return keys, nil

}

func (repo GithubRepository) CreateKey(username, token string, key *domain.Key) error {

	k := github.Key{
		Title: key.Title,
		Key:   key.Key,
	}

	ghK, _, err := repo.client.Users.CreateKey(&k)
	if err != nil {
		return err
	}

	key.ID = ghK.ID
	key.URL = ghK.URL

	return nil

}
