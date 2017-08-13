package upnp

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
	"github.com/fromkeith/gossdp"
	"strings"
	"time"
)

type lst struct {
}

func (b lst) Response(m gossdp.ResponseMessage) {}

// Port -- Currently active port
// set by server
var Port string

// DVPATH - path to device type
const DVPATH = "/device_description.xml"

//SN - serial number used to identify urn:..:basic:1 as RenoTTS in upnp
const SN = "5666482213265F"

// FN - the device 'friendly name'
var FN = coms.AppName + ": " + coms.GetOutboundIP().String()

// M - the device 'manufacturer'
const M = "Dustin Jorge"

// MN - the device model name
const MN = coms.AppVers

// UUIDB - UUID Base String
const UUIDB = "e658f044-7bf4-11e7-bb31-"

// UUID - UUID of Device
var UUID string

// MAXAGE - max age of UPNP broadcast validity
const MAXAGE = 3600

// REF - interval at which to re-broadcast in Seconds
const REF = 2000

// signal channel, close cast loop
var sig = make(chan int)

// ip address we're on, if this changes, restart server
var ip = ""

func init() {
	// generate UUID
	mcs, err := getMacStr()
	if err != nil {
		coms.Msg("Could not reliably determine mac address, unable to start UPNP.")
	}
	UUID = UUIDB + mcs

	// set port to default if nil
	if len(Port) > 0 {
		Port = config.Val("port")
	}
}

// Create - Start the UPNP server
func Create() {
	go mk()
}

func mk() {
	defer Stop()
	cast()
}

//Stop broadcasting and kill
func Stop() {
	sig <- 1
}

func cast() error {
	s, err := gossdp.NewSsdp(nil)
	if err != nil {
		return err
	}
	defer s.Stop()
	go s.Start()
	// store ip for future checks
	ip = coms.GetOutboundIP().String()
	// create server defaults
	serverDef := defs(ip)
	s.AdvertiseServer(serverDef) // library re-adverts correctly
	for {
		select {
		case <-sig:
			return nil
		default:
			time.Sleep(time.Second * 5)
			// check ip every 5 Seconds
			if ip != coms.GetOutboundIP().String() {
				s.RemoveServer(UUID)
				ip = coms.GetOutboundIP().String()
				serverDef = defs(ip)
				s.AdvertiseServer(serverDef)
			}
			break
		}
	}
}

func defs(ip string) gossdp.AdvertisableServer {
	return gossdp.AdvertisableServer{
		ServiceType: coms.DeviceType,
		DeviceUuid:  UUID,
		Location:    "http://" + ip + Port + DVPATH,
		MaxAge:      MAXAGE,
	}
}

func getMacStr() (string, error) {
	mcs, err := coms.GetMacAddr()
	if err != nil {
		return "", err
	}
	var mac string
	if len(mcs) > 0 {
		mac = mcs[1]
	}
	return strings.Replace(mac, ":", "", -1), nil
}
