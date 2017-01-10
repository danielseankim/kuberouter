package main

import (
    "net/http"
    "net/url"
    "fmt"
    "os"
    "strings"
    "regexp"
    "github.com/vulcand/oxy/forward"
)

type target struct {
    Name string
    Endpoint string
    Port string
}


func getAddrs() ([]target) {
    services := ParseRawServices(getEnviro())
    return services
}

func getEnviro() []string {
    return os.Environ()
}
//ParseRawServices returns a map of the services we need to proxy
// expecs environment variabls like: OUTSCORE_DEPLOYMENT_PORT_8080_TCP_ADDR=10.95.249.177
func ParseRawServices(env []string) []target {
    serviceMap := []target{}

    for _, val := range env {
        target := target{}
        if !strings.Contains(val, "ADDR") {
            continue
        }
        target.Name = strings.Split(val, "_")[0]
        portre := regexp.MustCompile("_(\\d*)_")
        port := strings.Replace(portre.FindString(val), "_", "", 2)
        ipre := regexp.MustCompile("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}")
        target.Port = port
        target.Endpoint = ipre.FindString(val)
        serviceMap = append(serviceMap, target)
    }
    return serviceMap
}

func PortParse(url string) string {
     return strings.Split(url, ":")[1]
}

func main() {

    targets := getAddrs()
    fmt.Printf("Targets found: %q\n", targets)
    target_size := len(targets)
    target_id := 1
    rs := []http.HandlerFunc{}
    for _, t := range targets {
        fmt.Printf("Serving: %q:%q\n", t.Endpoint, t.Port)

        fwd, _ := forward.New()
        target := fmt.Sprintf("%s:%s", t.Endpoint, t.Port)
        r := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                req.URL = &url.URL{
                    Scheme: "http",
                    Host: target,
                }
                fwd.ServeHTTP(w, req)
        })
        fmt.Printf("%q", r)
        rs = append(rs, r)
        if target_id != target_size {
            // go http.ListenAndServe(fmt.Sprintf(":%s", t.Port), &r)
        } else {
            // http.ListenAndServe(fmt.Sprintf(":%s", t.Port), &r)
        }
        target_id = target_id + 1
	}
    fmt.Printf("%q", rs)
}
