package serv

import (
	"encoding/json"
	"errors"
	"github.com/dustinmj/renotts/aud"
	"github.com/dustinmj/renotts/com"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/service"
	"github.com/dustinmj/renotts/upnp"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

var mPath string
var mPort string
var msgs map[string]string

// Create - start a listening server on port/path
// p: port to listen, must include (eg :8080)
// path: path to set up tts server at (eg /tts)
func Create() {
	// TODO error checking in port,path
	mPort = config.Val("port")
	mPath = config.Val("path")
	if rsvd(mPath) {
		com.Msg("Invalid path specified. ", mPath, " is reserved. Rewriting to /tts")
		mPath = "tts"
	}
	com.Msg("Server listening at http://" + com.GetOutboundIP().String() + mPort + "/" + mPath + "/polly/")
	com.Msg("Instructions on use at http://" + com.GetOutboundIP().String() + mPort + "/")
	sMux := http.NewServeMux()
	sMux.HandleFunc("/", handler)
	if err := http.ListenAndServe(mPort, sMux); err != nil {
		com.Exit(71, []byte("Cannot create webserver. "+err.Error()))
	}
}

func logg(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.RequestURI {
	case "/", "/help", "/help/":
		w.Write(config.Instruct)
		break
	case upnp.DVPATH:
		devType(w, r)
	case "/status", "/status/":
		status(w, r)
		break
	case "/services", "/services/":
		servicePath(w, r)
		break
	case "/" + mPath + "/polly", "/" + mPath + "/polly/":
		tts(w, r)
		break
	default:
		makeHead(w, http.StatusNotFound, "text/plain", "tts").Write([]byte("Endpoint not found. Please check your path configuration."))
	}
}

func instruct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Write(config.Instruct)
}

func servicePath(w http.ResponseWriter, r *http.Request) {
	s := map[string]string{}
	// show services
	for k := range service.AvailServs {
		s[k] = "/" + config.Val("path") + "/" + k + "/"
	}
	j, _ := json.Marshal(s)
	makeHead(w, http.StatusOK, "application/json", "services").Write(j)
}

func status(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type re struct {
		Status string
	}
	resp, _ := json.Marshal(&re{Status: "200 OK"})
	makeHead(w, http.StatusOK, "application/json", "status").Write(resp)
}

func devType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	makeHead(w, http.StatusOK, "application/xml", "device-type").Write(upnp.GetDD())
}

func tts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	t, err := typ(r)
	if err != nil {
		reply(w, http.StatusBadRequest, err.Error())
		return
	}
	eN, err := enGet(t)
	if err != nil {
		reply(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	rQ, err := mk(r, t)
	if err != nil {
		reply(w, http.StatusBadRequest, err.Error())
		return
	}
	sF, rsp := eN.Query(&rQ)
	if rsp.Err != nil {
		w = addHead(w, rsp.Heads)
		reply(w, http.StatusMethodNotAllowed, rsp.Err.Error())
		return
	}
	if eN.Caches() {
		//go aud.Play(sF)
		if err := aud.Play123(sF); err != nil {
			reply(w, http.StatusInternalServerError, err.Error())
		}
	}
	reply(w, rsp.Code, rsp.Msg)
}

func typ(in *http.Request) (string, error) {
	p := in.URL.Path[len(mPath)+1:]
	t := filepath.Base(p)
	// check for extra content in path (2 slashes)...
	if len(p) > len(t)+2 {
		return "", errors.New(config.Err["InvalidPath"])
	}
	return t, nil
}

func mk(in *http.Request, t string) (com.Rq, error) {
	bd, err := ioutil.ReadAll(in.Body)
	if err != nil {
		return com.Rq{}, errors.New(config.Err["ErrorReadingBody"])
	}
	out := com.Rq{Typ: t, Body: bd}
	err = json.Unmarshal(bd, &out.Param)
	if err != nil {
		return com.Rq{}, err
	}
	// trim text to 3k chars
	if len(out.Param.Text) > 3000 {
		out.Param.Text = out.Param.Text[:3000]
	}
	return out, nil
}

func rsvd(p string) bool {
	switch p {
	case "", "status", "services", "help":
		return true
	}
	return false
}

func reply(w http.ResponseWriter, code int, txt string) {
	makeHead(w, code, "text/plain", "tts").Write([]byte(txt))
}

func addHead(w http.ResponseWriter, h map[string]string) http.ResponseWriter {
	for k, v := range h {
		w.Header().Set(k, v)
	}
	return w
}

func makeHead(w http.ResponseWriter, c int, t string, a string) http.ResponseWriter {
	w.Header().Set("Server", com.AppName)
	w.Header().Set("Content-Type", t)
	w.Header().Set("Action", a)
	w.Header().Set("Status", getStatus(c))
	//w.Header().Set("Content-Location", config.Val("path") + "/polly/")
	w.WriteHeader(c)
	return w
}

func getStatus(code int) string {
	return strconv.Itoa(code) + " " + string(http.StatusText(code))
}
