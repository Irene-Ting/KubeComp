package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
  	"net/http"
  	"io/ioutil"
  	"encoding/json"
	"gopkg.in/yaml.v3"
    "strconv"
	"sync"
	"reflect"
	"sort"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	// pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

const (
	resourceName        	string = "falcon.com/gpu"
	falconSocket          	string = "falcon.sock"
	KubeletSocket 			string = "kubelet.sock"
	DevicePluginPath 		string = "/var/lib/kubelet/device-plugins/"
	DevicePluginConfigPath 	string = "/etc/kubernetes/device-plugin-config.yaml"
	NodeIP 					string = "NODE_IP"
)

// FalconServer is a device plugin server
type FalconServer struct {
	srv         *grpc.Server
	devices     map[string]*pluginapi.Device
	notify      chan bool
	ctx         context.Context
	cancel      context.CancelFunc
	restartFlag bool 
	localIP		string
	hostPort	string
	gpuLookUp 	map[string]string
	// apiEndpoint string
	apiEndpoints [2]string
	deviceCheckInterval time.Duration
}

type  DevicePair struct {
	dev_ID		string
	GPU_GUID	string
}

func NewFalconServer() *FalconServer {
	ctx, cancel := context.WithCancel(context.Background())

	// read device-plugin-config.yaml
	buf, err := ioutil.ReadFile(DevicePluginConfigPath)
	if err != nil {
		panic(err.Error())
	}

	var config map[string]string
	err = yaml.Unmarshal(buf, &config)
	host_port := ""
	if err != nil {
		panic(err.Error())
	} else {
		ips := strings.Split(config["local_ips"], ",")
		host_ports := strings.Split(config["host_ports"], ",")
		log.Infoln("ips: ", ips)
		log.Infoln("host_ports: ", host_ports)
		for k, v := range(ips) {
			if v == os.Getenv(NodeIP) {
				host_port = host_ports[k]
				break
			}
		}
	}

	if host_port == "" {
		log.Printf("Node IP is not in the config file")
	}

	return &FalconServer{
		devices:     	make(map[string]*pluginapi.Device),
		srv:         	grpc.NewServer(grpc.EmptyServerOption{}),
		notify:      	make(chan bool),
		ctx:         	ctx,
		cancel:      	cancel,
		restartFlag: 	false,
		localIP:	 	os.Getenv(NodeIP),
		hostPort:		host_port,
		gpuLookUp: 		make(map[string]string),
		// apiEndpoint: 	config["api_endpoint"],
		apiEndpoints: 	[2]string{"http://9.2.130.147:80/h3api/v1/resources", "http://9.2.130.148:80/h3api/v1/resources"},
		deviceCheckInterval: 20 * time.Second,
	}
}

func (s *FalconServer) Run() error {
	err := s.listDevice()
	if err != nil {
		log.Fatalf("list device error: %v", err)
	}

	pluginapi.RegisterDevicePluginServer(s.srv, s)
	err = syscall.Unlink(DevicePluginPath + falconSocket)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	l, err := net.Listen("unix", DevicePluginPath+falconSocket)
	if err != nil {
		return err
	}

	go func() {
		lastCrashTime := time.Now()
		restartCount := 0
		for {
			log.Printf("start GPPC server for '%s'", resourceName)
			err = s.srv.Serve(l)
			if err == nil {
				break
			}

			log.Printf("GRPC server for '%s' crashed with error: $v", resourceName, err)

			if restartCount > 5 {
				log.Fatal("GRPC server for '%s' has repeatedly crashed recently. Quitting", resourceName)
			}
			timeSinceLastCrash := time.Since(lastCrashTime).Seconds()
			lastCrashTime = time.Now()
			if timeSinceLastCrash > 3600 {
				restartCount = 1
			} else {
				restartCount++
			}
		}
	}()

	// listen to reconfig event
	var config *rest.Config
	var clientset *kubernetes.Clientset
    config, err = rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }

    clientset, err = kubernetes.NewForConfig(config)
	opts := metav1.ListOptions{
       FieldSelector: "involvedObject.kind=Pod",
       Watch:         true,
    }
    watcher, err := clientset.CoreV1().Events("").Watch(context.TODO(), opts)
    if err != nil {
        panic(err.Error())
    }
    
    go func() {
        for {
			select {
			case event := <- watcher.ResultChan():
				ev, _ := event.Object.(*v1.Event)

				if ev.Reason != "Reconfig" {
					continue
				}

				name := ev.InvolvedObject.Name
            	namespace := ev.InvolvedObject.Namespace
				cur_pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
				
				// skip if the pod does not exist
				if err != nil {
					continue
				}  
				
				// if already scheduled then skip
				if cur_pod.Status.Phase != "Pending" {
					continue
				}

				log.Printf("reconfig event: %s, %s, %s\n", ev.InvolvedObject.Kind, ev.InvolvedObject.Name, ev.Reason)
				s.deviceCheckInterval = 0 * time.Second
				time.Sleep(180 * time.Second)
				s.deviceCheckInterval = 20 * time.Second
			}
		}
    }()

	// Wait for server to start by lauching a blocking connection
	conn, err := s.dial(falconSocket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

// register to Kubelet
func (s *FalconServer) RegisterToKubelet() error {
	socketFile := filepath.Join(DevicePluginPath + KubeletSocket)

	conn, err := s.dial(socketFile, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	req := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(DevicePluginPath + falconSocket),
		ResourceName: resourceName,
	}
	log.Infof("Register to kubelet with endpoint %s", req.Endpoint)
	_, err = client.Register(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

// GetDevicePluginOptions returns options to be communicated with Device
// Manager
func (s *FalconServer) GetDevicePluginOptions(ctx context.Context, e *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	log.Infoln("GetDevicePluginOptions called")
	return &pluginapi.DevicePluginOptions{PreStartRequired: true}, nil
}

// ListAndWatch returns a stream of List of Devices
// Whenever a Device state change or a Device disappears, ListAndWatch
// returns the new list
func (s *FalconServer) ListAndWatch(e *pluginapi.Empty, srv pluginapi.DevicePlugin_ListAndWatchServer) error {
	log.Infoln("ListAndWatch called")
	
	devs := make([]*pluginapi.Device, len(s.devices))
	i := 0
	for _, dev := range s.devices {
		devs[i] = dev
		i++
	}

	log.Infoln("devs: ", devs)
	err := srv.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
	if err != nil {
		log.Errorf("ListAndWatch send device error: %v", err)
		return err
	}

	old_devs := make([]*pluginapi.Device, len(s.devices))
	for {
		select {
		case <-s.ctx.Done():
			log.Info("ListAndWatch exit")
			return nil
		default:
			startTime := time.Now()
			err := s.listDevice()
			if err != nil {
				continue
			}

			devs := make([]*pluginapi.Device, len(s.devices))
			keys := make([]string, len(s.devices))

			// Collect keys of the map
			i := 0
			for k, _ := range s.devices {
				keys[i] = k
				i++
			}
			sort.Strings(keys)
			// Iterate over the map by key with an order
			i = 0
			for _, k := range keys {
				devs[i] = s.devices[k]
				i++
			}
			
			endTime := time.Now()
			duration := endTime.Sub(startTime)
			
			if !reflect.DeepEqual(devs, old_devs) {
				log.Infoln("num:", len(devs), "time:", duration, "next check:", s.deviceCheckInterval)
				log.Infoln(keys, devs)
				srv.Send(&pluginapi.ListAndWatchResponse{Devices: devs})	
			}
			old_devs = make([]*pluginapi.Device, len(s.devices))
			copy(old_devs, devs)
		}
		<-time.After(s.deviceCheckInterval)
	}
}

// Allocate is called during container creation so that the Device
// Plugin can run device specific operations and instruct Kubelet
// of the steps to make the Device available in the container
func (s *FalconServer) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	log.Infoln("Allocate called")
	resps := &pluginapi.AllocateResponse{}
	for _, req := range reqs.ContainerRequests {
		log.Infof("received request: %v", strings.Join(req.DevicesIDs, ","))
		gpuIDs := []string{}
		for _, devID := range req.DevicesIDs {
			gpuIDs = append(gpuIDs, s.gpuLookUp[devID])
		}
		resp := pluginapi.ContainerAllocateResponse{
			Envs: map[string]string{
				"FALCON_DEVICES": strings.Join(req.DevicesIDs, ","),
				"NVIDIA_VISIBLE_DEVICES": strings.Join(gpuIDs, ","),
			},
		}
		resps.ContainerResponses = append(resps.ContainerResponses, &resp)
	}
	return resps, nil
}

// PreStartContainer is called, if indicated by Device Plugin during registeration phase,
// before each container start. Device plugin can run device specific operations
// such as reseting the device before making devices available to the container
func (s *FalconServer) PreStartContainer(ctx context.Context, req *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	log.Infoln("PreStartContainer called")
	log.Infof("received request: %v", req)
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (s *FalconServer) get_resource_api(val int, endpoint string, ch chan DevicePair, ch_err chan error, wg *sync.WaitGroup) {
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
			if res["Host_Port#"] == s.hostPort {
				dp := DevicePair{
					dev_ID: strconv.Itoa(val) + "-" + res["Dev_GID"],
					GPU_GUID: res["OB_GPU_GUID"],
				}
				ch <- dp
			}
		}  
	}
	ch_err <- nil
	return
}

// discover device
func (s *FalconServer) listDevice() error {
	s.devices = make(map[string]*pluginapi.Device)
	result := make(chan DevicePair, 8)
	err := make(chan error, 2)

	wg := sync.WaitGroup{}

	for i, endpoint := range s.apiEndpoints {
		wg.Add(1)
		go s.get_resource_api(i, endpoint, result, err, &wg)
	}
	wg.Wait()
	close(result)
	close(err)

	for e := range err {
		if e != nil {
			return e
		}
	}

	for dp := range result {
		s.gpuLookUp[dp.dev_ID] = dp.GPU_GUID
		s.devices[dp.GPU_GUID] = &pluginapi.Device{
			ID:     dp.dev_ID,
			Health: pluginapi.Healthy,
		}
	}
	return nil
}

func (s *FalconServer) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *FalconServer) GetPreferredAllocation(context.Context, *v1beta1.PreferredAllocationRequest) (*v1beta1.PreferredAllocationResponse, error) {
	return &v1beta1.PreferredAllocationResponse{}, nil
}
