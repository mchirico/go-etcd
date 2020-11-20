package handles

import (
	"bytes"
	"github.com/mchirico/go-etcd/pkg/etcdutils"
	"github.com/mchirico/zcovid/pkg/echarts/gauge"
	"github.com/mchirico/zcovid/pkg/echarts/heatmap"
	"github.com/mchirico/zcovid/pkg/echarts/line"


	"net/http"
)

var Count = 0
var CountStatus = 0

func BaseRoot(w http.ResponseWriter, r *http.Request) {

	//version := "v0.0.1"

	switch r.Method {
	case "GET":
		Count += 1
		e := etcdutils.NewETC("/certs")
		msg := e.EtcdRun()
		w.Write([]byte(msg))
	case "POST":
		// msg := fmt.Sprintf("Hello world: POST: %v", r.FormValue("user"))
		w.Write([]byte("post"))
	default:
		w.Write([]byte(`"Sorry, only GET and POST methods are supported."`))
	}

}

func Gauge(w http.ResponseWriter, r *http.Request) {

	g := gauge.GaugeExamples{}
	buf := bytes.NewBufferString("")
	g.Examples(buf)
	w.Write(buf.Bytes())

}

func Line(w http.ResponseWriter, r *http.Request) {
	e := line.LineExamples{}
	buf := bytes.NewBufferString("")
	e.Examples(buf)
	w.Write(buf.Bytes())

}

func Heatmap(w http.ResponseWriter, r *http.Request) {

	h := heatmap.HeatmapExamples{}
	buf := bytes.NewBufferString("")
	h.Examples(buf)
	w.Write(buf.Bytes())

}
