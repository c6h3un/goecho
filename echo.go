package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// listen port
var (
	PORT = flag.String("p", "8888", "service port")
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
	//fmt.Fprintf(w, "You are from %s, ", r.RemoteAddr)
	fmt.Fprintf(w, "This is %s ", GetLocalIP())
	fmt.Fprintf(w, "(%s)\n", hostname)
}

func main() {
	flag.Parse()
	fmt.Println("server starting on port ", *PORT)

	http.HandleFunc("/info", handler)
	http.HandleFunc("/dumpPacket", dumpPacket)
	http.HandleFunc("/waitSeconds", waitSeconds)
	http.ListenAndServe(":"+*PORT, nil)
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
