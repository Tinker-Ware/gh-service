package usecases

import (
	"fmt"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (interactor GHInteractor) ShowKeys(username string) ([]domain.Key, error) {
	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	ghKeys, _, err := client.Users.ListKeys(username, nil)
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

func (interactor GHInteractor) CreateKey(username string, key *domain.Key) error {

	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	k := github.Key{
		Title: key.Title,
		Key:   key.Key,
	}

	ghK, _, err := client.Users.CreateKey(&k)
	if err != nil {
		return err
	}

	key.ID = ghK.ID
	key.URL = ghK.URL

	return nil
}

func (interactor GHInteractor) ShowKey(username string, id int) (*domain.Key, error) {
	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	ghkey, _, err := client.Users.GetKey(id)
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
