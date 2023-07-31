package handlers

import (
	"chunk_upload/cache"
	"chunk_upload/utils"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type ChunkMergeReq struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Chunk    int    `json:"chunk"`
}

//ChunkMerge
// 文件合并请求
func ChunkMerge() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// 获取body内容，{"filename": "sdafkslkdjf.zip", "hash": "adasdasdasd", "chunk": 10}
		b, _ := ioutil.ReadAll(r.Body)
		var req ChunkMergeReq
		if err := json.Unmarshal(b, &req); err != nil {
			log.Println(err)
			utils.Error(w, utils.ErrCodeBadReqeust, nil)
			return
		}
		// 创建文件
		f, err := os.Create("./uploads/" + req.Filename)
		if err != nil {
			log.Println(err)
			utils.Error(w, utils.ErrCodeSystemError, nil)
			return
		}
		var remove bool
		defer func() {
			f.Close()
			if remove {
				os.Remove("./uploads/" + req.Filename)
			}
		}()

		filecache, ok := cache.Get(req.Hash)
		if !ok {
			log.Println("未找到文件上传信息")
			remove = true
			utils.Error(w, utils.ErrCodeBadReqeust, nil)
			return
		}
		hasher := md5.New()
		// 根据分片数量，按照顺序遍历获取文件并进行合并
		for i := 0; i < req.Chunk; i++ {
			//按照顺序获取到分片的hash值
			chunkhash, ok := filecache[i]
			if !ok {
				log.Println("未找到分片文件", i)
				remove = true
				utils.Error(w, utils.ErrCodeBadReqeust, nil)
				return
			}
			//根据分片hash值，从文件夹中获取文件内容
			chunkfile := "./tmp/" + req.Hash + "/" + chunkhash
			cf, err := os.Open(chunkfile)
			if err != nil {
				log.Println("打开分片文件失败", err)
				remove = true
				utils.Error(w, utils.ErrCodeSystemError, nil)
				return
			}
			//使用teereader，copy的时候同时copy到hash
			tee := io.TeeReader(cf, hasher)
			_, err = io.Copy(f, tee)
			cf.Close() // 关闭文件资源
			if err != nil {
				log.Println(err)
				remove = true
				utils.Error(w, utils.ErrCodeSystemError, nil)
				return
			}
		}
		// 对上传文件进行hash验证
		fhash := fmt.Sprintf("%x", hasher.Sum(nil))
		if fhash != req.Hash {
			remove = true
			utils.Error(w, utils.ErrCodeHashNotMath, nil)
			return
		}
		cache.Del(req.Hash)
		utils.Success(w, nil)
	}
}
