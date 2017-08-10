package config

//Instruct - instructions sent to the user
var Instruct = []byte(`
***** EasyTTSBox Configuration

Configuration file lives either in the same directory
as the binary, or in $HOME/.renotts

Configuration file should be
named easytts.toml containing:

// **** begin renotts.toml

/*
* port:      port to listen on, e.g. 8080
* cachepath: path to cache location, must be writable
aws-config-profile // TODO INSTRUCT ABOUT
*/

port = "<port>"
cachepath = "<cachepath>"

// *** end renotts.toml


***** Using RenoTTS

Input: JSON data in POST
Output: text/plain, message
Response Codes:
  200 (ok, from cache)
  205 (ok, queried/downloaded)
  400 (failed, bad request)

JSON Input Format (enclose all values in "")
{
  //* text data to translate to speech
"text":"<text>",
  //* aws voice to use in translation
"voice":"<voicename>",
  //* sample rate for returned file
"sampleRate":"<samplerate>"
}

***** Valid "voice" Values (Joanna is the best, by far)

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

***** Valid "sampleRate" Values

8000 | 16000 | 22050

***** Valid "text" Values

Plain text max 3000 characters, will
be pruned if sent data exceeds 3k chars


`)
