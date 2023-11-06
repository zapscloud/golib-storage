package storage_repository

// AWS S3 Storage Implementations

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/zapscloud/golib-storage/storage_common"
	"github.com/zapscloud/golib-utils/utils"
)

// AWSStorageServices - AWS Storage Service structure
type AWSStorageServices struct {
	awsS3Region  string
	awsS3Bucket  string
	awsAccessKey string
	awsSecretKey string
}

func (p *AWSStorageServices) InitializeStorageService(props utils.Map) error {

	if _, dataOk := props[storage_common.STORAGE_AWS_S3_REGION]; !dataOk {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "Parameter S3 Region is not received"}
		return err
	} else if _, dataOk := props[storage_common.STORAGE_AWS_S3_BUCKET]; !dataOk {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "Parameter S3 Bucket is not received"}
		return err
	} else if _, dataOk := props[storage_common.STORAGE_AWS_S3_ACCESSKEY]; !dataOk {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "Parameter S3 AccessKey is not received"}
		return err
	} else if _, dataOk := props[storage_common.STORAGE_AWS_S3_SECRETKEY]; !dataOk {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "Parameter S3 SecretKey is not received"}
		return err
	}

	// Store the Parameter to member variable
	p.awsS3Region = props[storage_common.STORAGE_AWS_S3_REGION].(string)
	p.awsS3Bucket = props[storage_common.STORAGE_AWS_S3_BUCKET].(string)
	p.awsAccessKey = props[storage_common.STORAGE_AWS_S3_ACCESSKEY].(string)
	p.awsSecretKey = props[storage_common.STORAGE_AWS_S3_SECRETKEY].(string)

	return nil
}

// Get SignedURL from S3
func (p *AWSStorageServices) GetSignedURL(method string, fileKey string) (string, error) {

	// Validate Input Method
	method = strings.ToUpper(method)
	if method != storage_common.GET_OBJECT && method != storage_common.PUT_OBJECT {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "Parameter method should either GET or PUT"}
		return "", err
	}

	// Create New Session
	sess, err := p.createNewSession()
	if err != nil {
		log.Println("Error while creating NewSession:: ", err)
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Error while creating NewSession", ErrorDetail: err.Error()}
		return "", err
	}

	var req *request.Request
	// Create S3 service client
	svc := s3.New(sess)

	if method == storage_common.PUT_OBJECT {
		req, _ = svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(p.awsS3Bucket),
			Key:    aws.String(fileKey),
		})
	} else if method == storage_common.GET_OBJECT {
		req, _ = svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(p.awsS3Bucket),
			Key:    aws.String(fileKey),
		})

	}

	// Presign the URL
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Error while getting URL:: ", err)
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Error while presign", ErrorDetail: err.Error()}
		return "", err
	}

	// Everything success, return the url
	return urlStr, nil
}

func (p *AWSStorageServices) UploadFile(fileName string, fileKey string) (string, string, error) {

	// The session the S3 Uploader will use
	sess := session.Must(p.createNewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		log.Println("golib-storage::UploadFile::failed to open file ", fileName, err)
		return "", "", err
	}
	defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(p.awsS3Bucket),
		Key:    aws.String(fileKey),
		Body:   f,
	})

	if err != nil {
		log.Println("golib-storage::UploadFile::failed to upload file", err)
		return "", "", err
	}

	log.Println("golib-storage::UploadFile::file uploaded to => ", result.Location)

	return fileKey, result.Location, nil
}

func (p *AWSStorageServices) DownloadFile(fileName string, fileKey string) error {

	// The session the S3 Uploader will use
	sess := session.Must(p.createNewSession())

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	f, err := os.Create(fileName)
	if err != nil {
		log.Println("golib-storage::DownloadFile::failed to open file ", fileName, err)
		return err
	}
	defer f.Close()

	// Upload the file to S3.
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(p.awsS3Bucket),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		log.Println("golib-storage::DownloadFile::failed to upload file", err)
		return err
	}

	log.Println("golib-storage::DownloadFile::file downloaded to => ", fileName, n)

	return nil
}

func (p *AWSStorageServices) createNewSession() (*session.Session, error) {
	// Assign Credentials
	s3Creds := credentials.NewStaticCredentials(p.awsAccessKey, p.awsSecretKey, "")
	// Create Configuration
	s3Cfg := aws.NewConfig().WithRegion(p.awsS3Region).WithCredentials(s3Creds)

	// Create New Session
	return session.NewSession(s3Cfg)
}
