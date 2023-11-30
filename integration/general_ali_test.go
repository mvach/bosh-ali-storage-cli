package integration_test

import (
	"bytes"
	"os"

	"github.com/cloudfoundry/bosh-ali-storage-cli/config"
	"github.com/cloudfoundry/bosh-ali-storage-cli/integration"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("General testing for all Ali regions", func() {

	Describe("Invoking `put`", func() {
		var blobName string
		var configPath string
		var contentFile string

		BeforeEach(func() {
			blobName = integration.GenerateRandomString()
			configPath = integration.MakeConfigFile(&defaultConfig)
			contentFile = integration.MakeContentFile("foo")
		})

		AfterEach(func() {
			defer func() { _ = os.Remove(configPath) }()
			defer func() { _ = os.Remove(contentFile) }()
		})

		It("uploads a file", func() {
			defer func() {
				cliSession, err := integration.RunCli(cliPath, configPath, "delete", blobName)
				Expect(err).ToNot(HaveOccurred())
				Expect(cliSession.ExitCode()).To(BeZero())
			}()

			cliSession, err := integration.RunCli(cliPath, configPath, "put", contentFile, blobName)
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(BeZero())

			cliSession, err = integration.RunCli(cliPath, configPath, "exists", blobName)
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(BeZero())

			Expect(string(cliSession.Err.Contents())).To(MatchRegexp("File '" + blobName + "' exists in bucket '" + bucketName + "'"))
		})

		It("overwrites an existing file", func() {
			defer func() {
				cliSession, err := integration.RunCli(cliPath, configPath, "delete", blobName)
				Expect(err).ToNot(HaveOccurred())
				Expect(cliSession.ExitCode()).To(BeZero())
			}()

			tmpLocalFile, _ := os.CreateTemp("", "ali-storage-cli-download")
			tmpLocalFile.Close()
			defer func() { _ = os.Remove(tmpLocalFile.Name()) }()

			contentFile = integration.MakeContentFile("initial content")
			cliSession, err := integration.RunCli(cliPath, configPath, "put", contentFile, blobName)
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(BeZero())

			cliSession, err = integration.RunCli(cliPath, configPath, "get", blobName, tmpLocalFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(BeZero())

			gottenBytes, _ := os.ReadFile(tmpLocalFile.Name())
			Expect(string(gottenBytes)).To(Equal("initial content"))

			contentFile = integration.MakeContentFile("updated content")
			cliSession, err = integration.RunCli(cliPath, configPath, "put", contentFile, blobName)
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(BeZero())

			cliSession, err = integration.RunCli(cliPath, configPath, "get", blobName, tmpLocalFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(BeZero())

			gottenBytes, _ = os.ReadFile(tmpLocalFile.Name())
			Expect(string(gottenBytes)).To(Equal("updated content"))
		})

		It("returns the appropriate error message", func() {
			cfg := &config.AliStorageConfig{
				AccessKeyID:     accessKeyID,
				AccessKeySecret: accessKeySecret,
				Endpoint:        endpoint,
				BucketName:      "not-existing",
			}

			configPath = integration.MakeConfigFile(cfg)

			cliSession, err := integration.RunCli(cliPath, configPath, "put", contentFile, blobName)
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(Equal(1))

			consoleOutput := bytes.NewBuffer(cliSession.Err.Contents()).String()
			Expect(consoleOutput).To(ContainSubstring("upload failure"))
		})
	})
	Describe("Invoking `-v`", func() {
		It("returns the cli version", func() {
			configPath := integration.MakeConfigFile(&defaultConfig)
			defer func() { _ = os.Remove(configPath) }()

			cliSession, err := integration.RunCli(cliPath, configPath, "-v")
			Expect(err).ToNot(HaveOccurred())
			Expect(cliSession.ExitCode()).To(Equal(0))

			consoleOutput := bytes.NewBuffer(cliSession.Out.Contents()).String()
			Expect(consoleOutput).To(ContainSubstring("version"))
		})
	})
})
