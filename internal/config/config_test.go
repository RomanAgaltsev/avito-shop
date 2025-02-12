package config_test

import (
	"flag"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/RomanAgaltsev/avito-shop/internal/config"
)

type testCase struct {
	envNam string
	envVal string
	flgNam string
	flgVal string
	defVal string
}

var _ = Describe("Config", func() {
	var cfg *config.Config
	var err error

	defaultArgs := os.Args
	defaultCommandLine := flag.CommandLine

	AfterEach(func() {
		os.Clearenv()
		os.Args = defaultArgs
		flag.CommandLine = defaultCommandLine
	})

	// Run address
	ra := &testCase{
		envNam: "RUN_ADDRESS",
		envVal: "localhost:8081",
		flgNam: "-a",
		flgVal: "localhost:8091",
		defVal: "localhost:8080",
	}

	DescribeTable("Run address",
		func(envName, envVal, flgName, flgVal, def, expected string) {
			setEnv(envName, envVal)
			setFlag(flgName, flgVal)

			cfg, err = config.Get()

			Expect(err).Should(BeNil())
			Expect(cfg.RunAddress).To(Equal(expected))
		},

		EntryDescription("When env %s=%s, flag %s=%s and default=%s"),
		Entry(nil, ra.envNam, ra.envVal, ra.flgNam, ra.flgVal, ra.defVal, ra.envVal),
		Entry(nil, ra.envNam, ra.envVal, ra.flgNam, "", ra.defVal, ra.envVal),
		Entry(nil, ra.envNam, ra.envVal, "", "", ra.defVal, ra.envVal),
		Entry(nil, ra.envNam, "", ra.flgNam, ra.flgVal, ra.defVal, ra.flgVal),
		Entry(nil, "", "", ra.flgNam, ra.flgVal, ra.defVal, ra.flgVal),
		Entry(nil, ra.envNam, "", ra.flgNam, "", ra.defVal, ""),
		Entry(nil, "", "", ra.flgNam, "", ra.defVal, ""),
		Entry(nil, ra.envNam, "", ra.flgNam, "", ra.defVal, ""),
		Entry(nil, ra.envNam, "", "", "", ra.defVal, ra.defVal),
		Entry(nil, "", "", "", "", ra.defVal, ra.defVal),
		Entry(nil, "", "", "", "", ra.defVal, ra.defVal),
	)

	// Database URI
	du := &testCase{
		envNam: "DATABASE_URI",
		envVal: "postgres://postgres:12345@localhost:5432/avitoshop?sslmode=disable",
		flgNam: "-d",
		flgVal: "postgres://postgres:12346@localhost:5433/avitoshop?sslmode=disable",
		defVal: "",
	}

	DescribeTable("Database URI",
		func(envName, envVal, flgName, flgVal, def, expected string) {
			setEnv(envName, envVal)
			setFlag(flgName, flgVal)

			cfg, err = config.Get()

			Expect(err).Should(BeNil())
			Expect(cfg.DatabaseURI).To(Equal(expected))
		},

		EntryDescription("When env %s=%s, flag %s=%s and default=%s"),
		Entry(nil, du.envNam, du.envVal, du.flgNam, du.flgVal, du.defVal, du.envVal),
		Entry(nil, du.envNam, du.envVal, du.flgNam, "", du.defVal, du.envVal),
		Entry(nil, du.envNam, du.envVal, "", "", du.defVal, du.envVal),
		Entry(nil, du.envNam, "", du.flgNam, du.flgVal, du.defVal, du.flgVal),
		Entry(nil, "", "", du.flgNam, du.flgVal, du.defVal, du.flgVal),
		Entry(nil, du.envNam, "", du.flgNam, "", du.defVal, ""),
		Entry(nil, "", "", du.flgNam, "", du.defVal, ""),
		Entry(nil, du.envNam, "", du.flgNam, "", du.defVal, ""),
		Entry(nil, du.envNam, "", "", "", du.defVal, du.defVal),
		Entry(nil, "", "", "", "", du.defVal, du.defVal),
		Entry(nil, "", "", "", "", du.defVal, du.defVal),
	)

	// Secret key
	sk := &testCase{
		envNam: "SECRET_KEY",
		envVal: "very secret key",
		defVal: "secret",
	}

	DescribeTable("Secret key",
		func(envName, envVal, def, expected string) {
			setEnv(envName, envVal)

			cfg, err = config.Get()

			Expect(err).Should(BeNil())
			Expect(cfg.SecretKey).To(Equal(expected))
		},

		EntryDescription("When env %s=%s and default=%s"),
		Entry(nil, sk.envNam, sk.envVal, sk.defVal, sk.envVal),
		Entry(nil, sk.envNam, "", sk.defVal, sk.defVal),
		Entry(nil, "", "", sk.defVal, sk.defVal),
	)
})

func setEnv(name, value string) {
	if name != "" {
		t := GinkgoT()
		t.Setenv(name, value)
	}
}

func setFlag(name, value string) {
	if name != "" {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = append([]string{"cmd"}, name, value)
	}
}
