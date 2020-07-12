package grpc

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/magicdvd/nacos-client"
	"google.golang.org/grpc/resolver"
)

const (
	modeHeartBeat string = "hb"
	modeSubscribe string = "sb"
)

var (
	ErrUnsupportSchema = errors.New("unsupport schema (nacos/nacoss)")
	ErrMissServiceName = errors.New("target miss service name")
	ErrMissGroupName   = errors.New("target miss group name")
	ErrMissNameSpaceID = errors.New("target miss namespace name")
	ErrMode            = errors.New("target mode err (hb/sb)")
	ErrMissInterval    = errors.New("target mode heartbeat miss interval")
	ErrInterval        = errors.New("target mode heartbeat interval error")
	ErrNoInstances     = errors.New("no valid instance")
)

type nacosResolver struct {
	nacosClient nacos.ServiceCmdable
	cc          resolver.ClientConn
	params      []nacos.Param
	mode        string
	close       chan bool
	interval    time.Duration
	serviceName string
}

func newNacosResolver(target resolver.Target, cc resolver.ClientConn) (*nacosResolver, error) {
	if target.Scheme != "nacos" && target.Scheme != "nacoss" {
		return nil, ErrUnsupportSchema
	}
	u, err := url.Parse("http://test.com" + target.Endpoint)
	if err != nil {
		return nil, err
	}
	schema := "http"
	if target.Scheme == "nacoss" {
		schema = "https://"
	}
	client, err := nacos.NewServiceClient(schema + target.Authority + "/" + u.Host + "/" + u.Path)
	if err != nil {
		return nil, err
	}
	params := make([]nacos.Param, 0)
	params = append(params, nacos.ParamHealthy(true))
	values := u.Query()
	if values.Get("cs") != "" {
		params = append(params, nacos.ParamClusters(strings.Split(values.Get("cs"), ",")))
	}
	serviceName := values.Get("s")
	if serviceName == "" {
		return nil, ErrMissServiceName
	}
	if values.Get("n") == "" {
		return nil, ErrMissNameSpaceID
	}
	params = append(params, nacos.ParamNameSpaceID(values.Get("n")))
	if values.Get("g") == "" {
		return nil, ErrMissGroupName
	}
	params = append(params, nacos.ParamGroupName(values.Get("g")))
	mode := values.Get("m")
	if mode != modeHeartBeat && mode != modeSubscribe {
		return nil, ErrMode
	}
	var interval time.Duration
	if mode == modeHeartBeat {
		s := values.Get("d")
		if s == "" && mode == modeHeartBeat {
			return nil, ErrMissInterval
		}
		tmp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		if tmp <= 0 {
			return nil, ErrInterval
		}
		interval = time.Duration(tmp) * time.Millisecond
	}
	params = append(params, nacos.ParamGroupName(values.Get("g")))
	c := &nacosResolver{
		nacosClient: client,
		cc:          cc,
		params:      params,
		mode:        mode,
		close:       make(chan bool),
		interval:    interval,
		serviceName: serviceName,
	}
	return c, nil
}

func (c *nacosResolver) start() {
	if c.mode == modeHeartBeat {
		tick := time.NewTicker(c.interval)
		service, err := c.nacosClient.GetService(c.serviceName, false, c.params...)
		if err != nil {
			c.cc.ReportError(err)
		} else {
			addrs, err := c.getInstances(service)
			if err != nil {
				c.cc.ReportError(err)
			} else {
				c.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}
		for {
			select {
			case <-tick.C:
				service, err := c.nacosClient.GetService(c.serviceName, false, c.params...)
				if err != nil {
					c.cc.ReportError(err)
				} else {
					addrs, err := c.getInstances(service)
					if err != nil {
						c.cc.ReportError(err)
					} else {
						c.cc.UpdateState(resolver.State{Addresses: addrs})
					}
				}
			case <-c.close:
				return
			}
		}
	} else {
		c.nacosClient.Subscribe(c.serviceName, func(service *nacos.Service) {
			addrs, err := c.getInstances(service)
			if err != nil {
				c.cc.ReportError(err)
			} else {
				c.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}, c.params...)
	}
}

func (c *nacosResolver) getInstances(service *nacos.Service) ([]resolver.Address, error) {
	if len(service.Instances) == 0 {
		return nil, ErrNoInstances
	}
	l := len(service.Instances)
	ret := make([]resolver.Address, l)
	for i := 0; i < l; i++ {
		ins := service.Instances[i]
		ret = append(ret, resolver.Address{
			Addr:       fmt.Sprintf("%s:%d", ins.Ip, ins.Port),
			ServerName: c.serviceName,
		})
	}
	return ret, nil
}

func (c *nacosResolver) ResolveNow(o resolver.ResolveNowOptions) {
	//directly get service
}

func (c *nacosResolver) Close() {
	close(c.close)
}
