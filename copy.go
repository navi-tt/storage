package storage

import (
	"io"
)

func Copy(wa io.WriterAt, r io.Reader) error {
	var (
		buf       = make([]byte, 1024)
		off int64 = 0
	)

	for {
		n, err1 := r.Read(buf)
		buf := buf[:n]

		n, err2 := wa.WriteAt(buf, off)
		if err2 != nil {
			return err2
		}

		off += int64(n)
		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			return err1
		}
	}
	return nil
}
