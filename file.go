package main

import (
	"io/ioutil"
	"log"
)

//WriteData, new or truncate
func WriteData(filename string, bytes []byte) {
	err := ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		panic("write file failed " + filename)
	}

}

//ReadFile the given filename
func ReadFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("cannot read " + filename)
		//return false, err
	}
	return string(data)
}

//WriteFile writes the data as filename, , with create or truncate
func WriteFile(filename string, data string) {
	err := ioutil.WriteFile(filename, []byte(data), 0666)
	if err != nil {
		log.Fatal("cannot read " + filename)
		//return false, err
	}
}

//FindFile in given directory
func FindFile(filename string) (bool, error) {
	files, err := ioutil.ReadDir(".") //"." for currentdir
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	for _, f := range files {
		//fmt.Println(f.Name())
		if f.Name() == filename {
			return true, nil
		}
	}
	return false, nil
}
