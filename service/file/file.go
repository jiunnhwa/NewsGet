package file

import (
	"log"
	"os"
	"sync"
)

//Append text to the file.
func Appender(fname, text string) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		log.Println(err)
	}
}
