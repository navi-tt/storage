package storage

type storageType string

const (
	FS  storageType = "FS"
	COS storageType = "COS"
	QS  storageType = "QS"
)

// 本地文件系统存储服务
type CFG_FS struct {
	BaseDir string
}

// 青云存储服务
type CFG_QS struct {
	AccesskeyId     string
	SecretAccessKey string
	Zone            string
	Bucket          string
	Protocol        string
	Host            string
	Port            int
}

//腾讯云COS存储
type CFG_COS struct {
	SecretID  string
	SecretKey string
	Host      string
	Bucket    string
	Protocol  string
}
