package configuration

import (
	"fmt"
	"time"
)

var current Application
var dateStarted time.Time

func init() {
	current.Name = "api-klit"
	current.Version = "0.0.1"
	dateStarted = time.Now()
}

type Application struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Build   struct {
		Number string `json:"number"`
		Date   string `json:"date"` // ISO8601 UTC ?
	} `json:"build"`
	UpTime        string `json:"upTime"`
	Configuration struct {
		Source string `json:"source"`
		Data   Data   `json:"data"`
	} `json:"configuration"`
}

func SetBuildInfo(gitHash string, date string) {
	current.Build.Number = gitHash
	current.Build.Date = date
}

func setConfig(source string, data Data) {
	current.Configuration.Source = source
	current.Configuration.Data = data
}

func Get() Application {
	var app = current
	app.UpTime = time.Since(dateStarted).String()
	return app
}

func String() string {
	return fmt.Sprintf("%v %v (%v %v)", current.Name, current.Version, current.Build.Number, current.Build.Date)
}
