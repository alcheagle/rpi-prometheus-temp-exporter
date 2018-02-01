package main

import (
  log "github.com/sirupsen/logrus"
  "os"

  "github.com/urfave/cli"

  "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

  "net/http"

  "github.com/alcheagle/rpi-prometheus-temp-exporter/temperatureSensors"
)

var(
  // merge the metrics with different labels
  gpuTemperatureDesc = prometheus.NewDesc(
		"rpi_gpu_temperature_celsius",
		"temperature in celsius",
		nil,
    nil)

  cpuTemperatureDesc = prometheus.NewDesc(
		"rpi_cpu_temperature_celsius",
		"temperature in celsius",
		nil,
    nil)
)

type TemperatureCollector struct {

}

func NewTemperatureCollector() (*TemperatureCollector) {
  return &TemperatureCollector{}
}

func (l *TemperatureCollector) Describe(ch chan<- *prometheus.Desc) {
  ch <- cpuTemperatureDesc
  ch <- gpuTemperatureDesc
}

func (me *TemperatureCollector)Collect(ch chan<- prometheus.Metric) {
  //measure the temperature

  metric1, err := prometheus.NewConstMetric(
    cpuTemperatureDesc,
    prometheus.GaugeValue,
    TemperatureSensors.MeasureCPUTemperature())

  if err == nil {
    ch <- metric1
  } else {
    log.Fatal(err)
  }

  metric2, err := prometheus.NewConstMetric(
    gpuTemperatureDesc,
    prometheus.GaugeValue,
    TemperatureSensors.MeasureGPUTemperature())

  if err == nil {
    ch <- metric2
  } else {
    log.Fatal(err)
  }
}

func main() {
  app := cli.NewApp()
  app.Usage = "export rpi metrics to docker"

  app.Flags = []cli.Flag {
    cli.StringFlag {
      Name:   "http-server, s",
      Value:  ":9101",
      Usage:  "host:port combination to bind the http service to",
      EnvVar: "HTTP",
    },
    cli.StringFlag {
      Name:   "logging-level, log",
      Value:  "INFO",
      Usage:  "The logging level for this application",
      EnvVar: "LOGGING_LEVEL",
    },
  }

  app.Action = func(c *cli.Context) error {
    logging_level_flag := c.String("logging-level")
    logging_level, err := log.ParseLevel(logging_level_flag)
    if err != nil {
      log.Fatalf("logging level: %s doesn't exist", logging_level_flag)
    }
    log.SetLevel(logging_level)
    
    http.Handle("/metrics", promhttp.Handler())

    temperatureCollector := NewTemperatureCollector()

    prometheus.MustRegister(temperatureCollector)

    http_server := c.String("http-server")
    log.Infof("listening on: [%s]", http_server)
    log.Fatal(http.ListenAndServe(http_server, nil))
    return nil
  }

  app.Run(os.Args)
}
