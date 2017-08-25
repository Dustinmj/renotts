package tmplt

//BootHTML - template for showing boot configuration clues
var BootHTML = `
{{define "content"}}
    <h3>For Debian Based Systems <small>Rasbian, Jesse, Chip OS, Ubuntu, Debian, etc.</small></h3>
    <p>RenoTTS <strong>should not be run as root user</strong>. Please follow these directions to run RenoTTS as <{{.User}}> each time the system boots.</p>
    <p><strong>1) Check and Test your setup:</strong>
    <br /><br /><a href="{{.ConfigCheckURL}}">{{.ConfigCheckURL}}</a>
    <br /><a href="{{.TestURL}}">{{.TestURL}}</a></pre></p>
    <p><strong>2) Open /etc/rc.local:</strong>
    <pre><code>sudo nano /etc/rc.local</code></pre>
    </p>
    <p><strong>3) Copy/paste to rc.local (before 'exit 0'):</strong>
    <pre><code># Start RenoTTS as {{.User}}
su {{.User}} -c '{{.AppPath}} >> {{.LogFile}} 2>&1 &'</code></pre>
    </p>
    <p><strong>4) Save and close rc.local</strong></p>
    <p><strong>5) Make sure {{.User}} is a member of the audio group.</strong>
    <pre><code>sudo adduser {{.User}} audio</code></pre></p>
    <p><strong>6) Reboot and test</strong>
    <pre><code>sudo reboot now</code></pre>
    <p><strong>Please note:</strong></p>
    <ol>
        <li><strong>If you change the location of RenoTTS, you will need to re-configure your boot settings.</strong></li>
        <li>If RenoTTS crashes, it will not restart itself until next boot.</li>
        <li>If you are having problems, check the log file at {{.LogFile}}. This file will be erased after every boot.</li>
        <li>RenoTTS configuration file will live at <strong>{{.ConfigFile}}</strong>.</li>
        <li>Reno cache path will be as set in RenoTTS configuration file and must be writable by {{.User}}.</li>
        <li>If you haven't specified a specific port in your configuration, you will not know what port RenoTTS is on but UPnP implementations (such as the Smarthings handlers) will be able to find RenoTTS just fine.</li>
    </ol>
    <p>sudo alsactl init, sudo adduser dustin audio</p>
{{end}}
`
