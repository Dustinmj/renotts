package server

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/fromkeith/gossdp"
	"strings"
	"time"
)

type lst struct {
}

func (b lst) Response(m gossdp.ResponseMessage) {}

const (
	//upnpdtp - path to device type
	upnpdtp = "/device_description.xml"
	//upnpman - the device 'manufacturer'
	upnpman = "Dustin Jorge"
	//upnpmod - the device model name
	upnpmod = coms.AppVers
	//upnpdt - the device type
	upnpdt = "urn:schemas-dustinjorge-com:device:TTSEngine:1"
	//upnpuuidb - UUID Base String
	upnpuuidb = "e658f044-7bf4-11e7-bb31-"
	//upnpmaxage - max age of UPNP broadcast validity
	upnpmaxage = 3600
)

// UUID - UUID of Device
var upnpuuid string

// fn - the device 'friendly name'
var upnpfn = coms.AppName + ": " + getOutboundIP().String()

// signal channel, close cast loop
var sig = make(chan int)

func init() {
	// generate UUID
	mcs, err := getMacStr()
	if err != nil {
		coms.Msg("Could not reliably determine mac address, unable to start UPNP.")
	}
	upnpuuid = upnpuuidb + mcs
}

//StartUPNP - Start the UPNP server
func StartUPNP(ttsPort string) {
	go func() {
		defer func() {
			sig <- 1
		}()
		cast(ttsPort)
	}()
}

func cast(ttsPort string) error {
	s, err := gossdp.NewSsdp(nil)
	if err != nil {
		return err
	}
	defer s.Stop()
	go s.Start()
	// store ip for future checks
	ip = getOutboundIP().String()
	// create server defaults
	serverDef := defs(ip, ttsPort)
	s.AdvertiseServer(serverDef) // library re-adverts correctly
	for {
		select {
		case <-sig:
			return nil
		default:
			time.Sleep(time.Second * 5)
			// check ip every 5 Seconds
			if ip != getOutboundIP().String() {
				s.RemoveServer(upnpuuid)
				ip = getOutboundIP().String()
				serverDef = defs(ip, ttsPort)
				s.AdvertiseServer(serverDef)
			}
			break
		}
	}
}

func defs(ip string, port string) gossdp.AdvertisableServer {
	return gossdp.AdvertisableServer{
		ServiceType: upnpdt,
		DeviceUuid:  upnpuuid,
		Location:    "http://" + ip + port + upnpdtp,
		MaxAge:      upnpmaxage,
	}
}

func getMacStr() (string, error) {
	mcs, err := getMacAddr()
	if err != nil {
		return "", err
	}
	var mac string
	if len(mcs) > 0 {
		mac = mcs[1]
	}
	return strings.Replace(mac, ":", "", -1), nil
}
