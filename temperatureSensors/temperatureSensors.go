package TemperatureSensors

import (
  log "github.com/sirupsen/logrus"
  "regexp"
  "strconv"
  "os/exec"
  // "bytes"
  "io/ioutil"
)

var (
  temperatureRegex = regexp.MustCompile(`^temp=(?P<temperature>[0-9]*(?:\.[0-9]*)?)'C$`)
)

func MeasureGPUTemperature() float64 {
  out, err  := exec.Command("vcgencmd", "measure_temp").Output()
  if err != nil {
    log.Fatal(err)
  }

  res := temperatureRegex.FindStringSubmatch(string(out[:len(out)-1]))
  value, err := strconv.ParseFloat(res[1], 64)

  if err != nil {
    log.Fatal(err)
  }
  return value
}

func MeasureCPUTemperature() float64 {
  file, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp") // For read access.
  if err != nil {
  	log.Fatal(err)
  }
  temp, _ := strconv.ParseFloat(string(file[:len(file)-1]), 64)
  return temp/1000
}
