package cos

import (
	"fmt"
	"github.com/navi-tt/storage"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"testing"
)

var s storage.Storage

func TestMain(m *testing.M) {
	storageUrl := fmt.Sprintf("http://%s.%s", "kst-zxm-1304077072", "cos.ap-beijing.myqcloud.com")

	u, _ := url.Parse(storageUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "AKID74k9NB8NkEJCoCsAk7ipgHbhyzkI6ZIv",
			SecretKey: "13oIVW5NN9eutfpkOltFJwIM9gwwNWNu",
		},
	})

	storage.Register(storage.COS,&tensenCos{
		client: c,
	})

	s = &tensenCos{
		client: c,
	}
	m.Run()
}

func TestTensenCos_Put(t *testing.T) {
	err := storage.PutByPath(storage.COS,"/kst/test.wav", "../testdata/test_wav.wav")
	if err != nil {
		t.Fatalf("put error : %s", err.Error())
		return
	}

	t.Logf("finish")
}

func TestTensenCos_Get(t *testing.T) {
	_, _, err := storage.FileStream(storage.COS,"/kst/test.wav")
	if err != nil {
		t.Fatalf("get error : %s", err.Error())
		return
	}

	t.Logf("finish")
}

func TestTensenCos_IsExist(t *testing.T) {
	existed, err := storage.IsExist(storage.COS,"/kst/test.txt")
	if err != nil {
		t.Fatalf("is existed error : %s", err.Error())
		return
	}

	t.Logf("existed : %v", existed)
}

func TestTensenCos_Size(t *testing.T) {
	size, err := storage.Size(storage.COS,"/kst/test.txt")
	if err != nil {
		t.Fatalf("size error : %s", err.Error())
		return
	}

	t.Logf("size : %v", size)
}

func TestTensenCos_Del(t *testing.T) {
	err := storage.Del(storage.COS,"/kst/test.txt")
	if err != nil {
		t.Fatalf("delete error : %s", err.Error())
		return
	}

	t.Logf("finish")
}
