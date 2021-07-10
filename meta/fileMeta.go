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

func UpdateFileMetaDB(meta FileMeta) bool {
	return db.UploadFile(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
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

func DeleteFileMetaDB(sha1 string) bool {
	return db.DeleteFileMeta(sha1)
}
