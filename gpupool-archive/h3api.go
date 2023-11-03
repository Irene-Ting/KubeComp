package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
)

type Response struct {
	Action string
	Result []bool
	Remark string
}

func assign(gep string, port_gid string, dev_gid string, endpoint string) bool {
    /*
    POST http://localhost:8080/h3api/v1/allocation
    {"gep": "0", "port_gid": "30h", "dev_gid": "001300"}
    */
    log.Printf("assign:")
    log.Printf("gep: %v, port_gid: %v, dev_gid: %v, endpoint: %v", gep, port_gid, dev_gid, endpoint)
    
    method := "POST"
    param := fmt.Sprintf(`{"gep" : "%s", "port_gid" : "%s", "dev_gid" : "%s"}`, gep, port_gid, dev_gid)
    payload := strings.NewReader(param)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        panic(err.Error())
    }
    req.Header.Add("Content-Type", "text/plain")

    http_cli := &http.Client{}
    res, err := http_cli.Do(req)
    if err != nil {
        panic(err.Error())
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err.Error())
    }

    var resp Response
	if err := json.Unmarshal([]byte(string(body)), &resp); err != nil {  // Parse []byte to the go struct pointer
		panic(err.Error())
	}

    log.Printf("unassign resp body: %v", string(body))
    log.Printf("unassign resp: %v", resp)
    return resp.Result[0]
}

func unassign(gep string, dev_gid string, endpoint string) bool {
    /*
    DELETE http://localhost:8080/h3api/v1/allocation
    {"gep": "0", "dev_gid": "001300"}
    */
    log.Printf("unassign:")
    log.Printf("gep: %v, dev_gid: %v, endpoint: %v", gep, dev_gid, endpoint)

    method := "DELETE"
    param := fmt.Sprintf(`{"gep" : "%s", "dev_gid" : "%s"}`, gep, dev_gid)
    payload := strings.NewReader(param)
    fmt.Println(payload)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        panic(err.Error())
    }
    req.Header.Add("Content-Type", "text/plain")

    http_cli := &http.Client{}
    res, err := http_cli.Do(req)
    if err != nil {
        panic(err.Error())
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err.Error())
    }

    var resp Response
	if err := json.Unmarshal([]byte(string(body)), &resp); err != nil {  // Parse []byte to the go struct pointer
		panic(err.Error())
	}

    log.Printf("unassign resp body: %v", string(body))
    log.Printf("unassign resp: %v", resp)
    return resp.Result[0]
}

func check_assign(gep string, port_gid string, dev_gid string, endpoint string) bool {
	method := "GET"
	client := &http.Client {}
    
    var res_endpoint string
    if endpoint == "http://9.2.130.147:80/h3api/v1/allocation" {
        res_endpoint = "http://9.2.130.147:80/h3api/v1/resources"
    } else {
        res_endpoint = "http://9.2.130.148:80/h3api/v1/resources"
    }

    req, err := http.NewRequest(method, res_endpoint, nil)
    if err != nil {
        fmt.Println(err)
        return false
    }
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return false
    }
	defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return false
    }
        
    // parsing
    var result []map[string]string
    if err := json.Unmarshal([]byte(string(body)), &result); err != nil {  // Parse []byte to the go struct pointer
        log.Printf("Can not unmarshal JSON")
    }

    var host string = ""
    for _, res := range result {
        dev, ok := res["Dev_GID"]
        if ok && dev == dev_gid {
            if dev == dev_gid {
                host = res["Host_Port#"]
                break
            }
        }      
    }

    for _, res := range result {
        port, ok := res["Port#"]
        if ok && port == host {
            return res["Port_GID"] == port_gid
        }      
    }
    return false
}

func main() {	
	var ok bool
    // ok = unassign("0", "001300", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)
	// ok = unassign("0", "001800", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = unassign("0", "000C00", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = unassign("0", "000500", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)

	// ok = assign("0", "30h", "001300", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = assign("0", "30h", "001800", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = assign("0", "00h", "000C00", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = assign("0", "00h", "000500", "http://9.2.130.147:80/h3api/v1/allocation")
	// fmt.Println(ok)

    // ok = unassign("0", "002F00", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)
	// ok = unassign("0", "004200", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = unassign("0", "001A00", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = unassign("0", "000500", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)

    // ok = assign("0", "30h", "002F00", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = assign("0", "30h", "004200", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = assign("0", "30h", "001A00", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)
    // ok = assign("0", "30h", "000500", "http://9.2.130.148:80/h3api/v1/allocation")
	// fmt.Println(ok)

    ok = check_assign("0", "00h", "001300", "http://9.2.130.147:80/h3api/v1/allocation") 
    fmt.Println(ok)
}