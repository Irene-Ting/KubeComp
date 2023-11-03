package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	buf, err := ioutil.ReadFile("/home/lsalab/gpupool-with-k8s/playground/testfile.yaml")
	if err != nil {
		fmt.Errorf("in file %q: %w", "./testfile.yaml", err)
	}

	var myData map[string]string
	err = yaml.Unmarshal(buf, &myData)
	if err != nil {
		fmt.Errorf("in file %q: %w", "./testfile.yaml", err)
	} else {
		// fmt.Print(myData["api_endpoint"])
		ips := strings.Split(myData["local_ips"], ",")
		geps := strings.Split(myData["geps"], ",")
		for k, v := range(ips) {
			if v == "9.2.131.129" {
				fmt.Print(geps[k])
			}
		}
	}
}
