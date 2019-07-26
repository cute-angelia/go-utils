package api

import (
	"net/http"
	"encoding/json"
	"log"
)

var ApiMakeLog *ApiLog

var filterPaths = []string{
	"/log/operation/listsByPage",
}

type ApiLog struct {
	create chan interface{}
	stop   chan struct{}
}

func NewApiLog() *ApiLog {
	h := &ApiLog{
		create: make(chan interface{}),
	}
	return h
}

func (h *ApiLog) Start() {
	go h.loop()
}

func (h *ApiLog) Stop() {
	close(h.stop)
}

func (h *ApiLog) loop() {
	for {
		select {
		case p := <-h.create:
			h.insertLog(p)
		case <-h.stop:
			return
		}
	}
}

func (h *ApiLog) insertLog(createL interface{}) {
	//apiClient := operationLogPb.NewOperationService(core.SRV_MANAGE, client.DefaultClient)
	//apiClient.Create(context.TODO(), createL)
	log.Println("to do insert")
}

// api 调用
func (h *ApiLog) CreateLog(r *http.Request, data map[string]interface{}) {
	uid := r.Header.Get("uid")

	r.ParseForm()

	requestParams, _ := json.Marshal(r.PostForm)
	response, _ := json.Marshal(data)

	path := r.URL.Path

	for _, v := range filterPaths {
		if v == path {
			return
		}
	}

	createL := map[string]string{
		"OperatorUserId": uid,
		"RequestPath":    r.URL.Path,
		"RequestIp":      r.RemoteAddr,
		"RequestParams":  string(requestParams),
		"Response":       string(response),
	}

	h.create <- createL
}
