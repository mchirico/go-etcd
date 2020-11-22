package socket

import "net/http"

type HANDLE struct {
	broadcast chan []byte
}

func (h HANDLE) BaseRoot(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		msg := "here"
		w.Write([]byte(msg))
		h.broadcast <- []byte("Sending message from base\n")
	case "POST":
		// msg := fmt.Sprintf("Hello world: POST: %v", r.FormValue("user"))
		w.Write([]byte("post"))
	default:
		w.Write([]byte(`"Sorry, only GET and POST methods are supported."`))
	}

}
