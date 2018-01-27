package main

import (
  "os"
  "fmt"
  "github.com/urfave/cli"
  "net/http"
  "os/exec"
  "bytes"
)

type MetricsHandler struct {
  metrics []string
}

func (me *MetricsHandler)webHandler(w http.ResponseWriter, r *http.Request) {
  //execute command for getting the temperature
  values := me.ObtainMetrics()
  for _,  value := range values {
    fmt.Println(value)
    fmt.Fprintf(w, value)
  }
}

func (me *MetricsHandler)ObtainMetrics() (map[string]string) {
  results := make(map[string]string)

  for _,  metric := range me.metrics {
    cmd  := exec.Command("vcgencmd", metric)
    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput

    err := cmd.Run()
    if err != nil {
      //TODO an error has occourred
    }
    results[metric] = string(cmdOutput.Bytes())
    fmt.Println(results[metric])
  }
  fmt.Println(me.metrics)
  return results
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
    // cli.StringSliceFlag {
    //   Name:   "metrics, m",
    //   Usage:  "the metrics to expose",
    //   Value:  &cli.StringSlice{"measure_temp"},
    //   EnvVar: "METRICS",
    // },
  }

  app.Action = func(c *cli.Context) error {
    // me := c.StringSlice("metrics")
    me := []string{"measure_temp"}
    fmt.Println(me)
    mhandler := &MetricsHandler{metrics: me}

    http.HandleFunc("/metrics", mhandler.webHandler)

    http_server := c.String("http-server")
    fmt.Printf("listening on: %s\n", http_server)
    err := http.ListenAndServe(http_server, nil)

    fmt.Println(err);
    return nil
  }

  app.Run(os.Args)
}
