package tmplt

//DevDescFl - template for UPNP device description
var DevDescFl = `<?xml version="1.0" encoding="utf-8"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
    <specVersion>
        <major>1</major>
        <minor>0</minor>
    </specVersion>
    <device>
        <deviceType>{{.DeviceType}}</deviceType>
        <friendlyName>{{.FriendlyName}}</friendlyName>
        <manufacturer>{{.Manufacturer}}</manufacturer>
        <modelNumber>{{.ModelNumber}}</modelNumber>
        <UDN>{{.UUID}}</UDN>
        <iconList>
             <icon>
                 <mimetype>image/png</mimetype>
                 <width>120</width>
                 <height>120</height>
                 <depth>24</depth>
                 <url>https://icon.renotts.com/renotts.bl.2x.png</url>
             </icon>
        </iconList>
        <serviceList>
            {{- range $k,$v := .Services -}}
            <service>
                <serviceType>urn:schemas-dustinjorge-com:service:SpeakTTS:1</serviceType>
                <controlURL>{{$v.Path}}</controlURL>
            </service>
            {{- end -}}
        </serviceList>
        <presentationURL>/</presentationURL>
    </device>
</root>`
