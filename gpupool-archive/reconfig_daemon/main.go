package main

import (
    "bytes"
	"context"
	"fmt"
    "log"
    "strings"
    "time"
    "net/http"
    "io/ioutil"
    "encoding/json"
    // "gopkg.in/yaml.v3"
    "strconv"
    "sort"
    "sync"
    
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/remotecommand"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/apimachinery/pkg/util/sets"
    "k8s.io/apimachinery/pkg/types"
)

const (
    ReconfigDaemonConfigPath 	string = "/etc/kubernetes/reconfig-daemon-config.yaml"
    UseFalcon                   string = "use_falcon"
    // PermitReconfig              string = "permit_reconfig"
    gpuPoolScheduler            string = "my-scheduler"
)

type PodInfo struct {
    name        string
    namespace   string
    gids       []string // Falcon devices
}

type  DeviceLocation struct {
	dev_ID		string
	host_port	string
}

type ReconfigDaemon struct {
    deviceAlloc      map[string]string
    http_cli         *http.Client
    config           *rest.Config
    clientset        *kubernetes.Clientset
    schedulePods     sets.Set[types.UID]
    ignorePods      sets.Set[types.UID]
    schedulePodInfo  map[types.UID]PodInfo
    apiEndpoint      string
    hostPort         map[string]string
    port_convert     map[string]string
}

type Response struct {
	Action string
	Result []bool
	Remark string
}

func newReconfigDaemon() *ReconfigDaemon {
    // read reconfig-daemon-config.yaml
	// buf, err := ioutil.ReadFile(ReconfigDaemonConfigPath)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// var config map[string]string
	// err = yaml.Unmarshal(buf, &config)

    d := ReconfigDaemon{}
    d.deviceAlloc = make(map[string]string) // Dev_GID -> Host_Port
    d.http_cli = &http.Client{}
    d.schedulePods = sets.Set[types.UID]{} // ?
    d.ignorePods = sets.Set[types.UID]{}
    d.schedulePodInfo = make(map[types.UID]PodInfo) // ?
    d.hostPort = make(map[string]string) // nodeName -> Host_Port
    d.port_convert = make(map[string]string)

    d.hostPort["css-host-128"] = "48"
    d.hostPort["css-host-129"] = "0"
    
    d.port_convert["48"] = "30h"
    d.port_convert["0"] = "00h"

    // d.apiEndpoint = config["api_endpoint"]

    var err error
    d.config, err = rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }

    d.clientset, err = kubernetes.NewForConfig(d.config)
    if err != nil {
        panic(err.Error())
    }

    return &d
}

func (d *ReconfigDaemon) get_resource_api(val int, endpoint string, ch chan DeviceLocation, ch_err chan error, wg *sync.WaitGroup) {
	defer wg.Done()
    method := "GET"
	client := &http.Client {}

	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		fmt.Println(err)
        ch_err <- err
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
        ch_err <- err
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
        ch_err <- err
		return
	}
	
	// parsing
	var result []map[string]string
	if err := json.Unmarshal([]byte(string(body)), &result); err != nil {  // Parse []byte to the go struct pointer
		fmt.Printf("Can not unmarshal JSON")
	}
	for _, res := range result {
		dev_type, ok := res["dev_type"]
		if ok && dev_type == "GPU" {
            dl := DeviceLocation{
                dev_ID: strconv.Itoa(val) + "-" + res["Dev_GID"],
                host_port: res["Host_Port#"],
            }
            ch <- dl
		}  
	}
    ch_err <- nil
	return
}

func (d *ReconfigDaemon) updateDevice() error {
    ok := false
    for !ok {
        result := make(chan DeviceLocation, 8)
        err := make(chan error, 2)

        wg := sync.WaitGroup{}
        endpoints := [2]string{"http://9.2.130.147:80/h3api/v1/resources", "http://9.2.130.148:80/h3api/v1/resources"}
        for i, endpoint := range endpoints {
            wg.Add(1)
            go d.get_resource_api(i, endpoint, result, err, &wg)
        }
        wg.Wait()
        close(result)
        close(err)

        ok = true
        for e := range err {
            if e != nil {
                ok = false
            }
        }

        if ok {
            for dl := range result {
                d.deviceAlloc[dl.dev_ID] = dl.host_port
            }
            fmt.Printf("dev: %s\n", d.deviceAlloc)
        }
    }

	return nil
}

func (d *ReconfigDaemon) check_assign(gep string, port_gid string, dev_gid string, endpoint string) bool {
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
        fmt.Printf("Can not unmarshal JSON")
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

func (d *ReconfigDaemon) assign(gep string, port_gid string, dev_gid string, endpoint string) (bool, error) {
    /*
    POST http://localhost:8080/h3api/v1/allocation
    {"gep": "0", "port_gid": "30h", "dev_gid": "001300"}
    */
    fmt.Printf("assign -> gep: %v, port_gid: %v, dev_gid: %v, endpoint: %v\n", gep, port_gid, dev_gid, endpoint)
    
    method := "POST"
    param := fmt.Sprintf(`{"gep" : "%s", "port_gid" : "%s", "dev_gid" : "%s"}`, gep, port_gid, dev_gid)
    payload := strings.NewReader(param)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        fmt.Printf("new request err: %v\n", err)  
        return false, err
    }
    req.Header.Add("Content-Type", "text/plain")

    res, err := d.http_cli.Do(req)
    if err != nil {
        fmt.Printf("http err: %v\n", err)  
        return false, err
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("read body err: %v\n", err)  
        return false, err
    }

    var resp Response
	if err := json.Unmarshal([]byte(string(body)), &resp); err != nil {  // Parse []byte to the go struct pointer
		fmt.Printf("unmarshal err: %v\n", err)  
        return false, err
	}
    // var ok bool = resp.Result[0]
    // log.Printf("assign result: %v", ok)
    // // return resp.Result[0]

    // // double check becasue the api is weird
    // if !ok {
    //     ok = d.check_assign(gep, port_gid, dev_gid, endpoint)
    //     log.Printf("assign double check result: %v", ok)
    // }
    // return ok, nil

    return true, nil
}

func (d *ReconfigDaemon) unassign(gep string, dev_gid string, endpoint string) (bool, error) {
    /*
    DELETE http://localhost:8080/h3api/v1/allocation
    {"gep": "0", "dev_gid": "001300"}
    */
    fmt.Printf("unassign -> gep: %v, dev_gid: %v, endpoint: %v\n", gep, dev_gid, endpoint)
    method := "DELETE"
    param := fmt.Sprintf(`{"gep" : "%s", "dev_gid" : "%s"}`, gep, dev_gid)
    payload := strings.NewReader(param)

    req, err := http.NewRequest(method, endpoint, payload)
    if err != nil {
        fmt.Printf("new request err: %v\n", err)  
        return false, err
    }
    req.Header.Add("Content-Type", "text/plain")

    res, err := d.http_cli.Do(req)
    if err != nil {
        fmt.Printf("http err: %v\n", err)  
        return false, err
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Printf("read body err: %v\n", err)  
        return false, err
    }

    var resp Response
	if err := json.Unmarshal([]byte(string(body)), &resp); err != nil {  // Parse []byte to the go struct pointer
		fmt.Printf("unmarshal err: %v\n", err)  
        return false, err
	}

    // log.Printf("unassign result: %v", resp.Result[0])
    return resp.Result[0], nil
}

func (d *ReconfigDaemon) getGID(name string, namespace string) []string {
    restClient := d.clientset.CoreV1().RESTClient()
    cmd := []string{
        "sh",
        "-c",
        "echo $FALCON_DEVICES",
    }

    req := restClient.Post().Resource("pods").Name(name).Namespace(namespace).SubResource("exec")

    option := &v1.PodExecOptions{
        Command: cmd,
        Stdin:   false,
        Stdout:  true,
        Stderr:  false,
        TTY:     false,
    }

    req.VersionedParams(
        option,
        scheme.ParameterCodec,
    )
    
    exec, err := remotecommand.NewSPDYExecutor(d.config, "POST", req.URL())
    if err != nil {
        fmt.Printf("remotecommand error: %v\n", err);
        return []string{"-1"}
    }

    var stdout bytes.Buffer
    err = exec.Stream(remotecommand.StreamOptions{
        Stdin: nil, // os.Stdin,
        Stdout: &stdout, // os.Stdout,
        Stderr: nil, // stderr,
    })
    if err != nil {
        fmt.Printf("Stream error: %v\n", err);
        return []string{"-1"}
    }

    env_var := stdout.String() 
    if env_var == "" {
        // no gpu
        return []string{"-1"}
    }
    env_var = env_var[:len(env_var)-1] // remove \n
    gids := strings.Split(env_var, ",")
    return gids
}

func (d *ReconfigDaemon) reconfig(nodeName string, demand int) bool {
    fmt.Printf("dst: %v, require: %v\n", nodeName, demand)

    update_device_startTime := time.Now()
    d.updateDevice()
    update_device_endTime := time.Now()
    update_device_duration := update_device_endTime.Sub(update_device_startTime)
    log.Printf("update device duration: %v\n", update_device_duration)

    used_gpu := sets.Set[string]{}
    for po, _ := range d.schedulePods {
        for _, id := range d.schedulePodInfo[po].gids {
            used_gpu.Insert(id)
        }
    }
    
    var optionGPUs []struct {
		dev_gid     string
		host_port   string
        score       int
	}

    for dev, node := range d.deviceAlloc {
        if !used_gpu.Has(dev) && node != d.hostPort[nodeName]{
            optionGPUs = append(optionGPUs, struct {
                dev_gid    string
                host_port  string
                score      int
            }{dev, node, 0})
        }
    }

    // get score
    fmt.Printf("usedGPUs: %v, optionGPUs: %v\n", used_gpu, optionGPUs)
    gpuCounts := make(map[string]int)
    for _, gpuOption := range optionGPUs {
        gpuCounts[gpuOption.host_port]++
    }
    for i := range optionGPUs {
        optionGPUs[i].score = gpuCounts[optionGPUs[i].host_port]
    }

    sort.Slice(optionGPUs, func(i, j int) bool {
		return optionGPUs[i].score < optionGPUs[j].score
	})

    // unassign first and then assign
    for _, dev := range optionGPUs {
        if demand <= 0 {
            break
        }

        var endpoint string
        if string(dev.dev_gid[0]) == "0" {
            endpoint = "http://9.2.130.147:80/h3api/v1/allocation"
        } else {
            endpoint = "http://9.2.130.148:80/h3api/v1/allocation"
        }
        dev.dev_gid = dev.dev_gid[2:]
        dst_node := d.port_convert[d.hostPort[nodeName]]

        unassign_startTime := time.Now()
        ok, err := d.unassign("0", dev.dev_gid, endpoint)
        unassign_endTime := time.Now()
        unassign_duration := unassign_endTime.Sub(unassign_startTime)
        log.Printf("unassign duration: %v\n", unassign_duration)
        
        if err != nil {
            continue
        }

        assign_startTime := time.Now()
        ok, _ = d.assign("0", dst_node, dev.dev_gid, endpoint)
        assign_endTime := time.Now()
        assign_duration := assign_endTime.Sub(assign_startTime)
        log.Printf("assign duration: %v\n", assign_duration)
        
        if ok {
            demand--
        }

        fmt.Printf("current demand: %v\n", demand)
    }
    return demand == 0
}

func (d *ReconfigDaemon) podUseFalcon(name string, namespace string) bool {
    po, err := d.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        fmt.Printf("get pod err: %v\n", err)  
    }  
    return po.ObjectMeta.Annotations[UseFalcon] == "true"
}

func (d *ReconfigDaemon) podIsScheduled(name string, namespace string) bool {
    po, err := d.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        fmt.Printf("get pod err: %v\n", err)  
    }  
    for _, cond := range po.Status.Conditions {
        if cond.Type == "PodScheduled" {
            // fmt.Printf("pod %v is scheduled\n", name)  
            return cond.Status == "True"
        }
    }
    return false
}

func (d *ReconfigDaemon) getPodAnnotation(name string, namespace string, annotation string) string {
    po, err := d.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        fmt.Printf("get pod err: %v\n", err)  
    }  
    return po.ObjectMeta.Annotations[annotation]
}

// func (d *ReconfigDaemon) podPermitReconfig(name string, namespace string) bool {
//     po, err := d.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
//     if err != nil {
//         fmt.Printf("get pod err: %v\n", err)  
//     }  
//     return po.ObjectMeta.Annotations[PermitReconfig] == "true"
// }

func (d *ReconfigDaemon) watchReconfigEvent(eventChan <-chan watch.Event) error {
    for {
        select {
        case event := <- eventChan:
            ev, _ := event.Object.(*v1.Event)
            name := ev.InvolvedObject.Name
            namespace := ev.InvolvedObject.Namespace
            
            cur_pod, err := d.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
            
            // skip if the pod does not exist
            if err != nil {
                // fmt.Printf("skip because of err: %v\n", err)  
                continue
            }  
            
            // if already scheduled then skip
            if cur_pod.Status.Phase != "Pending" {
                // fmt.Printf("already scheduled or done \n")  
                continue
            }
            
            log.Printf("reconfig event: %s, %s, %s\n", ev.InvolvedObject.Kind, ev.InvolvedObject.Name, ev.Reason)

            // get to know which gpu is used
            update_GPU_startTime := time.Now()
            ready := false
            d.schedulePods = sets.Set[types.UID]{}
            for !ready {
                ready = true
                allPods, _ := d.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
                for _, po := range allPods.Items {
                    uid := po.ObjectMeta.UID
                    // only keep other pods that use GPU attached on Falcon and are not scheduled yet
                    if d.ignorePods.Has(uid) {
                        continue
                    }
                    // fmt.Printf("check pod: %v\n", po.ObjectMeta.Name)
                    if !d.podIsScheduled(po.ObjectMeta.Name, po.ObjectMeta.Namespace) {
                        // should not insert!!!
                        continue
                    }
                    if !d.podUseFalcon(po.ObjectMeta.Name, po.ObjectMeta.Namespace) {
                        d.ignorePods.Insert(uid)
                        continue
                    }
                    if (po.ObjectMeta.Name == name) && (po.ObjectMeta.Namespace == namespace){
                        // should not insert!!!
                        continue
                    }
                    if (po.Status.Phase == "Succeeded") || (po.Status.Phase == "Failed"){
                        d.ignorePods.Insert(uid)
                        continue
                    }
                    d.schedulePods.Insert(uid)

                    if po.Status.Phase == "Running" && len(d.schedulePodInfo[uid].gids) ==  0 {
                        // the pod is running but the gpu it used haven't been recorded
                        // runningPods.Insert(uid)
                        var info PodInfo
                        info.name = po.ObjectMeta.Name
                        info.namespace = po.ObjectMeta.Namespace
                        info.gids = d.getGID(info.name, info.namespace)
                        d.schedulePodInfo[uid] = info
                    } else if po.Status.Phase == "Pending" {
                        ready = false
                    } 
                }
            }
            
            update_gpu_endTime := time.Now()
	        update_gpu_duration := update_gpu_endTime.Sub(update_GPU_startTime)
            log.Printf("update GPU availability: %v\n", update_gpu_duration)

            nodeName := d.getPodAnnotation(name, namespace, "dst_node")
            demand := d.getPodAnnotation(name, namespace, "gpu_demand")
            demand_cnt, _ := strconv.Atoi(demand)
            _ = d.reconfig(nodeName, demand_cnt)

            reconfig_endTime := time.Now()
	        reconfig_duration := reconfig_endTime.Sub(update_GPU_startTime)
            log.Printf("reconfig total time: %v\n", reconfig_duration)
        }
    }
    return nil
}

func (d *ReconfigDaemon) garbageCollect() error {
    for {
        for uid, _ := range d.schedulePodInfo {
            if !d.schedulePods.Has(uid) {
                delete(d.schedulePodInfo, uid)
            }
        }
        time.Sleep(600 * time.Second)
    }
    return nil
}

func main() {
    d := newReconfigDaemon()
    d.updateDevice()  

    opts := metav1.ListOptions{
       FieldSelector: "involvedObject.kind=Pod",
       Watch:         true,
    }
    watcher, err := d.clientset.CoreV1().Events("").Watch(context.TODO(), opts)
    if err != nil {
        panic(err.Error())
    }
    
    reconfigChan := make(chan watch.Event)

    go func() {
        err := d.watchReconfigEvent(reconfigChan)
        if err != nil {
			fmt.Printf("watchReconfigEvent err: %v\n", err)
		}
    }()

    // go func() {
    //     err := d.garbageCollect()
    //     if err != nil {
	// 		fmt.Printf("garbageCollect err: %v\n", err)
	// 	}
    // }()

    for {
        select {
        case event := <- watcher.ResultChan():
            ev, _ := event.Object.(*v1.Event)
            if ev.Reason == "Reconfig" {
                reconfigChan <- event
            }
        }
    }
}