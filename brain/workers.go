package brain

import (
	"github.com/zenthangplus/goccm"
)

func learnInBackground() {
	c := goccm.New(2000)
	for text := range textChan {
		c.Wait()

		go func(t string) {
			defer c.Done()

			Learn(t, false)  
        }(text)
	}
}