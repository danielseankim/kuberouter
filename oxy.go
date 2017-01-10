package main

import (
    "net/http"
    "fmt"
    "os"
    "strings"
    "regexp"
    "net/http/httputil"
    "strconv"
)


func getAddrs() (map[string]string) {
    services := ParseRawServices(getEnviro())
    return services
}

func getEnviro() []string {
    return os.Environ()
}
//ParseRawServices returns a map of the services we need to proxy
// expecs environment variabls like: OUTSCORE_DEPLOYMENT_PORT_8080_TCP_ADDR=10.95.249.177
func ParseRawServices(env []string) map[string]string {
    targets := map[string]string{}
    for _, val := range env {
        if !strings.Contains(val, "ADDR") {
            continue
        }
        if (strings.Contains(val, "ROUTER") || strings.Contains(val, "KUBERNETES")) {
            continue
        }
        fmt.Println("%s", val)
        portre := regexp.MustCompile("_(\\d*)_")
        port := strings.Replace(portre.FindString(val), "_", "", 2)
        ipre := regexp.MustCompile("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}")
        portI, _ := strconv.Atoi(port)
        if (portI > 65535) {
            continue
        }
        targets[port] = ipre.FindString(val)
    }
    return targets
}

func PortParse(url string) string {
     return strings.Split(url, ":")[1]
}

func NewMultipleHostReverseProxy(targets map[string]string) *httputil.ReverseProxy {
    director := func(req *http.Request) {
            requestPort := PortParse(req.Host)
            newHost := fmt.Sprintf("%s:%s", targets[requestPort], requestPort)
            req.URL.Scheme = "http"
            req.URL.Host = newHost
        }
        return &httputil.ReverseProxy{Director: director}
}

func main() {
    targets := getAddrs()
    fmt.Printf("Targets found: %q\n", targets)
    proxy := NewMultipleHostReverseProxy(targets)
    targetSize := len(targets)
    targetID := 1
    for port, _ := range targets {
        serve := fmt.Sprintf(":%s", port)
        if (targetID == targetSize) {
            http.ListenAndServe(serve, proxy)
        } else {
            go func() {
                http.ListenAndServe(serve, proxy)
            }()
        }
        targetID = targetID + 1
    }
}
