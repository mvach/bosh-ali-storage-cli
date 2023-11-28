package client_test

import (
	"errors"
	"github.com/cloudfoundry/bosh-ali-storage-cli/client"
	"github.com/cloudfoundry/bosh-ali-storage-cli/client/clientfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	Context("Put", func() {
		It("uploads a file to a blob", func() {
			storageClient := clientfakes.FakeStorageClient{}

			aliBlobstore, err := client.New(&storageClient)
			Expect(err).ToNot(HaveOccurred())

			aliBlobstore.Put("source/file/path", "destination_object")

			Expect(storageClient.UploadCallCount()).To(Equal(1))
			sourceFilePath, destination := storageClient.UploadArgsForCall(0)

			Expect(sourceFilePath).To(BeAssignableToTypeOf("source/file/path"))
			Expect(destination).To(Equal("destination_object"))
		})
	})

	Context("Get", func() {
		It("get blob downloads to a file", func() {
			storageClient := clientfakes.FakeStorageClient{}

			aliBlobstore, err := client.New(&storageClient)
			Expect(err).ToNot(HaveOccurred())

			aliBlobstore.Get("source_object", "destination/file/path")

			Expect(storageClient.DownloadCallCount()).To(Equal(1))
			sourceObject, destinationFilePath := storageClient.DownloadArgsForCall(0)

			Expect(sourceObject).To(Equal("source_object"))
			Expect(destinationFilePath).To(Equal("destination/file/path"))
		})
	})

	Context("Delete", func() {
		It("delete blob deletes the blob", func() {
			storageClient := clientfakes.FakeStorageClient{}

			aliBlobstore, err := client.New(&storageClient)
			Expect(err).ToNot(HaveOccurred())

			aliBlobstore.Delete("blob")

			Expect(storageClient.DeleteCallCount()).To(Equal(1))
			object := storageClient.DeleteArgsForCall(0)

			Expect(object).To(Equal("blob"))
		})
	})

	Context("Exists", func() {
		It("returns blob.Existing on success", func() {
			storageClient := clientfakes.FakeStorageClient{}
			storageClient.ExistsReturns(true, nil)

			aliBlobstore, _ := client.New(&storageClient)
			existsState, err := aliBlobstore.Exists("blob")
			Expect(existsState == true).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			object := storageClient.ExistsArgsForCall(0)
			Expect(object).To(Equal("blob"))
		})

		It("returns blob.NotExisting for not existing blobs", func() {
			storageClient := clientfakes.FakeStorageClient{}
			storageClient.ExistsReturns(false, nil)

			aliBlobstore, _ := client.New(&storageClient)
			existsState, err := aliBlobstore.Exists("blob")
			Expect(existsState == false).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			object := storageClient.ExistsArgsForCall(0)
			Expect(object).To(Equal("blob"))
		})

		It("returns blob.ExistenceUnknown and an error in case an error occurred", func() {
			storageClient := clientfakes.FakeStorageClient{}
			storageClient.ExistsReturns(false, errors.New("boom"))

			aliBlobstore, _ := client.New(&storageClient)
			existsState, err := aliBlobstore.Exists("blob")
			Expect(existsState == false).To(BeTrue())
			Expect(err).To(HaveOccurred())

			object := storageClient.ExistsArgsForCall(0)
			Expect(object).To(Equal("blob"))
		})
	})
})
