package infraestructure_test

import (
	"io/ioutil"
	"os"

	. "github.com/gh-service/infraestructure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Infraestructure", func() {

	Describe("Read a YAML", func() {
		Context("Create a YAML on the fly", func() {

			yaml := `---
clientID: asdfg
clientSecret: aoihcou
port: 1000
scopes:
  - "user:email"
  - repo
  - "admin:public_key"`

			err := ioutil.WriteFile("conf.yaml", []byte(yaml), 0644)

			It("Should Return a Configuration struct with the correct values", func() {
				Ω(err).Should(BeNil())

				conf, err := GetConfiguration("conf.yaml")
				Ω(err).Should(BeNil())

				Ω(conf.Port).Should(Equal("1000"))
				Ω(conf.ClientID).Should(Equal("asdfg"))
				Ω(conf.ClientSecret).Should(Equal("aoihcou"))

				Ω(conf.Scopes).Should(ContainElement("user:email"))
				Ω(conf.Scopes).Should(ContainElement("repo"))
				Ω(conf.Scopes).Should(ContainElement("admin:public_key"))

				os.Remove("conf.yaml")

			})
		})
	})

})
