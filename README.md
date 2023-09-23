# framework-scheduler-networtraffic
对prometheus采集的数据进行打分，参与k8s调度，源码版本为1.26.7

需要在kubernetes 1.26.7的源码环境创建文件进行 go mod tidy  
root@jcrose:/usr/local/kubernetes-1.26.7/pkg/scheduler/framework/plugins/networtraffic# tree  
.  
├── Dockerfile  
├── main.go  
├── networktraffic-scheduler  
└── plugin  
    ├── networktraffic.go  
    └── prometheus.go  
 
1 directory, 5 files

