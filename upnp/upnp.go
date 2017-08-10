package upnp

import (
	"github.com/dustinmj/renotts/com"
	"github.com/dustinmj/renotts/config"
	"github.com/fromkeith/gossdp"
	"strings"
	"time"
)

type lst struct {
}

func (b lst) Response(m gossdp.ResponseMessage) {}

// DVPATH - path to device type
const DVPATH = "/device_description.xml"

//SN - serial number used to identify urn:..:basic:1 as RenoTTS in upnp
const SN = "5666482213265F"

// DT - the device type
const DT = "urn:schemas-dustinjorge-com:device:TTSEngine:1"

//const DT = "urn:schemas-upnp-org:device:ZonePlayer:1"

// FN - the device 'friendly name'
var FN = com.AppName + ": " + com.GetOutboundIP().String()

// M - the device 'manufacturer'
const M = "Dustin Jorge"

// MN - the device model name
const MN = com.AppVers

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

// server defaults, scoped for re-broadcast calls
var serverDef gossdp.AdvertisableServer

func init() {
	// generate UUID
	mcs, err := getMacStr()
	if err != nil {
		com.Msg("Could not reliably determine mac address, unable to start UPNP.")
	}
	UUID = UUIDB + mcs

	// create server defaults
	serverDef = gossdp.AdvertisableServer{
		ServiceType: DT,
		DeviceUuid:  UUID,
		Location:    "http://" + com.GetOutboundIP().String() + config.Val("port") + DVPATH,
		MaxAge:      MAXAGE,
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
	com.Msg("UPNP Server Started.")
	s.AdvertiseServer(serverDef) // library re-adverts correctly
	for {
		select {
		case <-sig:
			return nil
		default:
			time.Sleep(time.Second)
			break
		}
	}
}

func getMacStr() (string, error) {
	mcs, err := com.GetMacAddr()
	if err != nil {
		return "", err
	}
	var mac string
	if len(mcs) > 0 {
		mac = mcs[1]
	}
	return strings.Replace(mac, ":", "", -1), nil
}
