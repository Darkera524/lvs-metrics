package cron

import (
	"fmt"

	"github.com/google/seesaw/ipvs"
)

type RealServer struct {
	IP     string
	Port int
	VIP string
	VPort int

	//statistic
	ActiveConns uint32
	InactiveConns uint32
	Connections uint32
	PacketsIn uint32
	PacketsOut uint32
	BytesIn uint64
	BytesOut uint64
	CPS uint32
	PPSIn uint32
	PPSOut uint32
	BPSIn uint32
	BPSOut uint32
}

/**func NewRealServer(end string, actconn, inactconn int) *RealServer {
	return &RealServer{
		IP:     end,
		ActiveConn: actconn,
		InActConn:  inactconn,
	}
}**/

type VirtualIPPoint struct {
	IP            string
	Port          int
	ActiveConns   uint32
	InactiveConns uint32
	RealServerNum int

	// stats
	Connections uint32
	PacketsIn   uint32
	PacketsOut  uint32
	BytesIn     uint64
	BytesOut    uint64
	//Realservers     [](*RealServer)
}

func NewVirtualIPPoint(ip string, port int, actconn, inactconn uint32) *VirtualIPPoint {
	return &VirtualIPPoint{
		IP:            ip,
		Port:          port,
		ActiveConns:   actconn,
		InactiveConns: inactconn,
	}
}

func GetIPVSStats() (vips []*VirtualIPPoint, rips []*RealServer, err error) {
	svcs, err := ipvs.GetServices()
	if err != nil {
		return nil, nil, err
	}

	var vip *VirtualIPPoint
	var ActiveConns uint32
	var InactiveConns uint32
	var RsNum int
	//var PersistConns uint32
	for _, svc := range svcs {
		ActiveConns = 0
		InactiveConns = 0
		RsNum = len(svc.Destinations)

		ipstr := fmt.Sprintf("%v", svc.Address)

		for _, dest := range svc.Destinations {
			var rip *RealServer
			ripstr := fmt.Sprintf("%v", dest.Address)
			rip = &RealServer{
				IP : ripstr,
				Port: int(dest.Port),
				VIP: ipstr,
				VPort: int(svc.Port),
				ActiveConns: dest.Statistics.ActiveConns,
				InactiveConns: dest.Statistics.InactiveConns,
				Connections: dest.Statistics.Stats.Connections,
				PacketsIn: dest.Statistics.Stats.PacketsIn,
				PacketsOut: dest.Statistics.Stats.PacketsOut,
				BytesIn: dest.Statistics.Stats.BytesIn,
				BytesOut: dest.Statistics.Stats.BytesOut,
				CPS: dest.Statistics.Stats.CPS,
				PPSIn: dest.Statistics.Stats.PPSIn,
				PPSOut: dest.Statistics.Stats.PPSOut,
				BPSIn: dest.Statistics.Stats.BPSIn,
				BPSOut: dest.Statistics.Stats.BPSOut,
			}
			rips = append(rips, rip)

			ActiveConns += dest.Statistics.ActiveConns
			InactiveConns += dest.Statistics.InactiveConns
			//PersistConns += dest.Statistics.PersistConns
		}


		vip = &VirtualIPPoint{
			IP:            ipstr,
			Port:          int(svc.Port),
			ActiveConns:   ActiveConns,
			InactiveConns: InactiveConns,
			RealServerNum: RsNum,

			Connections: svc.Statistics.Connections,
			PacketsIn:   svc.Statistics.PacketsIn,
			PacketsOut:  svc.Statistics.PacketsOut,
			BytesIn:     svc.Statistics.BytesIn,
			BytesOut:    svc.Statistics.BytesOut,
		}

		vips = append(vips, vip)
	}

	return vips, rips, nil
}
