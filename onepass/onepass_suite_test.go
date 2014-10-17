package onepass_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTermpass(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Termpass Suite")
}
