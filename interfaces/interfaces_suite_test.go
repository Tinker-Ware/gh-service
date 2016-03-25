package interfaces_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInterfaces(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interfaces Suite")
}
