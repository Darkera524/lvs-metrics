# **lvs-metrics** <sup><sub>_lvs metrics collector for open-falcon_</sub></sup>

local version of lvs metrics collector for open-falcon.

changes:
 - Add the data between the rip and vip
 - Fix the bug of total data(it was rate data before)
 - Add the function of collect the info per process

## Dependencies

* [libnl3][]
* [ipvs][]
* [netlink][]
* [toolkits][]
* [glog][]

Dependencies are handled by [godep][], simple install it and type `godep restore` to fetch them.

## Install

#### The libnl3 needed in compile and production machine.

```console
# sudo yum install -y libnl3.x86_64
or
# sudo apt-get install libnl-3-dev libnl-genl-3-dev
```

#### Only needed in compile machine.
```console
$ git clone https://github.com/Darkera524/lvs-metrics.git
$ cd lvs-metrics
$ make bin
```


## Configuration

Edit cfg.json configuration file.

```console
{
    "debug": false,
    "attachtags": "",
    "http": {
        "enable": false,
        "listen": "0.0.0.0:1987"
    },
    "transfer": {
        "enable": true,
        "addr": "http://127.0.0.1:1988/v1/push", # Installed falcon agent in host.
        "interval": 30,
        "timeout": 1000
    }
}

or

{
    "debug": false,
    "attachtags": "",
    "http": {
        "enable": false,
        "listen": "0.0.0.0:1987"
    },
    "transfer": {
        "enable": true,
        "addr": "127.0.0.1:8433", # Send metrics to transfer direct.
        "interval": 30,
        "timeout": 1000
    }
}
```



## Metrics

| Counters | Type | Notes |
|-----|-----|-----|
| lvs.in.bytes | GUAGE | network in bytes per host |
| lvs.out.bytes | GUAGE | network out bytes per host |
| lvs.in.packets | GUAGE | network in packets per host |
| lvs.out.packets | GUAGE | network out packets per host |
| lvs.total.conns | GUAGE | lvs total connections per host |
| lvs.vip.active_conn | GUAGE | lvs active connections per vip now |
| lvs.vip.inact_conn | GUAGE | lvs inactive connections per vip now |
| lvs.vip.rs_num | GUAGE | lvs live realserver num per vip now |
| lvs.vip.conns | COUNTER->GAUGE | lvs conns counter from service start per vip |
| lvs.vip.inbytes | COUNTER->GAUGE | lvs inbytes counter from service start per vip |
| lvs.vip.outbytes | COUNTER->GAUGE | lvs outpkts counter from service start per vip |
| lvs.vip.inpkts | COUNTER->GAUGE | lvs inpkts counter from service start per vip |
| lvs.vip.outpkts | COUNTER->GAUGE | lvs outpkts counter from service start per vip |


[libnl3]: https://www.infradead.org/~tgr/libnl
[ipvs]: https://github.com/google/seesaw
[netlink]: https://github.com/google/seesaw
[toolkits]: https://github.com/toolkits
[glog]: https://github.com/golang/glog
[godep]: https://github.com/tools/godep


