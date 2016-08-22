package main

import (
    "testing"
    "github.com/vulcand/oxy/testutils"
)

func TestRawServiceParse(t *testing.T) {
    fail := false
    wants := map[string]string{}
    wants["8080"] = "10.95.249.177"
    wants["6061"] = "10.95.255.21"
    raw := []string{
        "OUTSCORE_DEPLOYMENT_PORT_8080_TCP_ADDR=10.95.249.177",
        "TUNNELTEST_DEPLOYMENT_PORT_6061_TCP_ADDR=10.95.255.21",
    }
    returned := ParseRawServices(raw)

    for key, val := range returned {
        if wants[key] != val {
            fail = true
        }
    }
    if fail {
            t.Errorf("ParseRawService failure == %q, should be: %q", returned, wants)
    }
}

func TestPortParse(t *testing.T) {
    url := testutils.ParseURI("http://localhost:63450")
    want := "63450"
    got := PortParse(url.Host)
    if (want != got) {
        t.Errorf("PortParse failure == %q, should be: %q", got, want)
    }
}
