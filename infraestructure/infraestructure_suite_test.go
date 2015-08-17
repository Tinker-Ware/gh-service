package infraestructure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInfraestructure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Infraestructure Suite")
}
