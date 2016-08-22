package main

import (
    "net/http"
    "net/url"
    "fmt"
    "os"
    "strings"
    "regexp"
    "sync"
    "github.com/vulcand/oxy/forward"
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
    serviceMap := map[string]string{}

    for _, val := range env {
        if !strings.Contains(val, "ADDR") {
            continue
        }
        portre := regexp.MustCompile("_(\\d*)_")
        port := strings.Replace(portre.FindString(val), "_", "", 2)
        ipre := regexp.MustCompile("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}")
        serviceMap[port] = ipre.FindString(val)
    }
    return serviceMap
}

func PortParse(url string) string {
     return strings.Split(url, ":")[1]
}

func main() {

    targets := getAddrs()
    fmt.Printf("Targets found: %q\n", targets)

    wg := &sync.WaitGroup{}

    for port, host := range targets {
        fmt.Printf("Serving: %q:%q\n", host,port)

        fwd, _ := forward.New()
        redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                target := fmt.Sprintf("%s:%s", host, port)
                req.URL = &url.URL{
                    Scheme: "http",
                    Host: target,
                }
                fwd.ServeHTTP(w, req)
        })

        wg.Add(1)
        go func() {
	           http.ListenAndServe(fmt.Sprintf(":%s", port), &redirect)
               wg.Done()
        }()
	}
    wg.Wait()
}
