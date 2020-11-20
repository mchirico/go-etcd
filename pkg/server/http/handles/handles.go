package handles

import (
	"bytes"
	"github.com/mchirico/zcovid/pkg/echarts/gauge"
	"github.com/mchirico/zcovid/pkg/echarts/heatmap"
	"github.com/mchirico/zcovid/pkg/echarts/line"

	"net/http"
)

var Count = 0
var CountStatus = 0

type HANDLE struct {
	Process func() string
}

func (h HANDLE) BaseRoot(w http.ResponseWriter, r *http.Request) {

	/*
		curl -H "Authorization: SomeToken" localhost:3000
		Value:  bob Revision:  1550
		Value:  555 Revision:  1551


		SomeToken
	*/

	reqToken := r.Header.Get("Authorization")

	switch r.Method {
	case "GET":
		Count += 1

		msg := h.Process()
		msg += "\n\n" + reqToken
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
