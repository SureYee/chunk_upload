package handlers

import (
	"chunk_upload/cache"
	"chunk_upload/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ChunkPrepare() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		hash := p.ByName("hash")
		c, _ := cache.Get(hash)
		data := make([]int, 0, len(c))
		for k := range c {
			data = append(data, k)
		}
		utils.Success(w, data)
	}
}
