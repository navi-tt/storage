package main

import (
	"bytes"
	"fmt"
	"github.com/navi-tt/storage"
	_ "github.com/navi-tt/storage/fs"
)

func main() {
	if _, err := storage.Init(storage.FS, `{
		"BaseDir":"./testdata"
	}`); err != nil {
		fmt.Printf("err : %s\n", err.Error())
		return
	}

	fsStorage, err := storage.GetStorage("fs")
	if err != nil {
		fmt.Printf("err : %s\n", err.Error())
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString("this is a test filesadfsdafdfadfadfasfsadfadsfafdasdfadfadgfasdgfdasgasdfadsfadsfjkasjdfipasddhjflakfjas;ldkjas;oidja;sdlkfa;osdia;sdgjva;sfijsdao;fa;fkda;lksdgbha;lfkha;ldjgkasl;dfasdk;ldsjk;ghas;dfljsa;ghsdlkal;sdghjsdka;gj;nsd;lkj; ncs;lfjaslcnfjklscnjklnjlkcmndjnka;sijfxa;;kf;l ;cxmjak;lmklfcanm;lknkjnc lhjnjkhsnxkjnlkxcjfhkclafdjkxnalknxl!")

	if err := fsStorage.Put("test_fs_2.txt", buf, int64(buf.Len())); err != nil {
		fmt.Printf("err : %s\n", err.Error())
		return
	}

	if err := fsStorage.Del("test_fs_2.txt"); err != nil {
		fmt.Printf("err : %s\n", err.Error())
		return
	}
}
