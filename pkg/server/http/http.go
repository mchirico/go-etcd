package http

import (
	"github.com/mchirico/go-etcd/pkg/etcdutils"
	"github.com/mchirico/go-etcd/pkg/server/http/handles"
	"log"
	"net/http"
	_ "time/tzdata"
)

func SetupHandles() {

	h := handles.HANDLE{}
	e := etcdutils.NewETC("/certs")
	h.Process = e.EtcdRun

	http.HandleFunc("/", h.BaseRoot)
	http.HandleFunc("/gauge", handles.Gauge)
	http.HandleFunc("/line", handles.Line)
	http.HandleFunc("/heatmap", handles.Heatmap)

}

func Server() {
	SetupHandles()
	log.Println("starting server... :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
