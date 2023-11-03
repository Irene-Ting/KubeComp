package main

import (
	"fmt"
  	"net/http"
  	"io/ioutil"
  	"encoding/json"
    "strconv"
	"sync"
	"time"
)

// FalconServer is a device plugin server
type FalconServer struct {
	localIP		string
	hostPort	string
	gpuLookUp 	map[string]string
	apiEndpoints [2]string
	// apiEndpoints [1]string
}

type  DevicePair struct {
	dev_ID		string
	GPU_GUID	string
}

func NewFalconServer() *FalconServer {
	return &FalconServer{
		localIP:	 	"9.2.131.129",
		hostPort:		"0",
		gpuLookUp: 		make(map[string]string),
		apiEndpoints: 	[2]string{"http://9.2.130.147:80/h3api/v1/resources", "http://9.2.130.148:80/h3api/v1/resources"},
		// apiEndpoints: 	[1]string{"http://9.2.130.147:80/h3api/v1/resources"},
	}
}

func (s *FalconServer) get_resource_api(val int, endpoint string, ch chan DevicePair, wg *sync.WaitGroup) error{
	method := "GET"
	client := &http.Client {}

	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	// parsing
	var result []map[string]string
	if err := json.Unmarshal([]byte(string(body)), &result); err != nil {  // Parse []byte to the go struct pointer
		fmt.Printf("Can not unmarshal JSON")
	}
	for _, res := range result {
		fmt.Println(res["dev_type"])
		dev_type, ok := res["dev_type"]
		if ok && dev_type == "GPU" {
			if res["Host_Port#"] == s.hostPort {
				dp := DevicePair{
					dev_ID: strconv.Itoa(val) + "-" + res["Dev_GID"],
					GPU_GUID: res["OB_GPU_GUID"],
				}
				ch <- dp
			}
		}  
	}
	wg.Done()
	return nil
}

// discover device
func (s *FalconServer) listDevice() error {
	result := make(chan DevicePair, 8)
	wg := sync.WaitGroup{}

	for i, endpoint := range s.apiEndpoints {
		wg.Add(1)
		go s.get_resource_api(i, endpoint, result, &wg)
	}
	wg.Wait()
	close(result)

	for dp := range result {
		s.gpuLookUp[dp.dev_ID] = dp.GPU_GUID
	}
		
	return nil
}


func main() {
	s := NewFalconServer()
	startTime := time.Now()
	s.listDevice()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("duration: ", duration)
}
