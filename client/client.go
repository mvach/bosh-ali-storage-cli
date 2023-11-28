package client

type AliBlobstore struct {
	storageClient StorageClient
}

func New(storageClient StorageClient) (AliBlobstore, error) {
	return AliBlobstore{storageClient: storageClient}, nil
}

func (client *AliBlobstore) Put(sourceFilePath string, destinationObject string) error {
	return client.storageClient.Upload(sourceFilePath, destinationObject)
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
