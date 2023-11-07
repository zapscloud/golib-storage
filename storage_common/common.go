package storage_common

import "github.com/zapscloud/golib-utils/utils"

// Enums
type StorageType byte

const (
	STORAGE_TYPE_NONE StorageType = iota
	STORAGE_TYPE_AWS_S3
	STORAGE_TYPE_MS_AZURE
	STORAGE_TYPE_GOOGLE
	STORAGE_TYPE_PLACEHOLDER_LAST // Only a place holder
)

const (
	STORAGE_TYPE = "storage_type"

	// Param for AWS_S3
	STORAGE_AWS_S3_REGION    = "aws_s3_region"
	STORAGE_AWS_S3_BUCKET    = "aws_s3_bucket"
	STORAGE_AWS_S3_ACCESSKEY = "aws_s3_accesskey"
	STORAGE_AWS_S3_SECRETKEY = "aws_s3_secretkey"
)

const (
	PUT_OBJECT = "PUT"
	GET_OBJECT = "GET"
)

// Default values
const (
	DEFAULT_CONTENT_TYPE = "image/png"
)

func GetStorageType(props utils.Map) (StorageType, error) {

	dataVal, dataOk := props[STORAGE_TYPE]

	// Convert it to String type
	storageType := dataVal.(StorageType)

	if !dataOk && (storageType <= STORAGE_TYPE_NONE || storageType >= STORAGE_TYPE_PLACEHOLDER_LAST) {

		err := &utils.AppError{ErrorStatus: 401, ErrorCode: "401", ErrorMsg: "Invalid StorageType", ErrorDetail: "Either StorageType value is not sent or Invalid"}
		return STORAGE_TYPE_NONE, err
	}

	return storageType, nil
}
