package handlers

import (
	"chunk_upload/cache"
	"chunk_upload/utils"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func ChunkUpload() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		chunkHash := r.Header.Get("X-Chunk-Hash") // 分片数据hash
		fileHash := r.Header.Get("X-File-Hash")   // 文件hash
		index := r.Header.Get("X-Chunk-Index")    // 分片序号
		// 判断请求头参数
		if chunkHash == "" || fileHash == "" || index == "" {
			utils.Error(w, utils.ErrCodeBadReqeust, nil)
			return
		}

		indexNum, err := strconv.Atoi(index)
		if err != nil {
			utils.Error(w, utils.ErrCodeBadReqeust, nil)
			return
		}
		//创建临时分片文件夹
		dir := "./tmp/" + fileHash
		_, err = os.Stat(dir)
		if os.IsNotExist(err) {
			//创建文件夹
			os.Mkdir(dir, os.ModePerm)
		}
		// 断点续传

		// 读取上传的file中的chunk内容
		f, _, err := r.FormFile("file")
		if err != nil {
			log.Println("文件上传失败")
			utils.Error(w, utils.ErrCodeBadReqeust, nil)
			return
		}

		// 创建文件准备开始写入
		chunkfile, err := os.Create(dir + "/" + chunkHash)
		if err != nil {
			log.Println(err)
			utils.Error(w, utils.ErrCodeSystemError, nil)
			return
		}
		var remove bool
		defer func() {
			chunkfile.Close()
			if remove {
				os.Remove(dir + "/" + chunkHash)
			}
		}()

		// 上传成功
		// 验证文件hash
		hasher := md5.New()
		tee := io.TeeReader(f, hasher)
		_, err = io.Copy(chunkfile, tee)
		if err != nil {
			log.Println(err)
			utils.Error(w, utils.ErrCodeSystemError, nil)
			return
		}
		fhash := fmt.Sprintf("%x", hasher.Sum(nil))
		if fhash != chunkHash {
			log.Println("文件hash不一致（文件已损坏）")
			// remove = true // hash验证失败，删除破损的文件
			utils.Error(w, utils.ErrCodeHashNotMath, nil)
			return
		}

		// 验证正确，对数据进行缓存
		cache.Set(fileHash, indexNum, chunkHash)
		utils.Success(w, nil)
	}
}
