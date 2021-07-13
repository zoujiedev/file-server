package handler

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	redislayer "github.com/zoujiepro/file-server/cache/redis"
	"github.com/zoujiepro/file-server/db"
	"github.com/zoujiepro/file-server/utils"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户参数
	token := r.Header.Get("token")
	username, success := db.GetUsernameByToken(token)
	if !success {
		fmt.Printf("未知的token，获取用户信息失败")
		utils.WriteFail(w, "获取用户信息失败")
		return
	}

	r.ParseForm()
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		utils.WriteFail(w, "the parameter filesize is invalid")
		return
	}
	//2. 获得redis连接
	redisConn := redislayer.RedisPool().Get()
	defer redisConn.Close()

	//3. 生成分块上传初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadId:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize / (5 * 1024 * 1024)))),
	}

	//3. 将初始化信息写入redis
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadId, "chunkcount", uploadInfo.ChunkCount)
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadId, "filehash", uploadInfo.FileHash)
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadId, "filesize", uploadInfo.FileSize)

	//4. 将初始化信息响应到客户端
	utils.WriteSuccess(w, uploadInfo)
}

func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析参数
	token := r.Header.Get("token")
	_, success := db.GetUsernameByToken(token)
	if !success {
		utils.WriteFail(w, "获取用户信息失败")
		return
	}

	r.ParseForm()
	uploadId := r.Form.Get(" ")
	chunkIndex := r.Form.Get("index")

	//2. 获取redis连接
	conn := redislayer.RedisPool().Get()
	defer conn.Close()

	//3. 获取文件句柄，存储分块文件内容
	fpath := "D:/tmp/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0700)
	fd, err := os.Create(fpath)
	if err != nil {
		fmt.Printf("create[%s] err: %s\n", fpath, err.Error())
		utils.WriteFail(w, "upload fail")
		return
	}
	defer fd.Close()

	buffer := make([]byte, 1024*1024)
	for {
		read, err := r.Body.Read(buffer)
		fd.Write(buffer[:read])
		if err != nil {
			break
		}
	}

	//4. 更新redis缓存
	conn.Do("HSET", "MP_"+uploadId, "chkidx_"+chunkIndex, 1)

	//5. 返回处理结果给客户端
	utils.WriteSuccess(w, nil)
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析参数
	token := r.Header.Get("token")
	username, success := db.GetUsernameByToken(token)
	if !success {
		utils.WriteFail(w, "获取用户信息失败")
		return
	}

	r.ParseForm()
	uploadId := r.Form.Get("uploadid")
	fileHash := r.Form.Get("filehash")
	fileSize := r.Form.Get("filesize")
	fileName := r.Form.Get("filename")
	//2. 获得redis连接池连接
	conn := redislayer.RedisPool().Get()
	defer conn.Close()

	//3. 通过uploadId查询redis并判断是否所有块分块上传完成
	data, err := redis.Values(conn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		fmt.Printf("query upload info from redis err: %s", err.Error())
		utils.WriteFail(w, "complete upload fail")
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, err = strconv.Atoi(v)
			if err != nil {
				fmt.Printf("redis totalCount convert err: %s", err.Error())
				utils.WriteFail(w, "complete upload fail")
				return
			}
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		fmt.Printf("totalConut[%d] not equal chunkCount[%d]", totalCount, chunkCount)
		utils.WriteFail(w, "complete upload fail")
		return
	}

	//4. todo 分块合并
	//5. 更新唯一文件表和用户文件表
	ffileSize, err := strconv.Atoi(fileSize)
	db.UploadFile(fileHash, fileName, int64(ffileSize), "")

	db.InsertUserFile(username, fileHash, int64(ffileSize), fileName)
	//6. 响应处理结果
	utils.WriteSuccess(w, nil)
}

func CancelUploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1. 删除已经上传的文件块
	//2. 删除redis缓存状态
	//3. 更新mysql文件status
}

func MultipartUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	//1. 检查分块上传是否有效
	//2. 获取分块初始化信息
	//3. 获取已上传的分块信息
}
