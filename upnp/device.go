package upnp

import (
	"github.com/dustinmj/renotts/coms"
	"github.com/dustinmj/renotts/config"
)

//GetDD - gets upnp device type description
func GetDD() []byte {
	dd := "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	dd += "<root xmlns=\"urn:schemas-upnp-org:device-1-0\">"
	dd += "<specVersion>"
	dd += "<major>1</major>"
	dd += "<minor>0</minor>"
	dd += "</specVersion>"
	dd += "<device>"
	dd += "<deviceType>" + coms.DeviceType + "</deviceType>"
	dd += "<friendlyName>" + FN + "</friendlyName>"
	dd += "<manufacturer>" + M + "</manufacturer>"
	dd += "<modelName>" + MN + "</modelName>"
	dd += "<serialNumber>" + SN + "</serialNumber>"
	dd += "<UDN>" + UUID + "</UDN>"
	dd += "<serviceList>"
	dd += "<service>"
	dd += "<serviceType>urn:schemas-dustinjorge-com:service:SpeakTTS:1</serviceType>"
	dd += "<controlURL>/" + config.Val("path") + "/polly/</controlURL>"
	dd += "</service>"
	dd += "</serviceList>"
	dd += "</device>"
	dd += "</root>"
	return []byte(dd)
}
