package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
)

func StartTestServer(handler func(w http.ResponseWriter, r *http.Request)) (server *httptest.Server, port string) {

	// Start a test http server to listen
	ts := httptest.NewServer(http.HandlerFunc(handler))
	//defer ts.Close()

	log.Println("Startng a test http server to recieve hits at : ", ts.URL)

	// Get port number
	re := regexp.MustCompile(`.*:(.*)$`)
	match := re.FindStringSubmatch(ts.URL)
	port = match[1]

	return ts, port
}

func formEncode(m map[string]string) string {
	form := url.Values{}
	for k, v := range m {
		form.Add(k, v)
	}
	return form.Encode()
}

const (
	PathPrefix  = "gosiege/"
	SessionPath = PathPrefix + "sessions/"
	NodePath    = PathPrefix + "nodes/"

	MethodNewSession    = "PUT"
	MethodStopSession   = "PUT"
	MethodUpdateSession = "PUT"
)

type NewSessionReq struct {
	Url        string
	Target     string `schema:"target"`
	Port       string `schema:"port"`
	Concurrent string `schema:"concurrent"`
	Delay      string `schema:"delay"`
}

func (r *NewSessionReq) Send() (resp *http.Response, err error) {
	m := map[string]string{
		"concurrent": r.Concurrent,
		"target":     r.Target,
		"port":       r.Port,
		"delay":      r.Delay,
	}

	form := formEncode(m)

	req, err := http.NewRequest("PUT", "http://localhost:8090/gosiege/sessions/new", strings.NewReader(form))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	client := &http.Client{}

	return client.Do(req)
}

type StopSessionReq struct {
	SessionId string
}

func (r StopSessionReq) Send() (resp *http.Response, err error) {
	url := "http://localhost:8090" + "/" + SessionPath + "stop/" + r.SessionId

	log.Println("Sending request : ", url)
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	return client.Do(req)
}

type UpdateSessionReq struct {
	Url        string
	SessionId  string `schema:"Id"`
	Target     string `schema:"target"`
	Port       int    `schema:"port"`
	Concurrent int    `schema:"concurrent"`
	Delay      string `schema:"delay"`
}

type EndSessionReq struct {
	SessionId string
}
