package ert_routing_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestErtRouting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ErtRouting Suite")
}
