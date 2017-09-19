# RenoTTS (BETA)
[http://renotts.com/](https://renotts.com/)

UPNP enabled TTS service for interfacing with AWS Polly to be installed on unix-based SOC or PC. Enables 'talking speaker' speech synthesis device.

Built for use with SmartThings using handlers available in my repo, but could be used for many other purposes.

## Features
- UPNP discoverable.
- Caches files for future playback.
- Plays files internally or externally.
- Handles all AWS communication.
- Simple configuration.

## Dependancies:
```
portaudio19-dev
libmpg123-dev
```

## Uses/Credits
- [viper](https://github.com/spf13/viper)
- [AWS Polly](https://github.com/aws/aws-sdk-go/tree/master/service/polly)
- [gossdp](https://github.com/fromkeith/gossdp)
- [go-portaudio](https://github.com/gordonklaus/portaudio)
- [go-mpg123](https://github.com/bobertlo/go-mpg123/)
