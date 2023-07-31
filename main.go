package main

import (
	"chunk_upload/handlers"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	rootPath, _ := os.Getwd()
	router := httprouter.New()
	router.PUT("/chunk-upload", handlers.ChunkUpload())                  //分片上传路由
	router.POST("/chunk-merge", handlers.ChunkMerge())                   // 文件合并路由
	router.ServeFiles("/static/*filepath", http.Dir(rootPath+"/static")) // 文件服务
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("listen at ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
