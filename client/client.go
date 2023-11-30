package client

import (
	"fmt"
	"log"
	"strings"
)

type AliBlobstore struct {
	storageClient StorageClient
}

func New(storageClient StorageClient) (AliBlobstore, error) {
	return AliBlobstore{storageClient: storageClient}, nil
}

func (client *AliBlobstore) Put(sourceFilePath string, destinationObject string) error {
	err := client.storageClient.Upload(sourceFilePath, destinationObject)
	if err != nil {
		return fmt.Errorf("upload failure: %w", err)
	}

	log.Println("Successfully uploaded file")
	return nil
}

func (client *AliBlobstore) Get(sourceObject string, destinationFilePath string) error {
	return client.storageClient.Download(sourceObject, destinationFilePath)
}

func (client *AliBlobstore) Delete(object string) error {
	return client.storageClient.Delete(object)
}

func (client *AliBlobstore) Exists(object string) (bool, error) {
	return client.storageClient.Exists(object)
}

func (client *AliBlobstore) Sign(object string, action string, expiredInSec int64) (string, error) {
	action = strings.ToUpper(action)
	switch action {
	case "PUT":
		return client.storageClient.SignedUrlPut(object, expiredInSec)
	case "GET":
		return client.storageClient.SignedUrlGet(object, expiredInSec)
	default:
		return "", fmt.Errorf("action not implemented: %s", action)
	}
}
