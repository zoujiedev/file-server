package meta

import (
	"github.com/zoujiepro/file-server/db"
)

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

func UpdateFileMetaDB(meta FileMeta) bool {
	return db.UploadFile(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
}

func GetFileMeta(sha1 string) FileMeta {
	return fileMetas[sha1]
}

func GetFileMetaDB(sha1 string) (FileMeta, error) {
	tableFile, err := db.GetFileMeta(sha1)
	if err != nil {
		return FileMeta{}, err
	}
	metaInfo := FileMeta{
		FileSha1: tableFile.FileHash,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAddr.String,
	}
	return metaInfo, nil
}

func DeleteFileMeta(sha1 string) {
	delete(fileMetas, sha1)
}
