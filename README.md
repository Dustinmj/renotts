# renotts 
[https://dustinmj.github.io/renotts/](https://dustinmj.github.io/renotts/)

UPNP enabled TTS service for interfacing with AWS Polly to be installed on unix-based SOC or PC. Enables 'talking speaker' speech synthesis device.

Built for use with SmartThings using handlers available in my repo, but could be used for many other purposes.

## Features
- UPNP
- Caching
- Plays files internally
- Handles all AWS communication

## Build dependancies (outside godeps):
```
portaudio19-dev
libmpg123-dev
```

## Uses/Credits
- [go-toml](https://github.com/pelletier/go-toml)
- [viper](https://github.com/spf13/viper)
- [go-homedir](https://github.com/mitchellh/go-homedir)
- [AWS Polly](https://github.com/aws/aws-sdk-go/tree/master/service/polly)
- [gossdp](https://github.com/fromkeith/gossdp)
- [go-portaudio](https://github.com/gordonklaus/portaudio)
- [go-mpg123](https://github.com/bobertlo/go-mpg123/)
