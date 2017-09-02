package tmplt

//SystemdFl - template for systemd file
var SystemdFl = `[Unit]
Description=RenoTTS Service
After=network-online.target sound.target
AssertFileIsExecutable={{.AppPath}}

[Service]
User={{.User}}
ExecStart={{.AppPath}}
Restart=on-failure
RestartSec=2
KillMode=process

[Install]
WantedBy=multi-user.target
`
