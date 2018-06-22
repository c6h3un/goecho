package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
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
	fmt.Fprintf(w, "You are from %s, ", r.RemoteAddr)
	fmt.Fprintf(w, "this is %s ", GetLocalIP())
	fmt.Fprintf(w, "(%s)", hostname)
}

func main() {
	flag.Parse()
	fmt.Println("server starting on port ", *PORT)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+*PORT, nil)
}
