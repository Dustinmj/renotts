package tmplt

//ConfigFl - default config template
var ConfigFl = `# renotts config options
# uncomment to set, comments are made with #

# port: Port the server will listen on. This does not necessarily have to be
# specified as RenoTTS is UPNP discoverable and ip/port should not be required to be
# specified statically and should be updated correctly by implementations as they change.
# If you wish to let the system decide port number, set port to 0.
# default: "0"

{{.Port}}

# path: Path to tts functionality on server. This is UPNP discoverable and is updated
# in device-description.xml. Users should not have to change this but are welcomed to.
# default: "tts"

{{.Path}}

# cache-path: Path to cache folder on system which contains mp3 files created from aws
# streams. This can get quite large if you have many distinct tts requests. You may
# want to move this cache to a USB drive if using a SOC like rasberry pi with limited
# drive space.
# default: "~/.renotts/cache"

{{.Cachepath}}

# exec-player: Uncomment this option to have RenoTTS execute an outside player rather
# than attempting to play mp3 files internally. The player must be available in $PATH and
# be able to play mp3 files. Examples are "mplayer", "mpg123", "aplay" etc.
# default: commented/play internally

{{.Execplayer}}

# aws-config-profile: You can specify the aws config profile you wish to use. AWS
# allows multiple profiles to be specified in config and credentials files. You may
# wish to give polly different keys than the default profile or you may have many
# aws profile key sets on your machine.
# default: commented/AWS default

{{.Awsconfigprofile}}

# force-portaudio-buffer-size: Buffer size for use when renotts decodes and plays files
# internally. Setting this to 0 causes the underlying PortAudio to choose the best
# buffer size for the given hardware. However, if you experience underruns or distortion,
# you may need manually set this value. common options range from 1000-10000 (max 10000). Increasing
# buffer size has the effect of increasing memory used during playback.
# default: 0 ... let PortAudio decide buffer size

{{.ForceBufferSize}}`
