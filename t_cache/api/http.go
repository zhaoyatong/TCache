package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"tcache/control"
)

const defaultBasePath = "/_tcache/"

// HTTPPool HTTP池
type HTTPPool struct {
	host     string
	basePath string
}

// NewHTTPPool 初始化HTTP池
func NewHTTPPool(host string) *HTTPPool {
	return &HTTPPool{
		host:     host,
		basePath: defaultBasePath,
	}
}

// Log 日志信息
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.host, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	cacheName := parts[0]
	key := parts[1]

	// 获取key
	if r.Method == "GET" {
		tCache := control.GetTCache(cacheName)
		if tCache == nil {
			http.Error(w, "no such group: "+cacheName, http.StatusNotFound)
			return
		}

		view, err := tCache.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")

		// 设置 Content-Type 为 JSON
		w.Header().Set("Content-Type", "application/json")

		// 直接写入 JSON（自动转为 []byte）
		result := map[string]string{
			"key":   key,
			"value": string(view.ByteSlice()),
		}
		json.NewEncoder(w).Encode(result)
	}
	// 添加/修改key
	if r.Method == "POST" {
		// 读取整个请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close() // 重要: 确保关闭

		tCache := control.GetTCache(cacheName)
		if tCache == nil {
			tCache = control.NewTCache(cacheName)
		}

		err = tCache.Set(key, body)
		if err != nil {
			http.Error(w, cacheName+key+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
