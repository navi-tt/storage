package cos

import (
	"context"
	"fmt"
	"github.com/navi-tt/storage"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"time"
)

//腾讯云COS存储
type COS struct {
	SecretID  string
	SecretKey string
	Host      string
	Bucket    string
	Protocol  string
}

func Init(cfg *COS) error {
	storageUrl := fmt.Sprintf("%s://%s.%s", cfg.Protocol, cfg.Bucket, cfg.Host)

	u, _ := url.Parse(storageUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	storage.Register("COS", &tensenCos{
		client: c,
	})
	return nil
}

type tensenCos struct {
	client *cos.Client
}

func (s *tensenCos) Put(key string, r io.Reader, contentLength int64) error {
	fmt.Printf("[COS PUT] object: %s \n", key)

	opt := &cos.ObjectPutOptions{}
	_, err := s.client.Object.Put(context.Background(), key, r, opt)
	if err != nil {
		return err
	}

	return nil
}

func (s *tensenCos) Get(key string, wa io.WriterAt) error {
	fmt.Printf("[COS GET] object: %s \n", key)

	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	opt := &cos.ObjectGetOptions{}

	v, err := s.client.Object.Get(context.Background(), key, opt)
	if err != nil {
		return err
	}

	return storage.Copy(wa, v.Body)
}

func (s *tensenCos) FileStream(key string) (io.ReadCloser, *storage.FileInfo, error) {
	if !storage.ValidKey(key) {
		return nil, nil, storage.ErrObjectKeyInvalid
	}

	opt := &cos.ObjectGetOptions{}

	output, err := s.client.Object.Get(context.Background(), key, opt)
	if err != nil {
		if output != nil {
			if output.StatusCode == http.StatusNotFound {
				return nil, nil, storage.ErrObjectNotFound
			}
		}
		return nil, nil, err
	}

	modTime, _ := time.Parse(time.RFC1123, output.Header.Get("Last-Modified"))

	return output.Body, &storage.FileInfo{
		Size:    output.ContentLength,
		ModTime: modTime.In(time.Local),
		Mode:    0666,
	}, nil
}

func (s *tensenCos) Stat(key string) (*storage.FileInfo, error) {
	if !storage.ValidKey(key) {
		return nil, storage.ErrObjectKeyInvalid
	}

	opt := &cos.ObjectGetOptions{}

	output, err := s.client.Object.Get(context.Background(), key, opt)
	if err != nil {
		if output != nil {
			if output.StatusCode == http.StatusNotFound {
				return nil, storage.ErrObjectNotFound
			}
		}
		return nil, err
	}

	modTime, _ := time.Parse(time.RFC1123, output.Header.Get("Last-Modified"))

	return &storage.FileInfo{
		Size:    output.ContentLength,
		ModTime: modTime.In(time.Local),
		Mode:    0666,
	}, nil
}

func (s *tensenCos) Size(key string) (int64, error) {
	if !storage.ValidKey(key) {
		return 0, storage.ErrObjectKeyInvalid
	}

	opt := &cos.ObjectGetOptions{}

	output, err := s.client.Object.Get(context.Background(), key, opt)
	if err != nil {
		if output != nil {
			if output.StatusCode == http.StatusNotFound {
				return 0, storage.ErrObjectNotFound
			}
		}
		return 0, err
	}

	if output.ContentLength == 0 {
		return 0, fmt.Errorf("failed to get object size with code %d", output.StatusCode)
	}

	return output.ContentLength, nil
}

func (s *tensenCos) IsExist(key string) (bool, error) {
	if !storage.ValidKey(key) {
		return false, storage.ErrObjectKeyInvalid
	}

	opt := &cos.ObjectGetOptions{}

	output, err := s.client.Object.Get(context.Background(), key, opt)
	if err != nil {
		if output != nil {
			if output.StatusCode == http.StatusNotFound {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

func (s *tensenCos) Del(key string) error {
	fmt.Printf("[COS DEL] object: %s \n", key)
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	resp, err := s.client.Object.Delete(context.Background(), key)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return storage.ErrObjectNotFound
	}

	return nil
}
