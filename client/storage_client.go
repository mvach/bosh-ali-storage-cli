package client

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cloudfoundry/bosh-ali-storage-cli/config"
	"log"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . StorageClient
type StorageClient interface {
	Upload(
		sourceFilePath string,
		destinationObject string,
	) error

	Download(
		sourceObject string,
		destinationFilePath string,
	) error
}

type DefaultStorageClient struct {
	storageConfig config.AliStorageConfig
}

func NewStorageClient(storageConfig config.AliStorageConfig) (StorageClient, error) {
	return DefaultStorageClient{storageConfig: storageConfig}, nil
}

func (dsc DefaultStorageClient) Upload(
	sourceFilePath string,
	destinationObject string,
) error {
	log.Println(fmt.Sprintf("Uploading %s/%s", dsc.storageConfig.BucketName, destinationObject))

	client, err := oss.New(dsc.storageConfig.Endpoint, dsc.storageConfig.AccessKeyID, dsc.storageConfig.AccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(dsc.storageConfig.BucketName)
	if err != nil {
		return err
	}

	return bucket.PutObjectFromFile(destinationObject, sourceFilePath)
}

func (dsc DefaultStorageClient) Download(
	sourceObject string,
	destinationFilePath string,
) error {
	log.Println(fmt.Sprintf("Downloading %s/%s", dsc.storageConfig.BucketName, sourceObject))

	client, err := oss.New(dsc.storageConfig.Endpoint, dsc.storageConfig.AccessKeyID, dsc.storageConfig.AccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(dsc.storageConfig.BucketName)
	if err != nil {
		return err
	}

	return bucket.GetObjectToFile(sourceObject, destinationFilePath)
}
