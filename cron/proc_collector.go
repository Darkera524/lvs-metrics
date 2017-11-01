package cron

import (
	"github.com/shirou/gopsutil/process"
	"fmt"
	"github.com/open-falcon/common/model"
	"github.com/mesos-utility/lvs-metrics/g"
	"time"
	"strings"
)

type ProcessInfo struct {
	pid int32
	cmdline string
	//excutablePath string
	//workingDirctory string
	CPUPercent float64
	MemPercent float32
	FileDescriptorNum int32
	ThreadNum int32
}

type MemInfo struct {

}

func Test(){
	fmt.Println("qweqe")
}

func CollectProc() {
	pids, err := process.Pids()


	if err != nil {
		//error handle
		fmt.Println("error:1")
		fmt.Println(err.Error())
		return
	}
	cpuInfoList,err := collect_info(pids)
	if err != nil {
		//error handle
		fmt.Println("error:2")
		fmt.Println(err.Error())
		return
	}

	proc_metrics ,_ := convirtProcessInfoToMetrics(cpuInfoList)
	g.SendMetrics(proc_metrics)
	/*for i:=0;i<len(cpuInfoList);i++{
		fmt.Print("pid:",cpuInfoList[i].pid,"cmdline:",cpuInfoList[i].cmdline ,"cpu:",cpuInfoList[i].CPUPercent,"mem:",cpuInfoList[i].MemPercent,"fdn:",cpuInfoList[i].FileDescriptorNum,"thread:",cpuInfoList[i].ThreadNum,"\n")
	}*/

}

func collect_info(pids []int32) (CPUInfoList []*ProcessInfo,err error) {
	for _, pid := range pids{
		proc, err := process.NewProcess(pid)
		if err != nil {
			fmt.Println(err.Error())
			return CPUInfoList,err
		}
		var singleInfo *ProcessInfo

		CPUPercent,err := proc.CPUPercent()
		if err !=nil {
			return CPUInfoList,err
		}

		cmdline, err := proc.Cmdline()
		if err !=nil {
			return CPUInfoList,err
		}

		/*excutablePath, err := proc.Exe()
		if err !=nil {
			return CPUInfoList,err
		}*/

		/*workingDerectory, err := proc.Cwd()
		if err !=nil {
			return CPUInfoList,err
		}*/

		memPercent, err := proc.MemoryPercent()
		if err !=nil {
			return CPUInfoList,err
		}

		fileDescriptiorNum, err := proc.NumFDs()
		if err !=nil {
			return CPUInfoList,err
		}

		threadNum, err := proc.NumThreads()
		if err !=nil {
			return CPUInfoList,err
		}

		singleInfo = &ProcessInfo{
			pid:pid,
			cmdline:cmdline,
			//excutablePath:excutablePath,
			//workingDirctory:workingDerectory,
			CPUPercent:CPUPercent,
			MemPercent:memPercent,
			FileDescriptorNum:fileDescriptiorNum,
			ThreadNum:threadNum,
		}
		CPUInfoList = append(CPUInfoList, singleInfo)


	}


	return CPUInfoList,nil
}

func convirtProcessInfoToMetrics(procInfo []*ProcessInfo)(metrics []*model.MetricValue, err error){
	hostname, _ := g.Hostname()
	now := time.Now().Unix()
	var tags string
	var attachtags = g.Config().AttachTags
	var interval int64 = g.Config().Transfer.Interval
	if attachtags != "" {
		tags = attachtags
	}
	for i:=0;i<len(procInfo);i++{
		cmdline := (strings.Split(procInfo[i].cmdline," "))[0]
		var tag string
		if tags != "" {
			tag = fmt.Sprintf("%s,pid=%d,cmdline=%s", tags, procInfo[i].pid, cmdline)
		} else {
			tag = fmt.Sprintf("pid=%d,cmdline=%s", procInfo[i].pid, cmdline)
		}
		singleMetric := &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.cpu.percent",
			Value:     procInfo[i].CPUPercent,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)

		singleMetric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.mem.percent",
			Value:     procInfo[i].MemPercent,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)

		singleMetric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.fd.num",
			Value:     procInfo[i].FileDescriptorNum,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)

		singleMetric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "proc.num.thread",
			Value:     procInfo[i].ThreadNum,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, singleMetric)
	}
	return metrics,nil
}

