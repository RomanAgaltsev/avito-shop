package end2end_test

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var serverAddr string

func init() {
	flag.StringVar(&serverAddr, "server-addr", "", "Address of the avito-shop server")
}

func TestEnd2end(t *testing.T) {
	if testing.Short() {
        t.Skip("skipping test in short mode.")
    }

	RegisterFailHandler(Fail)
	RunSpecs(t, "End2End Suite")
}

var _ = BeforeSuite(func() {
	Expect(serverAddr).NotTo(BeZero(), "Please make sure --server-addr is set correctly")
})
