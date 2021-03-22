package cos

import (
	"bytes"
	"github.com/navi-tt/storage"
	"io/ioutil"
	"testing"
)

var s storage.Storage

func TestMain(m *testing.M) {
	s, _ = storage.Init("cos", `{
		"SecretID":"AKID74k9NB8NkEJCoCsAk7ipgHbhyzkI6ZIv" 
		"SecretKey":"13oIVW5NN9eutfpkOltFJwIM9gwwNWNu"
		"Host":"cos.ap-beijing.myqcloud.com"     
		"Bucket":"kst-zxm-1304077072"   
		"Protocol":"http" 
	}`)

	m.Run()
}

func TestTensenCos_Put(t *testing.T) {

	data, _ := ioutil.ReadFile("info_result_2.txt")
	buf := bytes.NewBuffer(nil)
	_, _ = buf.Write(data)

	t.Logf("length : %d", buf.Len())

	err := s.Put("/kst/test.wav", buf, int64(buf.Len()))
	if err != nil {
		t.Fatalf("put error : %s", err.Error())
		return
	}

	t.Logf("finish")
}

func TestTensenCos_Get(t *testing.T) {
	_, _, err := s.FileStream("/kst/test.wav")
	if err != nil {
		t.Fatalf("get error : %s", err.Error())
		return
	}

	t.Logf("finish")
}

func TestTensenCos_IsExist(t *testing.T) {
	existed, err := s.IsExist("/kst/test.txt")
	if err != nil {
		t.Fatalf("is existed error : %s", err.Error())
		return
	}

	t.Logf("existed : %v", existed)
}

func TestTensenCos_Size(t *testing.T) {
	size, err := s.Size("/kst/test.txt")
	if err != nil {
		t.Fatalf("size error : %s", err.Error())
		return
	}

	t.Logf("size : %v", size)
}

func TestTensenCos_Del(t *testing.T) {
	err := s.Del("/kst/test.txt")
	if err != nil {
		t.Fatalf("delete error : %s", err.Error())
		return
	}

	t.Logf("finish")
}
