// Package main implements a gateway bewteen a Prometheus Alertmanager and Pushover.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gregdel/pushover"
	amt "github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiKey       = flag.String("api_key", "", "Pushover API key")
	recipient    = flag.String("recipient", "", "Pushover message recipient")
	port         = flag.Int("port", 5001, "Port to listen for alerts")
	dry          = flag.Bool("dry", false, "Dry run only, dont send to PB servers")
	bodyTemplate = template.Must(template.New("body").Parse(`{{.Annotations.summary}}

Labels:
{{ range .Labels.SortedPairs }} - {{.Name }} = {{ .Value }}
{{ end }}
Since: {{ .StartsAt.Format "02 Jan 06 15:04 MST" }}

Link: {{ .GeneratorURL }}
`))
)

var (
	pushoverTotal     prometheus.Gauge
	pushoverRemaining prometheus.Gauge
)

func main() {
	flag.Parse()
	if *apiKey == "" {
		panic("Missing 'api_key' flag")
	}
	pushoverTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pushover_monthly_limit",
		Help: "Maximum messages that can be sent per month to Pushover",
	})
	prometheus.MustRegister(pushoverTotal)
	pushoverRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pushover_monthly_remaining",
		Help: "Remaining messages that can be sent this month to Pushover",
	})
	prometheus.MustRegister(pushoverRemaining)

	app := pushover.New(*apiKey)
	fmt.Printf("ready.\n")
	http.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		var m amt.Data
		if err := dec.Decode(&m); err != nil {
			panic(err)
		}
		rec := pushover.NewRecipient(*recipient)
		for _, alert := range m.Alerts {
			title := fmt.Sprintf("[%s] %s (%s)", strings.ToUpper(alert.Status), alert.Labels["alertname"], alert.Labels["location"])
			var body bytes.Buffer
			if err := bodyTemplate.Execute(&body, alert); err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}
			fmt.Printf("Title: %s\nBody: %s\n", title, body.String())
			if *dry {
				fmt.Printf("Dry run, not sending.")
				continue
			}
			msg := pushover.NewMessageWithTitle(body.String(), title)
			resp, err := app.SendMessage(msg, rec)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}
			pushoverTotal.Set(float64(resp.Limit.Total))
			pushoverRemaining.Set(float64(resp.Limit.Remaining))
			fmt.Println(resp)
		}
	})
	http.Handle("/metrics", promhttp.Handler())
	panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
