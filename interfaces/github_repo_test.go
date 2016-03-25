package interfaces_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"net/url"

	"github.com/gh-service/domain"
	. "github.com/gh-service/interfaces"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GithubRepo", func() {
	var id = 0
	var keyID = 0
	var token = ""
	var username = "iasstest"
	var userToken = "123tamarindo"
	var clientID = "0d14937151de189d07a9"
	var clientSecret = "f37ca3601f3822ac37a02f51efe60843e528d4a9"
	var scopes = []string{"user:email", "delete_repo", "repo", "admin:public_key"}
	var repo = &GithubRepository{}

	BeforeSuite(func() {
		id, token, _ = getToken(clientID, clientSecret, username, userToken)
		repo, _ = NewGithubRepository(clientID, clientID, scopes)
		repo.SetToken(token)
	})

	reponame := "test"

	Describe("Test Oauth functionality", func() {
		Context("Get a oauth URL", func() {
			It("Should return a correctly formed url", func() {
				oauthURL, state := repo.GetOauthURL()
				Ω(oauthURL).Should(ContainSubstring(clientID))
				Ω(oauthURL).Should(ContainSubstring(state))
				Ω(oauthURL).Should(ContainSubstring(url.QueryEscape(scopes[0])))
				Ω(oauthURL).Should(ContainSubstring(url.QueryEscape(scopes[1])))
				Ω(oauthURL).Should(ContainSubstring(url.QueryEscape(scopes[2])))
				Ω(oauthURL).Should(ContainSubstring(url.QueryEscape(scopes[3])))
			})

		})
	})

	Describe("Test repo funcionality", func() {
		Context("Create a new repo", func() {

			It("Should create a new repo", func() {
				rp, err := repo.CreateRepo(username, reponame, "", false)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(*rp.FullName).Should(ContainSubstring(reponame))
			})

			It("Should retrieve a list of repositories", func() {
				repos, err := repo.GetAllRepos(username)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(repos).ShouldNot(HaveLen(0))
			})

		})
	})

	Describe("Test key functionality", func() {
		Context("With a test key", func() {
			strKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCz98siv2mHLiyk4MT1c6kA5BKlrLejRCpUOSHiCDcCxYN0aPbWfDRW7qMMyUrrCIcRXyd+ZPKn3O0FyDI/HKOFn3qn7PFawnG/1u6cg1H9TvPYmohQuNPt9gArmxdkecl9tFXamrSo3K3H2Uyb3RA9Q0c9NW4XDr/k1tSijSdZkhHf0tGgAuF28YGiXbri38oZsDPVkR24UajLQPfdHTFUAvmXjde7WKTU2I6zvOY/vEoaVSG5Tfnk+LsDp2L4wbl5SkMzZ6GjaQ/kn+6HBuznnSX3g0AEp9y9JiWd+YRAm46dKeRkzDm65dNP1FO/4Ovp2Xm599GB47su47DJ/2qV vagrant@web"

			key := domain.Key{
				Title: github.String("TestKey"),
				Key:   github.String(strKey),
			}

			It("Should Create a Key", func() {
				err := repo.CreateKey(username, &key)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(*key.ID).ShouldNot(BeZero())
				keyID = *key.ID
			})

			It("Should list all keys", func() {
				keys, err := repo.ShowKeys(username)

				Ω(err).ShouldNot(HaveOccurred())

				Ω(keys).Should(HaveLen(1))
			})

		})
	})

	Context("Add multiple files", func() {
		reponame := "multipleFiles"
		files := []domain.File{{
			Path:    "Bye.md",
			Content: []byte("# Bye"),
		},
			{
				Path:    "folder/Hello.md",
				Content: []byte("# Hello"),
			}}
		author := domain.Author{
			Author:  "iasstest",
			Message: "Hello",
			Branch:  "master",
			Email:   "infraestructuretest@gmail.com",
		}
		It("Should create files in the repo", func() {
			_, err := repo.CreateRepo(username, reponame, "", false)
			Ω(err).ShouldNot(HaveOccurred())

			err = repo.AddFiles(files, author, username, reponame)
			Ω(err).ShouldNot(HaveOccurred())

		})

	})

	AfterSuite(func() {
		deleteKey(keyID, userToken)
		deleteRepo("test", username, token)
		deleteRepo("multipleFiles", username, token)
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

func deleteKey(id int, userToken string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: userToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	client.Users.DeleteKey(id)
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
