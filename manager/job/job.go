package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"newsget/model"
	"path"
	"time"
)

//*************************************************************
// Job Control
//*************************************************************
var NOW time.Time //Sync timestamping, loggers can use global NOW, or time.Now()

//Status is a record to hold status information
type Status struct {
	Code    int
	Message string
}

type DataType struct {
	AsBytes  []byte
	AsString string
	AsRows   [][]string //csv
}

func (b *DataType) Hi() {
	fmt.Println(b)
}

//Job is a record to hold job details
type Job struct {
	ID                 string
	Name               string
	StartTime, EndTime time.Time
	Request            string
	Response           string
	Result             string
	Error              Status

	Feeds    []model.Feed
	Rows     [][]string //csv
	NewsData DataType   //input source
	SQLFile  string     //out file for insert to DB
	DataType
}

//NewJob constructor
func NewJob(name string) *Job {
	NOW = time.Now()
	return &Job{ID: CreateJobID(), Name: name, StartTime: NOW}
}

//Sequence generator
var nextNum = GetNextNum(0)

//CreateJobID creates a formatted jobID string. eg: 20210410-0001
func CreateJobID() string {
	//return fmt.Sprintf("%s", time.Now().Format("20060102-150405.000"))  //OptionA: 20210410-225106.481.json

	//Using alternate method with running sequence for each day.
	return fmt.Sprintf("%s-%04d", time.Now().Format("20060102"), nextNum()) //OptionB: 20210410-0001.json
}

//GetNextNum uses a closure to generate sequences, with reset each day
func GetNextNum(startNum int) func() int {
	Day := time.Now().Day()
	currNum := startNum
	return func() int {
		if time.Now().Day() != Day {
			Day = time.Now().Day()
			return currNum
		}
		currNum += 1
		return currNum
	}
}

//Capture stores in the Result field information of interest marshalled as json string
func (j *Job) Capture(result interface{}, err error) (interface{}, error) {
	data := struct {
		Result interface{}
		Err    error
	}{
		Result: result,
		Err:    err,
	}
	bytes, _ := json.Marshal(data)
	j.Result = string(bytes)
	return result, err
}

//Save persists the Job File to disk in the log dir. eg: \logs\20210413-0001.json
func (j *Job) Save() {
	bytes, _ := json.Marshal(j)
	err := ioutil.WriteFile(path.Join("logs", j.ID+".json"), bytes, 0666)
	if err != nil {
		log.Print(err)
	}
}

//End marks the end time, and saves the job.
func (j *Job) End() {
	j.EndTime = time.Now()
	j.Save()
}
