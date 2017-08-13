package coms

//Instruct - instructions sent to the user
var Instruct = []byte(`
RenoTTS Configuration:

Configuration file lives either in the same directory
as the binary, or in $HOME/.renotts. Configuration file should be
named easytts.toml.

// begin renotts.toml

port = <port>
cachepath = <cachepath>
aws-config-profile = <optional>

// end renotts.toml

UPnP:

UPnP Device Type: ` + DeviceType + `
UPnP Service Types: urn:schemas-dustinjorge-com:service:SpeakTTS:1

Usage:

Input: JSON data in POST, Output: text/plain, message
Response Codes:
  200 (ok, from cache)
  205 (ok, queried/downloaded)
  400 (failed, bad request)

JSON Input Format
{
  "text":"<text>",
  "voice":"<voicename>",
  "sampleRate":"<samplerate>"
}

Voices:

Geraint | Gwyneth | Mads | Naja |
Hans | Marlene | Nicole | Russell |
Amy | Brian | Emma | Raveena | Ivy |
Joanna | Joey | Justin | Kendra |
Kimberly | Salli | Conchita | Enrique |
Miguel | Penelope | Chantal | Celine |
Mathieu | Dora | Karl | Carla | Giorgio |
Mizuki | Liv | Lotte | Ruben | Ewa |
Jacek | Jan | Maja | Ricardo | Vitoria |
Cristiano | Ines | Carmen | Maxim | Tatyana |
Astrid | Filiz

sampleRates:

8000 | 16000 | 22050

Plain text max 3000 characters, will
be pruned if sent data exceeds 3k chars

`)
