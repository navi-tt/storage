package qingyun

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/navi-tt/storage"
	"github.com/qingstor/qingstor-sdk-go/v4/config"
	qsErrors "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	qs "github.com/qingstor/qingstor-sdk-go/v4/service"
	"io"
	"net"
	"net/http"
	"time"
)

// 青云存储服务
type QingStor struct {
	AccesskeyId     string
	SecretAccessKey string
	Zone            string
	Bucket          string
	Protocol        string
	Host            string
	Port            int
}

func (s *qingStor) Init(cfg string) (storage.Storage, error) {
	fmt.Printf("[QS Init] config: %s \n", cfg)
	qsConfig := &QingStor{}
	if err := json.Unmarshal([]byte(cfg), qsConfig); err != nil {
		return nil, err
	}

	qsCfg, err := config.New(qsConfig.AccesskeyId, qsConfig.SecretAccessKey)
	if err != nil {
		return nil, err
	}

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   5,
	}
	qsCfg.Connection = &http.Client{
		Transport: t,
		Timeout:   time.Second * 3,
	}

	if qsConfig.Host != "" {
		qsCfg.Host = qsConfig.Host
	}

	if qsConfig.Port > 0 {
		qsCfg.Port = qsConfig.Port
	}

	if qsConfig.Protocol != "" {
		qsCfg.Protocol = qsConfig.Protocol
	}
	qsSvc, err := qs.Init(qsCfg)
	if err != nil {
		return nil, err
	}

	bucket, err := qsSvc.Bucket(qsConfig.Bucket, qsConfig.Zone)
	if err != nil {
		return nil, err
	}

	qos = &qingStor{
		qsSvc:  qsSvc,
		bucket: bucket,
	}

	return qos, nil
}

var qos = &qingStor{}

func init() {
	storage.Register("qs", qos)
}

type qingStor struct {
	qsSvc  *qs.Service
	bucket *qs.Bucket
}

func (s *qingStor) Put(key string, r io.Reader, contentLength int64) error {
	fmt.Printf("[QS PUT] object: %s\n", key)
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	output, err := s.bucket.PutObject(key, &qs.PutObjectInput{
		ContentLength: qs.Int64(contentLength),
		ContentType:   qs.String("text/plain"),
		Body:          r,
	})

	if err != nil {
		return err
	}

	code := qs.IntValue(output.StatusCode)
	if code != http.StatusCreated {
		return fmt.Errorf("failed to put object with code %d", code)
	}

	return nil
}

func (s *qingStor) Get(key string, wa io.WriterAt) error {
	fmt.Printf("[QS GET] object: %s\n", key)

	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	output, err := s.bucket.GetObject(key, &qs.GetObjectInput{})

	if err != nil {
		if qsErr, ok := err.(*qsErrors.QingStorError); ok && qsErr.StatusCode == http.StatusNotFound {
			return storage.ErrObjectNotFound
		}
		return err
	}
	defer output.Close()

	return storage.Copy(wa, output.Body)
}

// 调用者需要关闭文件
func (s *qingStor) FileStream(key string) (io.ReadCloser, *storage.FileInfo, error) {
	if !storage.ValidKey(key) {
		return nil, nil, storage.ErrObjectKeyInvalid
	}

	output, err := s.bucket.GetObject(key, &qs.GetObjectInput{})
	if err != nil {
		if qsErr, ok := err.(*qsErrors.QingStorError); ok && qsErr.StatusCode == http.StatusNotFound {
			return nil, nil, storage.ErrObjectNotFound
		}
		return nil, nil, err
	}
	defer output.Close()

	return output.Body, &storage.FileInfo{
		Size:    *output.ContentLength,
		ModTime: qs.TimeValue(output.LastModified),
		Mode:    0666,
	}, nil
}

func (s *qingStor) Stat(key string) (*storage.FileInfo, error) {
	if !storage.ValidKey(key) {
		return nil, storage.ErrObjectKeyInvalid
	}

	output, err := s.bucket.GetObject(key, &qs.GetObjectInput{})
	if err != nil {
		if qsErr, ok := err.(*qsErrors.QingStorError); ok && qsErr.StatusCode == http.StatusNotFound {
			return nil, storage.ErrObjectNotFound
		}
		return nil, err
	}
	defer output.Close()

	return &storage.FileInfo{
		Size:    *output.ContentLength,
		ModTime: qs.TimeValue(output.LastModified),
		Mode:    0666,
	}, nil
}

func (s *qingStor) Size(key string) (int64, error) {
	if !storage.ValidKey(key) {
		return 0, storage.ErrObjectKeyInvalid
	}

	output, err := s.bucket.HeadObject(key, &qs.HeadObjectInput{})
	if err != nil {
		if qsErr, ok := err.(*qsErrors.QingStorError); ok && qsErr.StatusCode == http.StatusNotFound {
			return 0, storage.ErrObjectNotFound
		}
		return 0, err
	}

	if output.ContentLength == nil {
		return 0, fmt.Errorf("failed to get object size with code %d", qs.IntValue(output.StatusCode))
	}

	return *output.ContentLength, nil
}

func (s *qingStor) IsExist(key string) (bool, error) {
	if !storage.ValidKey(key) {
		return false, storage.ErrObjectKeyInvalid
	}

	_, err := s.bucket.HeadObject(key, &qs.HeadObjectInput{})
	if err != nil {
		if qsErr, ok := err.(*qsErrors.QingStorError); ok && qsErr.StatusCode == http.StatusNotFound {
			return false, storage.ErrObjectNotFound
		}
		return false, err
	}

	return true, nil
}

func (s *qingStor) Del(key string) error {
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	_, err := s.bucket.DeleteObject(key)
	if err != nil {
		if qsErr, ok := err.(*qsErrors.QingStorError); ok && qsErr.StatusCode == http.StatusNotFound {
			return storage.ErrObjectNotFound
		}
		return err
	}

	return nil
}
