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
	ghoauth "golang.org/x/oauth2/github"
)

type UserRepo interface {
	Store(user domain.User) error
	RetrieveByID(id int64) (*domain.User, error)
	RetrieveByUserName(username string) (*domain.User, error)
	Update(incUser domain.User) error
}

type GHInteractor struct {
	UserRepo UserRepo
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	oauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user:email", "repo", "admin:public_key"},
		Endpoint:     ghoauth.Endpoint,
	}
)

func (interactor GHInteractor) GHLogin() (string, string) {
	oauthStateString := randSeq(10)
	url := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	return url, oauthStateString

}

func (interactor GHInteractor) GHCallback(code, state, incomingState string) (*domain.User, error) {

	// TODO: Implement this later
	// if state != incomingState {
	// 	fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", incomingState, state)
	// 	return nil, errors.New("Invalid Oauth2 state")
	// }

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		// TODO: Log with interactor Logger not yet implemented
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err.Error())
		return nil, err
	}

	oauthClient := oauthConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, errors.New("Cannot retrieve User data")
	}

	usr := domain.User{
		Username:    *user.Login,
		AccessToken: token.AccessToken,
	}

	interactor.UserRepo.Store(usr)

	return &usr, nil
}

func (interactor GHInteractor) ShowUser(username string) (*domain.User, error) {

	user, err := interactor.UserRepo.RetrieveByUserName(username)
	if err != nil {
		return nil, err
	}

	return user, nil
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
