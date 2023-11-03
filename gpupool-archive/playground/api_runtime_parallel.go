package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
)

// var http_cli *http.Client
var cnt int

func unassign(gep string, dev_gid string, endpoint string) {
    http_cli := &http.Client{}
    method := "DELETE"
    param := fmt.Sprintf(`{"gep" : "%s", "dev_gid" : "%s"}`, gep, dev_gid)
    payload := strings.NewReader(param)
    fmt.Println(payload)

    req, err := http.NewRequest(method, endpoint, payload)    
    if err != nil {
        fmt.Printf("new request err: %s\n", err)  
    }
    req.Header.Add("Content-Type", "text/plain")

    maxRetries := 10
    var res *http.Response
    for retry := 0; retry < maxRetries; retry++ {
        res, err = http_cli.Do(req)
        if err != nil {
            fmt.Println("http err: %s\n", err)  
            req, err = http.NewRequest(method, endpoint, payload) 
            req.Header.Add("Content-Type", "text/plain")
        } else {
            break
        }
    }

    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("read body err: %s\n", err)  
    }
    fmt.Printf("body: %s\n", string(body)) 
}

func assign(gep string, port_gid string, dev_gid string, endpoint string) {
    http_cli := &http.Client{}
    method := "POST"
    param := fmt.Sprintf(`{"gep" : "%s", "port_gid" : "%s", "dev_gid" : "%s"}`, gep, port_gid, dev_gid)
    payload := strings.NewReader(param)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        fmt.Printf("new request err: %s\n", err)  
    }
    req.Header.Add("Content-Type", "text/plain")

    maxRetries := 10
    var res *http.Response
    for retry := 0; retry < maxRetries; retry++ {
        res, err = http_cli.Do(req)
        if err != nil {
            fmt.Println("http err: %s\n", err)  
            req, err = http.NewRequest(method, endpoint, payload)
            req.Header.Add("Content-Type", "text/plain")
        } else {
            break
        }
    }

    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("read body err: %s\n", err)  
    }

    fmt.Printf("body: %s\n", string(body)) 
}

func yourFunction() {
	method := "GET"
	client := &http.Client {
	}
	endpoint := "http://9.2.130.148:80/h3api/v1/resources"
	// "http://9.2.130.147:80/h3api/v1/resources"

	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		fmt.Println(err)
	}
    
    maxRetries := 10
    var res *http.Response
    for retry := 0; retry < maxRetries; retry++ {
        res, err = client.Do(req)
        if err != nil {
            fmt.Println("http err: %s\n", err)  
            req, err = http.NewRequest(method, endpoint, nil)
        } else {
            break
        }
    }

	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func api_measure(port_gid string, dev_gid string) {
    unassign_total := time.Duration(0)
    assign_total := time.Duration(0)
    for i := 1; i <= cnt; i++ {
        fmt.Printf("%s starts\n", dev_gid)
        startTime := time.Now()
        unassign("0", dev_gid, "http://9.2.130.148:80/h3api/v1/allocation")
        endTime := time.Now()
        duration := endTime.Sub(startTime)
        unassign_total += duration

        startTime = time.Now()
        assign("0", port_gid, dev_gid, "http://9.2.130.148:80/h3api/v1/allocation")
        endTime = time.Now()
        duration = endTime.Sub(startTime)
        assign_total += duration
    }
    fmt.Printf("%s unassign in average: %s\n", dev_gid, unassign_total.Seconds()/float64(cnt))
    fmt.Printf("%s assign in average: %s\n", dev_gid, assign_total.Seconds()/float64(cnt))
}

func main() {
	cnt = 1
	
    go func() {
        api_measure("30h", "001300")
	}()

    go func() {
        api_measure("30h", "001800")
	}()

    go func() {
        api_measure("00h", "000C00")
	}()

    go func() {
        api_measure("00h", "000500")
	}()
	
    time.Sleep(90 * time.Second)
}