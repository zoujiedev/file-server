package meta

type FileMeta struct {
	FileSha1   string
	FileName   string
	FileSize   int64
	Location   string
	UpdateTime string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(meta FileMeta) {
	fileMetas[meta.FileSha1] = meta
}

func GetFileMeta(sha1 string) FileMeta {
	return fileMetas[sha1]
}

func DeleteFileMeta(sha1 string) {
	delete(fileMetas, sha1)
}
