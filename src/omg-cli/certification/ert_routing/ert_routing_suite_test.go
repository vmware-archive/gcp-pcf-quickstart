package ert_routing_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestErtRouting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ErtRouting Suite")
}
