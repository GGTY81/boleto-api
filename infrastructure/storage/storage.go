package storage

import (
	"context"

	"github.com/mundipagg/boleto-api/config"
)

type IStorage interface {
	UploadAsJson(ctx context.Context, fullpath, payload string) (totalElapsedTimeInMilliseconds int64, err error)
}

// GetClient factory storage
func GetClient() (IStorage, error) {
	return NewAzureBlob(
		config.Get().AzureStorageAccount,
		config.Get().AzureStorageAccessKey,
		config.Get().AzureStorageContainerName,
		config.Get().DevMode,
	)
}
