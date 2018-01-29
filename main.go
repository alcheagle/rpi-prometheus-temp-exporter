package main

import (
  log "github.com/sirupsen/logrus"
  "os"
  "regexp"
  "strconv"
  "github.com/urfave/cli"

  "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

  "net/http"
  "os/exec"
  "bytes"
)

var(
  temperatureRegex = regexp.MustCompile(`^temp=([0-9]*\.[0-9]*)'C$`)
  temperatureDesc = prometheus.NewDesc(
		"vcgencmd_temperature_celsius",
		"temperature in celsius",
		nil,
    nil)
)

type VcgencmdCollector struct {

}

func NewVcgencmdCollector() (*VcgencmdCollector) {
  return &VcgencmdCollector{}
}

func (l *VcgencmdCollector) Describe(ch chan<- *prometheus.Desc) {
  ch <- temperatureDesc
}

func (me *VcgencmdCollector)Collect(ch chan<- prometheus.Metric) {
  cmd  := exec.Command("./vcgencmd", "measure_temp")
  cmdOutput := &bytes.Buffer{}
  cmd.Stdout = cmdOutput

  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }
  out := cmdOutput.Bytes()
  log.Debug(out)

  value, err := strconv.ParseFloat(string(temperatureRegex.Find(out)), 64)

  if err != nil {
    log.Fatal(err)
  }

  log.Debug(value)

  metric, err := prometheus.NewConstMetric(
    temperatureDesc,
		prometheus.GaugeValue,
		value)

  if err == nil {
    ch <- metric
  } else {
    log.Fatal(err)
  }
}

func main() {
  // temperature := prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "cpu_temperature_celsius",
	// 	Help: "Current temperature of the CPU.",
	// })

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
    // me := []string{"measure_temp"}
    http.Handle("/metrics", promhttp.Handler())
    // http.HandleFunc("/metrics", mhandler.webHandler)
    temperatureCollector := NewVcgencmdCollector()

    prometheus.MustRegister(temperatureCollector)

    http_server := c.String("http-server")
    log.Infof("listening on: [%s]\n", http_server)
    log.Fatal(http.ListenAndServe(http_server, nil))
    return nil
  }

  app.Run(os.Args)
}
