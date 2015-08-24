package usecases_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gh-service/usecases"
	ghoauth "golang.org/x/oauth2/github"
)

var _ = Describe("Usecases", func() {

	var id = 0
	var token = ""
	var username = "iasstest"
	var userToken = "123tamarindo"
	var clientID = "0d14937151de189d07a9"
	var clientSecret = "f37ca3601f3822ac37a02f51efe60843e528d4a9"

	BeforeSuite(func() {
		id, token, _ = getToken(clientID, clientSecret, username, userToken)
	})

	Describe("Test Repo funcionality", func() {

		Context("With an interactor and a test app", func() {

			oauth2client := &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Scopes:       []string{"user:email", "delete_repo", "repo", "admin:public_key"},
				Endpoint:     ghoauth.Endpoint,
			}

			interactor := GHInteractor{
				OauthConfig: oauth2client,
			}

			reponame := "test"
			It("Should create a repo with the test account", func() {
				repo, err := interactor.CreateRepo(username, token, reponame, "", false)

				Ω(err).ShouldNot(HaveOccurred())
				Ω(*repo.FullName).Should(ContainSubstring(reponame))
			})

			It("Should retrieve a list of repositories", func() {
				repos, err := interactor.ShowRepos(username, token)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(repos).Should(HaveLen(1))
			})

		})

	})

	AfterSuite(func() {
		deleteRepo("test", username, token)
		deleteToken(id, username, userToken)
	})

})

type request struct {
	Scopes       []string `json:"scopes"`
	Note         string   `json:"note"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Fingerprint  string   `json:"fingerprint"`
	Description  string   `json:"description"`
}

type response struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func getToken(clientID, clientSecret, username, userToken string) (int, string, error) {
	rq := request{
		Scopes:       []string{"user:email", "delete_repo", "repo", "admin:public_key"},
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Description:  "IAAS Testing",
	}

	b, _ := json.Marshal(rq)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewBuffer(b))
	req.SetBasicAuth(username, userToken)
	resp, err := client.Do(req)
	if err != nil {
		Fail(err.Error())
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	respT := response{}

	err = decoder.Decode(&respT)
	if err != nil {
		return 0, "", err
	}

	return respT.ID, respT.Token, nil

}

func deleteToken(id int, username, userToken string) {

	client := &http.Client{}

	path := fmt.Sprintf("https://api.github.com/authorizations/%d", id)

	req, _ := http.NewRequest("DELETE", path, nil)
	req.SetBasicAuth(username, userToken)

	client.Do(req)
}

func deleteRepo(reponame, username, userToken string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: userToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	client.Repositories.Delete(username, reponame)
}
