package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"github.com/dustinmj/renotts/file"
	"github.com/dustinmj/renotts/player"
	"github.com/dustinmj/renotts/tmplt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// error messages
const (
	errInvalidJSON = "cannot read JSON data in POST body"
	errReadingBody = "could not retrieve POST body"
	errBadServce   = "service does not exist"
	errInvalidPath = "invalid path specified"
	errNoText      = "no text specified"
)

//Param http request structure
type Param struct {
	Text, Voice, SampleRate, Padding string
}

//Rq http request
type request struct {
	Typ    string
	Param  Param
	Unique []byte
	Body   []byte
}

//Rsp http response
type response struct {
	Msg   string
	Err   error
	Code  int
	Heads map[string]string
}

// how long to wait for a busy player before returning an error, in seconds
const busyTimeout = 5

// how long to wait after unsuccessful ip address resolution to try against, in seconds
const ipDelay = 6

// how many times to check for ip address resolution before giving up
const ipCheckLimit = 20

var mPath string
var mPort string
var ip string
var baseURI string
var ttsEndpoint string
var msgs map[string]string
var cfg config.Cfg

// Create - initialize server
func Create(port string, path string, conf config.Cfg) {
	mPort = port
	mPath = path
	cfg = conf
}

// Serve - start listing... blocking
func Serve() {
	ip = determineIP().String() // blocking
	if rsvd(mPath) {
		coms.Msg("Invalid path specified. ", mPath, " is reserved. Rewriting to /tts")
		mPath = "tts"
	}
	sMux := http.NewServeMux()
	sMux.HandleFunc("/", handler)
	var p string
	if mPort == "0" {
		listener, err := net.Listen("tcp", ":"+mPort)
		if err != nil {
			coms.Msg(err.Error())
		} else {
			// if mPort is 0, that's what upnp will advertise, this won't working
			// without adjusting the library, we can just get the port and
			// rebind manually
			p = fmt.Sprintf(":%v", listener.Addr().(*net.TCPAddr).Port)
			listener.Close()
		}
	} else {
		p = fmt.Sprintf(":%v", mPort)
	}
	mPort = p
	baseURI = fmt.Sprintf("http://%v%v", ip, p)
	ttsEndpoint = fmt.Sprintf("%v/%v/polly/", baseURI, mPath)
	StartUPNP(mPort) // create upnp server now that we know port
	coms.Msg(fmt.Sprintf("Instructions/Options: visit %v in a browser.", baseURI))
	if err := http.ListenAndServe(p, sMux); err != nil {
		coms.Exit(71, []byte("Cannot create webserver. "+err.Error()))
	}
}

func determineIP() net.IP {
	// check to make sure we can get outbound ip...
	ip := getOutboundIP()
	try := 0
	for ip.String() == "127.0.0.1" {
		if try < ipCheckLimit {
			coms.Msg(fmt.Sprintf("Could not reliably determine ip, trying again in %v seconds...", ipDelay))
			time.Sleep(time.Second * ipDelay) // block
			try++
			ip = getOutboundIP()
		} else {
			coms.Msg("Could not determine IP address. Giving up.")
			os.Exit(2)
		}
	}
	return ip
}

func logg(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// handle preflights
	if r.Header.Get("Access-Control-Request-Method") == http.MethodPost {
		sendOptions(w, r)
		return
	}
	switch r.RequestURI {
	case upnpdtp:
		devType(w, r)
		break
	case "/boot", "/boot/":
		showBootSetup(w, r, cfg)
		break
	case "/status", "/status/":
		status(w, r)
		break
	case "/services", "/services/":
		servicePath(w, r, cfg.Val(config.PATH))
		break
	case "/" + mPath + "/polly", "/" + mPath + "/polly/":
		tts(w, r, cfg)
		break
	case "/test/", "/test":
		printTest(w, r)
		break
	case "/check", "/check/":
		printChecks(w, r)
		break
	case "/":
		printPaths(w, r)
		break
	default:
		makeHead(w, http.StatusNotFound, "text/plain", "tts").Write([]byte("Endpoint not found. Please check your path configuration."))
	}
}

func sendOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(nil)
}

func printTest(w http.ResponseWriter, r *http.Request) {
	host := fmt.Sprintf("%v%v", ip, mPort)
	path := "/" + mPath + "/polly/"
	data := tmplt.TestData{
		Path: path,
		Host: host,
		Common: tmplt.Common{
			Title:   "RenoTTS Tester",
			BaseURI: baseURI}}
	tmplt.ParseHTM(w, tmplt.TestHTML, data)
}

func printPaths(w http.ResponseWriter, r *http.Request) {
	var dat [][]string
	dat = append(dat, []string{"Test Interface", fmt.Sprintf("%v/test/", baseURI)})
	dat = append(dat, []string{"Boot Configuration Instructions", fmt.Sprintf("%v/boot/", baseURI)})
	dat = append(dat, []string{"Check Configuration", fmt.Sprintf("%v/check/", baseURI)})
	dat = append(dat, []string{"Ping Status", fmt.Sprintf("%v/status/", baseURI)})
	dat = append(dat, []string{"List Services", fmt.Sprintf("%v/services/", baseURI)})
	dat = append(dat, []string{"UPnP Device Description", fmt.Sprintf("%v%v", baseURI, upnpdtp)})
	dat = append(dat, []string{"TTS Endpoint", ttsEndpoint})
	data := tmplt.URLList{
		Data: dat,
		Common: tmplt.Common{
			Title:   "Available URI Endpoints:",
			BaseURI: baseURI}}
	tmplt.ParseHTM(w, tmplt.URLListHTML, data)
}

func printChecks(w http.ResponseWriter, r *http.Request) {
	out := config.ConfigChk.All()
	data := tmplt.List{
		Data: out,
		Common: tmplt.Common{
			Title:   "Configuration Check",
			BaseURI: baseURI}}
	tmplt.ParseHTM(w, tmplt.ListHTML, data)
}

func servicePath(w http.ResponseWriter, r *http.Request, path string) {
	s := map[string]string{}
	// show services
	for k := range AvailServs {
		s[k] = "/" + path + "/" + k + "/"
	}
	j, _ := json.Marshal(s)
	makeHead(w, http.StatusOK, "application/json", "services").Write(j)
}

func showBootSetup(w http.ResponseWriter, r *http.Request, cfg config.Cfg) {
	logFile := "/tmp/RenoTTS.log"
	data := tmplt.BootData{
		User:           cfg.User(),
		ConfigFile:     cfg.File(),
		AppPath:        cfg.AppPath(),
		LogFile:        logFile,
		ConfigCheckURL: fmt.Sprintf("%v/check/", baseURI),
		TestURL:        fmt.Sprintf("%v/test/", baseURI),
		ServiceFile:    file.SystemdPath,
		Common: tmplt.Common{
			Title:   "Startup Configuration",
			BaseURI: baseURI}}
	tmplt.ParseHTM(w, tmplt.BootHTML, data)
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
	data := tmplt.UpnpDevDesc{
		DeviceType:   upnpdt,
		FriendlyName: upnpfn,
		Manufacturer: upnpman,
		ModelNumber:  upnpmod,
		UUID:         upnpuuid,
		Services: []tmplt.UpnpService{
			tmplt.UpnpService{
				Path: ttsEndpoint}}}
	w = makeHead(w, http.StatusOK, "application/xml", "device-type")
	tmplt.ParseF(w, tmplt.DevDescFl, data)
}

func tts(w http.ResponseWriter, r *http.Request, cfg config.Cfg) {
	defer r.Body.Close()
	t, err := rtype(r)
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
	file, err := eN.Query(rQ, cfg)
	if err != nil {
		reply(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	// determine if we're padding
	before, after := pad(rQ)
	// if required, execute with player
	if eN.Caches() {
		p := player.GetPlayer(cfg)
		if err := playFile(p, *file, before, after, cfg); err != nil {
			reply(w, http.StatusFailedDependency, "error")
			return
		}
	}
	reply(w, http.StatusOK, "success")
}

func playFile(mpgPlayer player.SPlayer, file string, before bool, after bool, cfg config.Cfg) error {
	// if busy, queue the file for later
	if mpgPlayer.Busy() {
		mpgPlayer.Queue(file, before, after)
		return nil
	}
	if err := mpgPlayer.Play(file, before, after, cfg.Cache()); err != nil {
		coms.Msg("Unable to play ", file, err.Error())
		return err
	}
	return nil
}

func pad(req *request) (before bool, after bool) {
	switch req.Param.Padding {
	case "Both":
		before = true
		after = true
	case "Before":
		before = true
		break
	case "After":
		after = true
		break
	}
	return
}

func rtype(in *http.Request) (string, error) {
	p := in.URL.Path[len(mPath)+1:]
	t := filepath.Base(p)
	// check for extra content in path (2 slashes)...
	if len(p) > len(t)+2 {
		return "", errors.New(errInvalidPath)
	}
	return t, nil
}

func mk(in *http.Request, t string) (*request, error) {
	bd, err := ioutil.ReadAll(in.Body)
	if err != nil {
		return nil, errors.New(errReadingBody)
	}
	out := request{Typ: t, Body: bd}
	err = json.Unmarshal(bd, &out.Param)
	if err != nil {
		return nil, err
	}
	out.Param.Text = fmt.Sprintf("%s", out.Param.Text)
	// trim text to 3k chars
	if len(out.Param.Text) > 3000 {
		out.Param.Text = out.Param.Text[:3000]
	} else if len(out.Param.Text) < 1 {
		return nil, errors.New(errNoText)
	}
	// set unique
	p := out.Param
	p.Padding = ""
	p.Text = strings.ToLower(p.Text)
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	out.Unique = b
	return &out, nil
}

func rsvd(p string) bool {
	switch p {
	case "", "status", "services", "check", "test", "boot":
		return true
	}
	return false
}

type msg struct {
	Message string
}

func reply(w http.ResponseWriter, code int, txt string) {
	t := msg{
		Message: txt}
	b, _ := json.Marshal(t)
	makeHead(w, code, "application/json", "tts").Write(b)
}

func addHead(w http.ResponseWriter, h map[string]string) http.ResponseWriter {
	for k, v := range h {
		w.Header().Set(k, v)
	}
	return w
}

func makeHead(w http.ResponseWriter, c int, t string, a string) http.ResponseWriter {
	w.Header().Set("Server", coms.AppName)
	w.Header().Set("Content-Type", t)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", getOpenHeaders())
	w.Header().Set("Action", a)
	w.Header().Set("Status", getStatus(c))
	w.WriteHeader(c)
	return w
}

func getOpenHeaders() string {
	return "Server,Content-Type,Access-Control-Expose-Headers,Access-Control-Allow-Origin,Action,Status,Content-Length,Date"
}

func getStatus(code int) string {
	return strconv.Itoa(code) + " " + string(http.StatusText(code))
}

//getOutboundIP - Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.ParseIP("127.0.0.1")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

//getMacAddr - Get mac address of something on this system
func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
