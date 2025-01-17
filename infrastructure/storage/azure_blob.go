package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/mundipagg/boleto-api/config"
)

// AzureBlob represents a AzureBlob connection to a named container
type AzureBlob struct {
	accountName   string
	accessKey     string
	containerName string
	containerURL  azblob.ContainerURL
	devMode       bool
}

const MAX_UPLOAD_BLOCK_SIZE = 32 * 1024 * 1024 // 32MB
const CORE_MACHINE_NUMBER = 8

var jsonUploadOptions = azblob.UploadToBlockBlobOptions{
	BlockSize:   MAX_UPLOAD_BLOCK_SIZE,
	Parallelism: CORE_MACHINE_NUMBER,
}

// NewAzureBlob creates an instance of AzureBlob
// accountName is the AZURE_STORAGE_ACCOUNT
// accessKey is the AZURE_STORAGE_ACCESS_KEY
// containerName is the container name
func NewAzureBlob(accountName, accessKey, containerName string, devmode bool) (*AzureBlob, error) {
	if len(accountName) == 0 || len(accessKey) == 0 || len(containerName) == 0 {
		return nil, fmt.Errorf("either the AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY or Container name cannot be empty")
	}

	return &AzureBlob{
		accountName:   accountName,
		accessKey:     accessKey,
		containerName: containerName,
		devMode:       devmode,
	}, nil
}

func (ab *AzureBlob) connect() error {
	credential, err := azblob.NewSharedKeyCredential(ab.accountName, ab.accessKey)
	if err != nil {
		return fmt.Errorf("Invalid credentials with error: " + err.Error())
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL := ab.urlConnectionString()

	// pipeline to make requests.
	ab.containerURL = azblob.NewContainerURL(*URL, p)

	return nil
}

// Download a file from blob
func (ab *AzureBlob) Download(path, filename string) ([]byte, error) {
	ab.connect()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*20))
	defer cancel()

	blobURL := ab.containerURL.NewBlockBlobURL(fmt.Sprintf("%s/%s", path, filename))
	downloadResponse, err := blobURL.Download(
		ctx,
		0,
		azblob.CountToEnd,
		azblob.BlobAccessConditions{},
		false,
		azblob.ClientProvidedKeyOptions{},
	)
	if err != nil {
		return []byte(""), err
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return []byte(""), err
	}

	return downloadedData.Bytes(), nil
}

func (ab *AzureBlob) urlConnectionString() *url.URL {
	var URL *url.URL
	if ab.devMode {
		URL, _ = url.Parse(fmt.Sprintf("http://127.0.0.1:10000/%s/%s", ab.accountName, ab.containerName))
	} else {
		URL, _ = url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", ab.accountName, ab.containerName))
	}
	return URL
}

//UploadAsJson loads the payload into Azure Blob Storage
func (ab *AzureBlob) UploadAsJson(ctx context.Context, fileNamePrefix, payload string) (elapsedTime int64, err error) {
	if err = ab.connect(); err != nil {
		return
	}

	filename := fileNamePrefix + ".json"
	fullpath := config.Get().AzureStorageUploadPath + "/" + config.Get().AzureStorageFallbackFolder + "/" + filename

	elapsedTime, err = upload(ctx, ab, fullpath, payload)

	return
}

func upload(ctx context.Context, ab *AzureBlob, fullpath, payload string) (elapsedTime int64, err error) {
	data := []byte(payload)

	blobURL := ab.containerURL.NewBlockBlobURL(fullpath)

	start := time.Now()

	_, err = azblob.UploadBufferToBlockBlob(ctx, data, blobURL, jsonUploadOptions)

	elapsedTime = time.Since(start).Milliseconds()

	return
}
