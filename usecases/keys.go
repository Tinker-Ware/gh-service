package usecases

import (
	"fmt"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"
)

func (interactor GHInteractor) ShowKeys(username, token string) ([]domain.Key, error) {

	client := getClient(token)

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

func (interactor GHInteractor) CreateKey(username, token string, key *domain.Key) error {

	client := getClient(token)

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

func (interactor GHInteractor) ShowKey(username, token string, id int) (*domain.Key, error) {

	client := getClient(token)

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
