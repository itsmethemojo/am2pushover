Prometheus Alertmanager -> Pushover gateway
=====

A simple gateway that takes alert messages from [Prometheus Alertmanager](https://prometheus.io/docs/alerting/alertmanager/)
via the [webhook receiver](https://prometheus.io/docs/alerting/configuration/#webhook_config),
and forwards them to [Pushover](http://pushover.net/).

No templating or fancy features for now, but these may be included in due course,
perhaps also export metrics for alerts sent (and maybe Pushover's token allowance
left).

Usage
===

Configure as a webhook endpoint in alertmanager.yml:

```
- name: 'default-receiver'
  webhook_configs:
  - url: http://127.0.0.1:5001/alert
```

Set the host/port in the URL to match where this gateway is running.

Build the gateway:

```
$ go get github.com/prometheus/alertmanager/template
$ go get github.com/gregdel/pushover
$ go build am2pushover.go
```

Run the gateway with at least the API key flag, and user (recipient) flag.

```
$ ./am2pushover -api_key o.sdf923456fs765dfsfsdf -recipient oiusdf0u30fu89fu902fs
```

That's it!

Do note that Pushover limits pushes from regular accounts to 7500/month.

Licence
===

MIT.

