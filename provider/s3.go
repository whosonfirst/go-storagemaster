package provider

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/whosonfirst/go-storagemaster"
	_ "log"
	"path/filepath"
	"strings"
)

type S3Provider struct {
	storagemaster.Provider
	service *s3.S3
	bucket  string
	prefix  string
}

type S3Config struct {
	Bucket      string
	Prefix      string
	Region      string
	Credentials string // see notes below
}

func NewS3Provider(s3cfg S3Config) (*S3Provider, error) {

	// https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/

	cfg := aws.NewConfig()
	cfg.WithRegion(s3cfg.Region)

	if strings.HasPrefix(s3cfg.Credentials, "env:") {

		creds := credentials.NewEnvCredentials()
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(s3cfg.Credentials, "shared:") {

		details := strings.Split(s3cfg.Credentials, ":")

		if len(details) != 3 {
			return nil, errors.New("Shared credentials need to be defined as 'shared:CREDENTIALS_FILE:PROFILE_NAME'")
		}

		creds := credentials.NewSharedCredentials(details[1], details[2])
		cfg.WithCredentials(creds)

	} else if strings.HasPrefix(s3cfg.Credentials, "iam:") {

		// assume an IAM role suffient for doing whatever

	} else {

		return nil, errors.New("Unknown S3 config")
	}

	sess := session.New(cfg)

	if s3cfg.Credentials != "" {

		_, err := sess.Config.Credentials.Get()

		if err != nil {
			return nil, err
		}
	}

	service := s3.New(sess)

	c := S3Provider{
		service: service,
		bucket:  s3cfg.Bucket,
		prefix:  s3cfg.Prefix,
	}

	return &c, nil
}

func (conn *S3Provider) Exists(key string) (bool, error) {

	key = conn.prepareKey(key)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	_, err := conn.service.HeadObject(params)

	// check err here for 404-iness

	if err != nil {
		return false, err
	}

	return true, nil
}

func (conn *S3Provider) Get(key string) ([]byte, error) {

	key = conn.prepareKey(key)

	params := &s3.GetObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	rsp, err := conn.service.GetObject(params)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rsp.Body)

	return buf.Bytes(), nil
}

func (conn *S3Provider) Put(key string, body []byte, extras ...storagemaster.Extras) error {
     
	key = conn.prepareKey(key)

	params := &s3.PutObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
		// ACL:    aws.String("public-read"),
		// ContentType: aws.String(content_type),
	}

	if len(extras) == 1 {

		put_extras := extras[0]

		var v interface{}
		var e error
		
		v, e = put_extras.Get("content-type")

		if e == nil {
			params.SetContentType(v.(string))
		}

		v, e = put_extras.Get("acl")

		if e == nil {
			params.SetACL(v.(string))
		}	
	}

	_, err := conn.service.PutObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (conn *S3Provider) Delete(key string) error {

	key = conn.prepareKey(key)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(conn.bucket),
		Key:    aws.String(key),
	}

	_, err := conn.service.DeleteObject(params)

	if err != nil {
		return err
	}

	return nil
}

func (conn *S3Provider) prepareKey(key string) string {

	if conn.prefix == "" {
		return key
	}

	return filepath.Join(conn.prefix, key)
}
