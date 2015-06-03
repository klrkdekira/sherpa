package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/namsral/flag"
)

var host, upstream string

func init() {
	flag.StringVar(&upstream, "upstream", "", "upstream service url (example: https://popit.mysociety.org)")
	flag.StringVar(&host, "http", "0.0.0.0:8080", "<addr>:<port> to listen on")
	flag.Parse()
}

func main() {
	proxy := negroni.New(negroni.NewLogger(), negroni.NewRecovery())
	proxy.UseHandler(http.HandlerFunc(proxyHandler))
	proxy.Run(host)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	var req *http.Request
	var err error

	if err = r.ParseForm(); err != nil {
		log.Println(err)
	}

	req, err = http.NewRequest(r.Method, upstream+r.RequestURI, r.Body)
	if err != nil {
		log.Println(err)
	}
	req.Header = r.Header
	req.Header.Set("User-Agent", "sinar-api-mirror")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	for k, l := range resp.Header {
		if k == "Set-Cookie" {
			continue
		}

		for _, v := range l {
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("Server", "sinar-api-mirror")
	w.Write(body)
}
