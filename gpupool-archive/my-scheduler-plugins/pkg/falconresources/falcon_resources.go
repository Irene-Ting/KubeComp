package falconresources

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedv1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"encoding/json"
)

// FalconResources is a plugin that see the GPU as a composable device
type FalconResources struct {
	handle 	framework.Handle
	k8scli	*kubernetes.Clientset
}

var _ framework.PreFilterPlugin = &FalconResources{}
var _ framework.ScorePlugin = &FalconResources{}
var _ framework.PermitPlugin = &FalconResources{}
// var _ framework.PreBindPlugin = &FalconResources{}

// Name is the name of the plugin used in Registry and configurations.
const Name = "FalconResources"

func (gp *FalconResources) Name() string {
	return Name
}

// New initializes and returns a new FalconResources plugin.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	gp := FalconResources{}
	gp.handle = h

	config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }
	gp.k8scli, err = kubernetes.NewForConfig(config)
	if err != nil {
        panic(err.Error())
    }
	return &gp, nil
}

// filter the pod if the gpu count in the gpu pool is less than the required amount
func (gp *FalconResources) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	log.Printf("%s in Prefilter", pod.Name)
	// add annotation
	// k8scfg, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// k8scli, err := kubernetes.NewForConfig(k8scfg)
	// if err != nil {
	// 	panic(err)
	// }
	
	// if required gpu is more than all availabe gpu, then return failure early
	required_falcon_quantity := pod.Spec.Containers[0].Resources.Requests["falcon.com/gpu"]
	required_falcon, _ := required_falcon_quantity.AsInt64()

	nodeinfos, _ := gp.handle.SnapshotSharedLister().NodeInfos().List()
	total_falcon := int64(0)
	for _, nodeinfo := range nodeinfos {
		total_falcon += (nodeinfo.Allocatable.ScalarResources["falcon.com/gpu"] - nodeinfo.Requested.ScalarResources["falcon.com/gpu"])
	}

	log.Printf("%s requires %d gpu(s), and currently has %d gpus(s) in total\n", pod.Name, required_falcon, total_falcon)

	patchAnnotations := map[string]interface{} {
	"metadata": map[string]map[string]string{"annotations": {
		"use_falcon": "true",
	}}}
	newStatus := framework.NewStatus(framework.Success, "")

	// the pod doesn't use Falcon
	if required_falcon == 0 {
		patchAnnotations = map[string]interface{}{
		"metadata": map[string]map[string]string{"annotations": {
			"use_falcon": "false",
		}}}
	}

	// total gpu is less than required, so it's useless to do reconfig
	if total_falcon < required_falcon {
		// patchAnnotations = map[string]interface{}{
		// "metadata": map[string]map[string]string{"annotations": {
		// 	"use_falcon": "true",
		// 	"permit_reconfig": "false",
		// }}}
		reason := fmt.Sprintf("%s requires %d gpu but only %d gpu in the pool.", pod.Name, required_falcon, total_falcon)
		newStatus = framework.NewStatus(framework.Unschedulable, reason)
	}

	patchedAnnotationBytes, _ := json.Marshal(patchAnnotations)
	_, err := gp.k8scli.CoreV1().Pods(pod.Namespace).Patch(ctx, pod.Name, types.StrategicMergePatchType, patchedAnnotationBytes, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}	
	
	return nil, newStatus
}

// PreFilterExtensions returns a PreFilterExtensions interface if the plugin implements one.
func (gp *FalconResources) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}


// Score invoked at the score extension point.
func (gp *FalconResources) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := gp.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}

	required_falcon_quantity := pod.Spec.Containers[0].Resources.Requests["falcon.com/gpu"]
	required_falcon, _ := required_falcon_quantity.AsInt64()
	local_falcon := (nodeInfo.Allocatable.ScalarResources["falcon.com/gpu"] - nodeInfo.Requested.ScalarResources["falcon.com/gpu"])
	
	var score int64 = 0
	if local_falcon > required_falcon {
		score = int64(required_falcon*100 / local_falcon)	
	} else if local_falcon == required_falcon {
		score = 100
	} else {
		score = local_falcon - required_falcon
	}
	log.Printf("%s has %d gpu, %s requires %d gpu -> score: %v\n", nodeName, local_falcon, pod.Name, required_falcon, score)
	return score, nil
}


func (gp *FalconResources) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Find highest and lowest scores.
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Transform the highest to lowest score range to fit the framework's min to max node score range.
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	return nil
}

// ScoreExtensions of the Score plugin.
func (gp *FalconResources) ScoreExtensions() framework.ScoreExtensions {
	return gp
}

func (gp *FalconResources) GetGpuDemand(ctx context.Context, pod *v1.Pod, nodeName string) int {
	// k8scfg, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// k8scli, err := kubernetes.NewForConfig(k8scfg)
	// if err != nil {
	// 	panic(err)
	// }

	node, err := gp.k8scli.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	alloc_gpu_quantity := node.Status.Allocatable["falcon.com/gpu"]
	alloc_gpu, _ := alloc_gpu_quantity.AsInt64()

	nodeInfo, err := gp.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err)
		return 0
	}
	request_gpu := nodeInfo.Requested.ScalarResources["falcon.com/gpu"]

	required_falcon_quantity := pod.Spec.Containers[0].Resources.Requests["falcon.com/gpu"]
	required_falcon, _ := required_falcon_quantity.AsInt64()
	// log.Printf("%v -> alloc: %v, used: %v, required: %v", nodeName, alloc_gpu, request_gpu, required_falcon)
	
	if (alloc_gpu - request_gpu) >= required_falcon {
		return 0
	} else {
		demand := int(required_falcon - (alloc_gpu - request_gpu))
		return demand
	}
}

func (gp *FalconResources) Permit(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (*framework.Status, time.Duration) {
	var waitTime time.Duration = 0
	var retStatus *framework.Status
	retStatus = framework.NewStatus(framework.Success)

	// k8scfg, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// k8scli, err := kubernetes.NewForConfig(k8scfg)
	// if err != nil {
	// 	panic(err)
	// }

	demand := gp.GetGpuDemand(ctx, pod, nodeName)
	if demand > 0 {
		retStatus = framework.NewStatus(framework.Unschedulable)
		// annotate dst node
		patchAnnotations := map[string]interface{}{
		"metadata": map[string]map[string]string{"annotations": {
			"dst_node": nodeName,
			"gpu_demand": strconv.Itoa(demand),
		}}}
		patchedAnnotationBytes, _ := json.Marshal(patchAnnotations)
		_, err := gp.k8scli.CoreV1().Pods(pod.Namespace).Patch(ctx, pod.Name, types.StrategicMergePatchType, patchedAnnotationBytes, metav1.PatchOptions{})
		if err != nil {
			panic(err)
		}	
		
		// create an event
		scheme := runtime.NewScheme()
		_ = v1.AddToScheme(scheme)

		eventBroadcaster := record.NewBroadcaster()
		eventBroadcaster.StartStructuredLogging(0)
		eventBroadcaster.StartRecordingToSink(&typedv1core.EventSinkImpl{Interface: gp.k8scli.CoreV1().Events("")})
		eventRecorder := eventBroadcaster.NewRecorder(scheme, v1.EventSource{})
		eventRecorder.Event(pod, v1.EventTypeNormal, "Reconfig", fmt.Sprintf("Pod %v needs reconfiguration", pod.Name))
		defer eventBroadcaster.Shutdown()
		
		startTime := time.Now()
		for {
			new_demand := gp.GetGpuDemand(ctx, pod, nodeName)
			// log.Printf("new demand: %d\n", new_demand)
			if new_demand == 0 {
				retStatus = framework.NewStatus(framework.Success)
				break
			}
			if time.Since(startTime) > (10 + time.Duration(demand) * 30) * time.Second {
				break
			}
		}
	}
	return retStatus, waitTime
}
