package utils

import (
	"encoding/json"
	"net/http"
)

const (
	ErrCodeOk          = 1
	ErrCodeBadReqeust  = 2
	ErrCodeHashNotMath = 3
	ErrCodeUploadFail  = 4 //分片文件上传失败
	ErrCodeSystemError = 5
)

var errinfo = map[int]string{
	ErrCodeOk:          "上传成功",
	ErrCodeBadReqeust:  "请求参数错误",
	ErrCodeHashNotMath: "文件hash验证失败",
	ErrCodeUploadFail:  "文件上传失败",
	ErrCodeSystemError: "系统错误",
}

var errStatus = map[int]int{
	ErrCodeOk:          http.StatusOK,
	ErrCodeBadReqeust:  http.StatusBadRequest,
	ErrCodeHashNotMath: http.StatusBadRequest,
	ErrCodeUploadFail:  http.StatusBadRequest,
	ErrCodeSystemError: http.StatusInternalServerError,
}

func Success(w http.ResponseWriter) {
	Error(w, ErrCodeOk)
}

func Error(w http.ResponseWriter, code int) {
	b, _ := json.Marshal(map[string]any{
		"code":    code,
		"message": errinfo[code],
	})
	w.WriteHeader(errStatus[code])
	w.Write(b)
}
