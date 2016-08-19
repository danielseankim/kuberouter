package main

import (
    "net/http"
    "net/url"
    "fmt"
    "os/exec"
    "bufio"
    "strings"
    "bytes"
    "regexp"
    "github.com/vulcand/oxy/forward"
)


func getAddrs() (map[string]string) {
    services := ParseRawServices(GetRawEnv())
    fmt.Println(services)
    return services
}

func GetRawEnv() []byte {
    output, _ := exec.Command("env").Output()
    return output
}

func ParseRawServices(raw []byte) map[string]string {
    serviceMap := map[string]string{}
    scanner := bufio.NewScanner(bytes.NewReader(raw))

    scanner.Split(bufio.ScanLines)

    for scanner.Scan() {
        svc := scanner.Text()
        if !strings.Contains(svc, "ADDR") {
            continue
        }
        portre := regexp.MustCompile("_(\\d*)_")
        port := strings.Replace(portre.FindString(svc), "_", "", 2)
        ipre := regexp.MustCompile("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}")
        serviceMap[port] = ipre.FindString(svc)
    }
    return serviceMap
}

func PortParse(url string) string {
     return strings.Split(url, ":")[1]
}

func main() {
    fwd, _ := forward.New()

    servers := []http.Server{}
    targets := getAddrs()
    fmt.Printf("Targets found: %q", targets)
    for port, host := range targets {

        redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            req.URL = &url.URL{
                Scheme: "http",
                Host: fmt.Sprintf("%s:%s", host, port),
            }
            fwd.ServeHTTP(w, req)
        })
        servers = append(servers, http.Server{
            Addr: fmt.Sprintf(":%s", port),
            Handler: redirect,

        })
    }

    for _, server := range servers {
        server.ListenAndServe()
    }
}
