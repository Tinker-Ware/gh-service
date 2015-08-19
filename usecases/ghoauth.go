package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gh-service/domain"
	"github.com/google/go-github/github"

	"golang.org/x/oauth2"
)

type GHInteractor struct {
	OauthConfig *oauth2.Config
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func (interactor GHInteractor) GHLogin() (string, string) {
	oauthStateString := randSeq(10)
	url := interactor.OauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	return url, oauthStateString

}

func (interactor GHInteractor) GHCallback(code, state, incomingState string) (*domain.User, error) {

	// TODO: Implement this later
	// if state != incomingState {
	// 	fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", incomingState, state)
	// 	return nil, errors.New("Invalid Oauth2 state")
	// }

	token, err := interactor.OauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		// TODO: Log with interactor Logger not yet implemented
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err.Error())
		return nil, err
	}

	oauthClient := interactor.OauthConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, errors.New("Cannot retrieve User data")
	}

	usr := domain.User{
		Username:    *user.Login,
		AccessToken: token.AccessToken,
	}

	return &usr, nil
}

func (interactor GHInteractor) ShowUser(username, token string) (*domain.User, error) {

	client := getClient(token)
	user, _, err := client.Users.Get(username)

	if err != nil {
		return nil, err
	}

	usr := domain.User{
		ID:       *user.ID,
		Username: *user.Login,
	}

	return &usr, nil
}

func getClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	return client

}

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
