package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/axiomhq/msft-bing-speech"
	"github.com/axiomhq/msft-bing-speech/demo/recorder"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	token, err := speech.GetToken()
	checkErr(err)

	pr, pw := io.Pipe()
	go func() {
		checkErr(recorder.Record(pw))
	}()

	reco, err := speech.NewSpeechRequest("en-US", "linux", "websearch", "1d4b6030-9099-11e0-91e4-0800200c9a66", token, pr)
	checkErr(err)

	res, err := speech.Recognize(reco)
	checkErr(err)

	if len(res.Results) > 0 {
		f, err := strconv.ParseFloat(res.Results[0].Confidence, 64)
		checkErr(err)
		fmt.Printf("I heard:\n\tPhrase: %s\n\tConfidence: %d %%\n", res.Results[0].Name, int(f*100))
	} else {
		fmt.Println("I did not hear anything")
	}
}
