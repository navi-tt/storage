package qingyun

import (
	"bytes"
	"github.com/navi-tt/storage"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var s storage.Storage

func TestMain(m *testing.M) {
	//qsCfg, err := config.New("BKENTRVGSRFLLEFQKEBV", "gwVNTMlJaT5JCpqnjDSKGtIbss5aT8dWn0uAI8vv")
	//if err != nil {
	//	return
	//}
	//
	//tt := &http.Transport{
	//	Proxy: http.ProxyFromEnvironment,
	//	DialContext: (&net.Dialer{
	//		Timeout:   30 * time.Second,
	//		KeepAlive: 30 * time.Second,
	//		DualStack: true,
	//	}).DialContext,
	//	TLSClientConfig: &tls.Config{
	//		InsecureSkipVerify: true,
	//	},
	//	ForceAttemptHTTP2:     true,
	//	MaxIdleConns:          100,
	//	IdleConnTimeout:       90 * time.Second,
	//	TLSHandshakeTimeout:   10 * time.Second,
	//	ExpectContinueTimeout: 10 * time.Second,
	//	MaxIdleConnsPerHost:   5,
	//}
	//
	////cookieJar, _ := cookiejar.New(nil)
	//qsCfg.Connection = &http.Client{
	//	Transport: tt,
	//	//Jar:       cookieJar,
	//	Timeout: time.Second * 3,
	//}
	//
	//qsCfg.Host = "10.8.135.10"
	//qsCfg.Port = 53
	//qsCfg.Protocol = "http"
	//
	////_ = qsCfg.LoadDefaultConfig()
	//
	//qsSvc, err := qs.Init(qsCfg)
	//if err != nil {
	//	fmt.Println("init error :", err.Error())
	//	return
	//}
	//
	//bucket, err := qsSvc.Bucket("pxbaimcdianxiao", "")
	//if err != nil {
	//	fmt.Println("init error : ", err.Error())
	//	return
	//}
	//
	//s = &qingStor{
	//	qsSvc:  qsSvc,
	//	bucket: bucket,
	//}

	s, _ = storage.Init("qos", `{
		"AccesskeyId":"BKENTRVGSRFLLEFQKEBV",    
		"SecretAccessKey":"gwVNTMlJaT5JCpqnjDSKGtIbss5aT8dWn0uAI8vv",
		"Zone":"",           
		"Bucket":"pxbaimcdianxiao",         
		"Protocol":"http",       
		"Host":"10.8.135.10",           
		"Port":"53",           
	}`)

	m.Run()
}

func TestQingStor_Put(t *testing.T) {
	//return
	data, _ := ioutil.ReadFile("info_result_2.txt")
	buf := bytes.NewBuffer(nil)
	_, _ = buf.Write(data)

	t.Logf("length : %d", buf.Len())

	err := s.Put("/kst/info_result_1028.txt", buf, int64(buf.Len()))
	if err != nil {
		t.Errorf("put error : %s", err.Error())
		return
	}

	t.Log("success")
}

func TestQingStor_Get(t *testing.T) {
	fd, err := os.OpenFile("info_result_1028_1.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	err = s.Get("/kst/info_result_1028.txt", fd)
	if err != nil {
		if err == storage.ErrObjectNotFound {
			t.Log("file not found")
		} else {
			t.Fatal(err)
		}
	}
	byts := bytes.NewBuffer(nil)
	io.Copy(byts, fd)

	t.Logf("success")

}

func TestQingStor_IsExist(t *testing.T) {
	isExisted, err := s.IsExist("/kst/info_result_1025.txt")
	if err != nil {
		t.Errorf("delete error : %s", err.Error())
		return
	}

	t.Logf("isExisted : %v", isExisted)
	t.Log("success")
}

func TestQingStor_Size(t *testing.T) {
	size, err := s.Size("/kst/info_result_1025.txt")
	if err != nil {
		t.Errorf("delete error : %s", err.Error())
		return
	}

	t.Logf("file size : %v", size)
	t.Log("success")
}

func TestQingStor_Del(t *testing.T) {
	err := s.Del("/kst/info_result_1025.txt")
	if err != nil {
		t.Errorf("delete error : %s", err.Error())
		return
	}

	t.Log("success")
}
