package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"html/template"
)

// listen port
var (
	PORT = flag.String("p", "8888", "service port")
	ENABLE_TLS = flag.Bool("tls.enable", false, "enable tls")
)

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "You are from %s, ", r.RemoteAddr)
	fmt.Fprintf(w, "This is %s ", GetLocalIP())
	fmt.Fprintf(w, "(%s)\n", hostname)
}

func main() {
	flag.Parse()
	fmt.Println("server starting on port ", *PORT)
	http.HandleFunc("/", templatePage)
	http.HandleFunc("/info", handler)
	http.HandleFunc("/dumpPacket", dumpPacket)
	http.HandleFunc("/waitSeconds", waitSeconds)
	http.HandleFunc("/health", ok)
	http.HandleFunc("/ready", ok)
	http.HandleFunc("/echo/", echo)
	if *ENABLE_TLS {
		http.ListenAndServeTLS(":"+*PORT, "tls/server.crt", "tls/server.key", nil)
	}else{
		http.ListenAndServe(":"+*PORT, nil)
	}
}

func dumpPacket(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {

		log.Fatal(err)
	}
	fmt.Fprintf(w, "Hello, this is %s\n\n%s %s %s\nHost: %s\nBody:\n%s\nHeaders:\n", hostname, r.Method, r.URL, r.Proto, r.Host, string(body))
	for k, v := range r.Header {
		fmt.Fprintf(w, "\t%v: %v\n", k, v)
	}
}

func waitSeconds(w http.ResponseWriter, r *http.Request) {
	time.Sleep(26 * time.Second)
	fmt.Fprintf(w, "200 OK")
}

func ok(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, strings.TrimPrefix(r.RequestURI, "/echo/"))
}

func templatePage(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	background := os.Getenv("BG_COLOR") // https://www.w3schools.com/colors/colors_names.asp
	if len(background) == 0 {
		background = "LightGray"
	}
	tpl := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body style="background-color:`+background+`;">
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(tpl)
	check(err)
	appName := os.Getenv("APP_NAME")
	appMessage := os.Getenv("APP_MESSAGE")
	data := struct {
		Title string
		Items []string
	}{
		Title: appName,
		Items: []string{
			hostname,
			GetLocalIP(),
			appMessage,
		},
	}

	err = t.Execute(w, data)
	check(err)

	// noItems := struct {
	// 	Title string
	// 	Items []string
	// }{
	// 	Title: "My another page",
	// 	Items: []string{},
	// }

	// err = t.Execute(os.Stdout, noItems)
	// check(err)
}