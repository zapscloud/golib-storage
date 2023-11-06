package storage_services

import (
	"github.com/zapscloud/golib-storage/storage_common"
	"github.com/zapscloud/golib-storage/storage_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// StorageService - Storage Service
type StorageService interface {
	InitializeStorageService(props utils.Map) error
	GetSignedURL(method string, fileKey string) (string, error)
	UploadFile(fileName string, fileKey string) (string, string, error)
}

// NewStorageService - Contruct Storage Service
func NewStorageService(props utils.Map) (StorageService, error) {
	var storageClient StorageService = nil

	// Get StorageType from the Parameter
	storageType, err := storage_common.GetStorageType(props)
	if err != nil {
		return nil, err
	}

	// Get the Storage's Object based on StorageType
	switch storageType {
	case storage_common.STORAGE_TYPE_AWS_S3:
		storageClient = &storage_repository.AWSStorageServices{}
	case storage_common.STORAGE_TYPE_MS_AZURE:
		// *Not Implemented yet*
		storageClient = nil
	case storage_common.STORAGE_TYPE_GOOGLE:
		// *Not Implemented yet*
		storageClient = nil
	}

	if storageClient != nil {
		// Initialize the Dao
		err = storageClient.InitializeStorageService(props)
		if err != nil {
			return nil, err
		}
	}

	return storageClient, nil
}
