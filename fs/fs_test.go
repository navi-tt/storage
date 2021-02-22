package fs

import (
	"bytes"
	"github.com/navi-tt/storage"
	"io"
	"os"
	"testing"
)

var s storage.Storage

func TestMain(m *testing.M) {
	s = &fs{
		baseDir: "../testdata",
	}
	m.Run()
}

func TestPut(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("this is a test filesadfsdafdfadfadfasfsadfadsfafdasdfadfadgfasdgfdasgasdfadsfadsfjkasjdfipasddhjflakfjas;ldkjas;oidja;sdlkfa;osdia;sdgjva;sfijsdao;fa;fkda;lksdgbha;lfkha;ldjgkasl;dfasdk;ldsjk;ghas;dfljsa;ghsdlkal;sdghjsdka;gj;nsd;lkj; ncs;lfjaslcnfjklscnjklnjlkcmndjnka;sijfxa;;kf;l ;cxmjak;lmklfcanm;lknkjnc lhjnjkhsnxkjnlkxcjfhkclafdjkxnalknxl!")
	err := s.Put("test_fs.txt", buf, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("succ")
}

func TestGet(t *testing.T) {
	fd, err := os.OpenFile("./tmp/get_test", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	err = s.Get("test.txt", fd)
	if err != nil {
		if err == storage.ErrObjectNotFound {
			t.Log("file not found")
		} else {
			t.Fatal(err)
		}
	}
	byts := bytes.NewBuffer(nil)
	io.Copy(byts, fd)
	t.Log(byts.String())
	t.Log("succ")
}

func TestSize(t *testing.T) {
	size, err := s.Size("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(size)
}

func TestExists(t *testing.T) {
	exist, err := s.IsExist("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
