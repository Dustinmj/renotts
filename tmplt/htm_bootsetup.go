package tmplt

//BootHTML - template for showing boot configuration clues
var BootHTML = `
{{define "css"}}
body h2 small{display:block;font-size:.6em;font-weight:400;color:#009fdb}body>div>ul li{list-style-type:square;padding:.5em 0;color:#333}div.collapse{max-width:95%;padding:0 0 0 10px;background-color:#FFF;margin:1em 0;border:1px solid transparent}div.collapse.exp{border-color:#DDD;padding:0 0 10px 10px}div.collapse.exp p.clabel:after{content:"click to collapse"}div.collapse.exp p.clabel+div{display:block}div.collapse p.clabel:after,div.collapse.exp p.clabel:after{text-decoration:underline;font-size:70%;display:block;margin-top:1em}div.collapse p.clabel{margin-left:0;display:block;cursor:pointer;position:relative;left:-10px;max-width:50em;border-left:7px solid #009fdb;padding:1em;margin-top:0}div.collapse p.clabel small{display:block}div.collapse p.clabel:after{content:"click to expand"}div.collapse p.clabel+div{display:none}div.collapse ul li{list-style-type:none}div.collapse p>strong{display:block}div.collapse ol li{font-size:90%;max-width:50em}
{{end}}
{{define "content"}}
<h2>For Debian Based Systems <small>Rasbian, Jesse, Chip OS, Ubuntu, Debian, etc.</small></h2>
    <ul>
        <li>RenoTTS <mark>should not be run as root user</mark>.</li>
        <li>Please follow the directions for <mark>only one</mark> of the methods below.</li>
        <li>These allow RenoTTS to <mark>run as &lt;{{.User}}&gt;</mark> each time the system boots.</li>
    </ul>
    <div class='collapse'>
        <p class='clabel'>
            <strong>Method 1:</strong>
            Using systemctl
            <small>Recommended</small>
        </p>
        <div>
            <p><strong>1) Check and Test your setup:</strong>
            <ul>
                <li><a target="_blank" href="{{.ConfigCheckURL}}">{{.ConfigCheckURL}}</a></li>
                <li><a target="_blank" href="{{.TestURL}}">{{.TestURL}}</a></li>
            </ul>
            <p><strong>2) Copy the generated renotts.service to systemd:</strong>
                <code>sudo cp {{.ServiceFile}} /etc/systemd/system/</code>
            </p>
            <p><strong>3) Enable the service:</strong>
                <code>sudo systemctl enable renotts</code>
            </p>
            <p><strong>4) Make sure {{.User}} is a member of the audio group.</strong>
            <code>sudo adduser {{.User}} audio</code>
            </p>
            <p><strong>5) Reboot and test</strong>
            <code>sudo reboot</code>
            </p>
            <p><strong>6) Other Commands</strong>
            <code>sudo systemctl |enable|disable| renotts<br /><br />sudo systemctl |start|stop|restart| renotts<br /><br />systemctl status renotts<br /><br />journalctl -u renotts</code>
            </p>
            <p><strong>Please note:</strong></p>
            <ol>
                <li><mark>If you change the location of RenoTTS, you will need to re-configure your boot settings.</mark></li>
                <li>If you are having problems, check logs by using "journalctl -u renotts".</li>
                <li>RenoTTS configuration file will live at <strong>{{.ConfigFile}}</strong>.</li>
                <li>Reno cache path will be as set in RenoTTS configuration file and must be writable by {{.User}}.</li>
            </ol>
        </div>
    </div>
    <div class='collapse'>
        <p class='clabel'>
            <strong>Method 2:</strong>
            Using /etc/rc.local
            <small>Method 1 recommended</small>
        </p>
        <div>
            <p><strong>1) Check and Test your setup:</strong>
            <ul>
                <li><a target="_blank" href="{{.ConfigCheckURL}}">{{.ConfigCheckURL}}</a></li>
                <li><a target="_blank" href="{{.TestURL}}">{{.TestURL}}</a></li>
            </ul>
            <p><strong>2) Open /etc/rc.local:</strong>
            <code>sudo nano /etc/rc.local</code>
            </p>
            <p><strong>3) Copy/paste to rc.local (before 'exit 0'):</strong>
            <code># Start RenoTTS as {{.User}}<br />su {{.User}} -c '/tmp/go-build986520472/command-line-arguments/_obj/exe/renotts >> /tmp/RenoTTS.log 2>&1 &'</code>
            </p>
            <p><strong>4) Make sure {{.User}} is a member of the audio group.</strong>
            <code>sudo adduser {{.User}} audio</code>
            </p>
            <p><strong>5) Reboot and test</strong>
            <code>sudo reboot now</code>
            </p>
            <p><strong>Please note:</strong></p>
            <ol>
                <li><mark>If you change the location of RenoTTS, you will need to re-configure your boot settings.</mark></li>
                <li>If RenoTTS crashes, it will not restart itself until next boot.</li>
                <li>If you are having problems, check the log file at /tmp/RenoTTS.log. This file will be erased after every boot.</li>
                <li>RenoTTS configuration file will live at <strong>{{.ConfigFile}}</strong>.</li>
                <li>Reno cache path will be as set in RenoTTS configuration file and must be writable by {{.User}}.</li>
            </ol>
        </div>
    </div>
{{end}}{{define "javascript"}}
<script type="text/javascript">
var renotts=renotts||{};renotts.common={bootstrap:function(){var cols=document.getElementsByTagName('body')[0].querySelectorAll('div.collapse');this.setupCollapse(cols)},setupCollapse:function(cols){for(let collapse of cols){var toggle=collapse.querySelector('.clabel');if(toggle!='undefined'){toggle.addEventListener('click',function(){this.parentNode.classList.toggle('exp')}.bind(toggle))}}}}
window.addEventListener('DOMContentLoaded',function(){renotts.common.bootstrap()},!1);
</script>
{{end}}
`
