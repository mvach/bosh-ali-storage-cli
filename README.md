# Ali Storage CLI

The Ali Storage CLI is for uploading, fetching and deleting content to and from an Ali OSS.
It is highly inspired by the https://github.com/cloudfoundry/bosh-s3cli.

## Usage

Given a JSON config file (`config.json`)...

``` json
{
  "access_key_id":             "<string> (required)",
  "access_key_secret":         "<string> (required)",
  "endpoint":                  "<string> (required)",
  "bucket_name":               "<string> (required)"
}
```

``` bash
# Command: "put"
# Upload a blob to the blobstore.
./bosh-ali-storage-cli -c config.json put <path/to/file> <remote-blob>

# Command: "get"
# Fetch a blob from the blobstore.
# Destination file will be overwritten if exists.
./bosh-ali-storage-cli -c config.json get <remote-blob> <path/to/file>

# Command: "delete"
# Remove a blob from the blobstore.
./bosh-ali-storage-cli -c config.json delete <remote-blob>

# Command: "exists"
# Checks if blob exists in the blobstore.
./bosh-ali-storage-cli -c config.json exists <remote-blob>

# Command: "sign"
# Create a self-signed url for a blob in the blobstore.
./bosh-ali-storage-cli -c config.json sign <remote-blob> <get|put> <seconds-to-expiration>
```

### Using signed urls with curl
``` bash
# Uploading a blob:
curl -X PUT -T path/to/file <signed url>

# Downloading a blob:
curl -X GET <signed url>
```


