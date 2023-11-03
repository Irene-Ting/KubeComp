module my-device-plugin

go 1.13

replace (
	k8s.io/api => k8s.io/api v0.26.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.2-rc.0
	k8s.io/apiserver => k8s.io/apiserver v0.26.1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.26.1
	k8s.io/client-go => k8s.io/client-go v0.26.1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.26.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.1
	k8s.io/code-generator => k8s.io/code-generator v0.26.2-rc.0
	k8s.io/component-base => k8s.io/component-base v0.26.1
	k8s.io/cri-api => k8s.io/cri-api v0.26.2-rc.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.26.1
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.1
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.26.1
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.26.1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.26.1
	k8s.io/kubelet => k8s.io/kubelet v0.26.1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.26.1
	k8s.io/metrics => k8s.io/metrics v0.26.1
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.26.1
)

require (
	github.com/Rican7/retry v0.1.0 // indirect
	github.com/auth0/go-jwt-middleware v0.0.0-20170425171159-5493cabe49f7 // indirect
	github.com/bazelbuild/bazel-gazelle v0.0.0-20181012220611-c728ce9f663e // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/cespare/prettybench v0.0.0-20150116022406-03b8cfe5406c // indirect
	github.com/cloudflare/cfssl v0.0.0-20180726162950-56268a613adf // indirect
	github.com/clusterhq/flocker-go v0.0.0-20160920122132-2b8b7259d313 // indirect
	github.com/codedellemc/goscaleio v0.0.0-20170830184815-20e2ce2cf885 // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	github.com/containernetworking/cni v0.6.0 // indirect
	github.com/coreos/rkt v1.30.0 // indirect
	github.com/cpuguy83/go-md2man v1.0.4 // indirect
	github.com/d2g/dhcp4 v0.0.0-20170904100407-a1d1b6c41b1c // indirect
	github.com/d2g/dhcp4client v0.0.0-20170829104524-6e570ed0a266 // indirect
	github.com/docker/libnetwork v0.0.0-20180830151422-a9cd636e3789 // indirect
	github.com/fsnotify/fsnotify v1.6.0
	github.com/go-openapi/loads v0.17.2 // indirect
	github.com/go-openapi/spec v0.17.2 // indirect
	github.com/go-openapi/validate v0.18.0 // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible // indirect
	github.com/godbus/dbus v0.0.0-20151105175453-c7fdd8b5cd55 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d // indirect
	github.com/gophercloud/gophercloud v0.0.0-20190126172459-c818fa66e4c8 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/heketi/heketi v0.0.0-20181109135656-558b29266ce0 // indirect
	github.com/heketi/rest v0.0.0-20180404230133-aa6a65207413 // indirect
	github.com/heketi/tests v0.0.0-20151005000721-f3775cbcefd6 // indirect
	github.com/heketi/utils v0.0.0-20170317161834-435bc5bdfa64 // indirect
	github.com/jteeuwen/go-bindata v0.0.0-20151023091102-a0ff2567cfb7 // indirect
	github.com/kardianos/osext v0.0.0-20150410034420-8fef92e41e22 // indirect
	github.com/kr/fs v0.0.0-20131111012553-2788f0dbd169 // indirect
	github.com/lpabon/godbc v0.1.1 // indirect
	github.com/mattn/go-shellwords v0.0.0-20180605041737-f8471b0a71de // indirect
	github.com/mesos/mesos-go v0.0.9 // indirect
	github.com/mholt/caddy v0.0.0-20180213163048-2de495001514 // indirect
	github.com/mvdan/xurls v0.0.0-20160110113200-1b768d7c393a // indirect
	github.com/pkg/sftp v0.0.0-20160930220758-4d0e916071f6 // indirect
	github.com/pquerna/ffjson v0.0.0-20180717144149-af8b230fcd20 // indirect
	github.com/quobyte/api v0.1.2 // indirect
	github.com/robfig/cron v0.0.0-20170309132418-df38d32658d8 // indirect
	github.com/russross/blackfriday v0.0.0-20151117072312-300106c228d5 // indirect
	github.com/sigma/go-inotify v0.0.0-20181102212354-c87b6cf5033d // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/storageos/go-api v0.0.0-20180912212459-343b3eff91fc // indirect
	github.com/urfave/negroni v1.0.0 // indirect
	github.com/vmware/photon-controller-go-sdk v0.0.0-20170310013346-4a435daef6cc // indirect
	github.com/xanzy/go-cloudstack v0.0.0-20160728180336-1e2cbf647e57 // indirect
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	gonum.org/v1/gonum v0.0.0-20190331200053-3d26580ed485 // indirect
	google.golang.org/grpc v1.49.0
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.26.1
	k8s.io/apimachinery v0.26.1
	k8s.io/client-go v0.26.1
	k8s.io/heapster v1.2.0-beta.1 // indirect
	k8s.io/klog v0.3.1 // indirect
	k8s.io/kubelet v0.26.1
	k8s.io/kubernetes v1.26.1
	k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3 // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	vbom.ml/util v0.0.0-20160121211510-db5cfe13f5cc // indirect
)

replace k8s.io/component-helpers => k8s.io/component-helpers v0.26.1

replace k8s.io/controller-manager => k8s.io/controller-manager v0.26.1

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.26.1

replace k8s.io/kms => k8s.io/kms v0.26.2-rc.0

replace k8s.io/kubectl => k8s.io/kubectl v0.26.1

replace k8s.io/mount-utils => k8s.io/mount-utils v0.26.2-rc.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.26.1

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.26.1

replace k8s.io/sample-controller => k8s.io/sample-controller v0.26.1
