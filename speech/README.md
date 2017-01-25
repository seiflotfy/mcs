# Bing Speech Recogntition for Go


## Howto use



```go
token, _ := speech.GetToken()
file, _ := os.Open(os.Args[1])
reco, _ := speech.NewSpeechRequest("en-US", "linux", "websearch", "1d4b6030-9099-11e0-91e4-0800200c9a66", token, file)
res, _ := speech.Recognize(reco)
```

**Note:** Make sure the `MSFT_SPEECH_API_KEY` env variable is set.