package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
)

var http_cli *http.Client



func unassign(gep string, dev_gid string, endpoint string) {
    /*
    DELETE http://localhost:8080/h3api/v1/allocation
    {"gep": "0", "dev_gid": "001300"}
    */
    // log.Printf("unassign -> gep: %v, dev_gid: %v, endpoint: %v", gep, dev_gid, endpoint)

    method := "DELETE"
    param := fmt.Sprintf(`{"gep" : "%s", "dev_gid" : "%s"}`, gep, dev_gid)
    payload := strings.NewReader(param)
    fmt.Println(payload)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        fmt.Printf("new request err: %v\n", err)  
    }
    req.Header.Add("Content-Type", "text/plain")

    res, err := http_cli.Do(req)
    if err != nil {
        fmt.Printf("http err: %v\n", err)  
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("read body err: %v\n", err)  
    }
    fmt.Printf("body: %v\n", string(body)) 
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
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func assign(gep string, port_gid string, dev_gid string, endpoint string) {
    /*
    POST http://localhost:8080/h3api/v1/allocation
    {"gep": "0", "port_gid": "30h", "dev_gid": "001300"}
    */
    // log.Printf("assign -> gep: %v, port_gid: %v, dev_gid: %v, endpoint: %v", gep, port_gid, dev_gid, endpoint)
    
    method := "POST"
    param := fmt.Sprintf(`{"gep" : "%s", "port_gid" : "%s", "dev_gid" : "%s"}`, gep, port_gid, dev_gid)
    payload := strings.NewReader(param)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        fmt.Printf("new request err: %v\n", err)  
    }
    req.Header.Add("Content-Type", "text/plain")

    res, err := http_cli.Do(req)
    if err != nil {
        fmt.Printf("http err: %v\n", err)  
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("read body err: %v\n", err)  
    }

    fmt.Printf("body: %v\n", string(body)) 
}

func main() {
	http_cli = &http.Client{}
	cnt := 20
	unassign_total := time.Duration(0)
	assign_total := time.Duration(0)
	for i := 1; i <= cnt; i++ {
		startTime := time.Now()
		unassign("0", "001800", "http://9.2.130.148:80/h3api/v1/allocation")
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		unassign_total += duration

		startTime = time.Now()
		assign("0", "30h", "001800", "http://9.2.130.148:80/h3api/v1/allocation")
		endTime = time.Now()
		duration = endTime.Sub(startTime)
		assign_total += duration
	}
	fmt.Printf("unassign in average: %v\n", unassign_total.Seconds()/float64(cnt))
	fmt.Printf("assign in average: %v\n", assign_total.Seconds()/float64(cnt))
}