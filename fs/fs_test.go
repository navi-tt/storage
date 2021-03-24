package fs

import (
	"bytes"
	"github.com/navi-tt/storage"
	"testing"
)

var s storage.Storage

func TestMain(m *testing.M) {
	//s = &fs{
	//	baseDir: "../testdata",
	//}

	s, _ = storage.Init("fs", `{
		"BaseDir":"../testdata"
	}`)

	m.Run()
}

func TestFs_Put(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("this is a test filesadfsdafdfadfadfasfsadfadsfafdasdfadfadgfasdgfdasgasdfadsfadsfjkasjdfipasddhjflakfjas;ldkjas;oidja;sdlkfa;osdia;sdgjva;sfijsdao;fa;fkda;lksdgbha;lfkha;ldjgkasl;dfasdk;ldsjk;ghas;dfljsa;ghsdlkal;sdghjsdka;gj;nsd;lkj; ncs;lfjaslcnfjklscnjklnjlkcmndjnka;sijfxa;;kf;l ;cxmjak;lmklfcanm;lknkjnc lhjnjkhsnxkjnlkxcjfhkclafdjkxnalknxl!")
	err := s.Put("test_fs.txt", buf, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("succ")
}

func TestFs_PutByPath(t *testing.T) {
	err := s.PutByPath("test_fs_put_by_path.txt", "../testdata/test_fs.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("succ")
}

func TestFs_Get(t *testing.T) {
	err := s.GetToPath("test_fs_put_by_path.txt", "../testdata/test_get_to_path.txt")
	if err != nil {
		if err == storage.ErrObjectNotFound {
			t.Log("file not found")
		} else {
			t.Fatal(err)
		}
	}
	t.Log("succ")
}

func TestFs_Size(t *testing.T) {
	size, err := s.Size("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(size)
}

func TestFs_Exists(t *testing.T) {
	exist, err := s.IsExist("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}

func TestFs_Del(t *testing.T) {
	err := s.Del("test_get_to_path.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("deleted")

	exist, err := s.IsExist("test_get_to_path.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
