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

	"github.com/Tinker-Ware/gh-service/domain"
	"github.com/Tinker-Ware/gh-service/interfaces"
	. "github.com/Tinker-Ware/gh-service/usecases"
)

var _ = Describe("Usecases", func() {

	var id = 0
	var token = ""
	var username = "iasstest"
	var userToken = "123tamarindo"
	var clientID = "0d14937151de189d07a9"
	var clientSecret = "f37ca3601f3822ac37a02f51efe60843e528d4a9"
	var scopes = []string{"user:email", "delete_repo", "repo", "admin:public_key"}
	var keyID = 0
	var interactor = GHInteractor{}

	BeforeSuite(func() {
		id, token, _ = getToken(clientID, clientSecret, username, userToken)
		ghrepo, err := interfaces.NewGithubRepository(clientID, clientSecret, scopes)
		if err != nil {
			panic(err.Error())
		}

		ghrepo.SetToken(token)
		interactor = GHInteractor{
			GithubRepository: ghrepo,
		}

	})

	Describe("Test interactor functionality", func() {

		reponame := "test"

		Context("Test repo functionality", func() {

			It("Should create a repo with the test account", func() {
				repo, err := interactor.CreateRepo(username, token, reponame, "", false)

				Ω(err).ShouldNot(HaveOccurred())
				Ω(*repo.FullName).Should(ContainSubstring(reponame))
			})

			It("Should retrieve a list of repositories", func() {
				repos, err := interactor.ShowRepos(username)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(repos).Should(HaveLen(1))
			})

		})

		Context("Test public key functionality", func() {
			key := domain.Key{
				Title: github.String("TestKey"),
				Key:   github.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCz98siv2mHLiyk4MT1c6kA5BKlrLejRCpUOSHiCDcCxYN0aPbWfDRW7qMMyUrrCIcRXyd+ZPKn3O0FyDI/HKOFn3qn7PFawnG/1u6cg1H9TvPYmohQuNPt9gArmxdkecl9tFXamrSo3K3H2Uyb3RA9Q0c9NW4XDr/k1tSijSdZkhHf0tGgAuF28YGiXbri38oZsDPVkR24UajLQPfdHTFUAvmXjde7WKTU2I6zvOY/vEoaVSG5Tfnk+LsDp2L4wbl5SkMzZ6GjaQ/kn+6HBuznnSX3g0AEp9y9JiWd+YRAm46dKeRkzDm65dNP1FO/4Ovp2Xm599GB47su47DJ/2qV vagrant@web"),
			}

			It("Should Create a Key", func() {
				err := interactor.CreateKey(username, token, &key)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(*key.ID).ShouldNot(BeZero())
				keyID = *key.ID
			})

			It("Should list all keys", func() {
				keys, err := interactor.ShowKeys(username)

				Ω(err).ShouldNot(HaveOccurred())

				Ω(keys).Should(HaveLen(1))
			})

		})

		Context("Test file functionality", func() {

			dt := "# Hello"
			file := domain.File{
				Path:    "test.md",
				Content: []byte(dt),
			}

			author := domain.Author{
				Author:  "iasstest",
				Message: "Add file",
				Branch:  "master",
				Email:   "infrastructuretest@gmail.com",
			}

			It("Should create a file in the repo", func() {
				err := interactor.CreateFile(file, author, username, reponame)
				Ω(err).ShouldNot(HaveOccurred())

			})

		})

		Context("Add multiple files", func() {
			repo := "multipleFiles"
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
				Email:   "infrastructuretest@gmail.com",
			}
			It("Should create files in the repo", func() {
				_, err := interactor.CreateRepo(username, token, repo, "", false)
				Ω(err).ShouldNot(HaveOccurred())

				err = interactor.AddFiles(files, author, username, repo)
				Ω(err).ShouldNot(HaveOccurred())

			})

		})

	})

	AfterSuite(func() {
		deleteRepo("test", username, token)
		deleteRepo("multipleFiles", username, token)
		deleteKey(keyID, userToken)
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
