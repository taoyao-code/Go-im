package middleware

import (
	"bytes"
	"log"
	"net/http"
	"reptile-go/util/vlog"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type ResponseWithRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *ResponseWithRecorder) WriteHeader(statusCode int) {
	rec.ResponseWriter.WriteHeader(statusCode)
	rec.statusCode = statusCode
}

func (rec *ResponseWithRecorder) Write(d []byte) (n int, err error) {
	n, err = rec.ResponseWriter.Write(d)
	if err != nil {
		return
	}
	rec.body.Write(d)

	return
}

// 日志记录
func AccessLogging(f http.Handler) http.Handler {
	// 创建一个新的handler包装http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//buf := new(bytes.Buffer)
		//buf.ReadFrom(r.Body)
		logEntry := vlog.AccessLog.WithFields(logrus.Fields{
			"ip":     r.RemoteAddr,
			"method": r.Method,
			"path":   r.RequestURI,
			"query":  r.URL.RawQuery,
			//"request_body": buf.String(),
			"request_body": r.PostForm.Encode(),
		})
		wc := &ResponseWithRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			//body:           bytes.Buffer{},
		}
		// 调用下一个中间件或者最终的handler处理程序
		f.ServeHTTP(wc, r)
		defer logEntry.WithFields(logrus.Fields{
			"status":        wc.statusCode,
			"response_body": wc.body.String(),
		}).Info()

	})
}

// 记录每个URL请求的执行时长
func Logging() mux.MiddlewareFunc {
	//	创建中间件
	return func(f http.Handler) http.Handler {
		//	创建一个新的handler包装http.HandlerFunc
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//	中间件的处理逻辑
			start := time.Now()
			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()
			// 调用下一个中间件或者最终的handler处理程序
			f.ServeHTTP(w, r)
		})
	}
}
