package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/icrowley/fake"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
    RegisterURL = "http://host.docker.internal:5555/register"
    MeURL = "http://host.docker.internal:5555/me?curr=rub"
    DepositURL = "http://host.docker.internal:5555/deposit"
    WithdrawURL = "http://host.docker.internal:5555/withdraw"
)

type Metrics struct {
    Users       prometheus.Gauge
    Requests    prometheus.HistogramVec
}

func newMetrics(reg prometheus.Registerer) *Metrics {
    metrics := &Metrics{
        Users: prometheus.NewGauge(prometheus.GaugeOpts{
            Namespace: "client",
            Name: "users_counter",
            Help: "Number of users",
        }),
        Requests: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
            Namespace: "client",
            Name: "requests_latency",
            Help: "Latency of the requests",
            Buckets: []float64{0.05, 0.1, 0.15, 0.2, 0.3, 0.4, 0.5},
        }, []string{"status"}),
    }
    reg.MustRegister(metrics.Users, metrics.Requests)
    return metrics
}

type Sender struct {
    ctx context.Context
    m   *Metrics
    wg  *sync.WaitGroup
}

func main() {
    reg := prometheus.NewRegistry()
    metrics := newMetrics(reg)

    ctx, cancel := context.WithCancel(context.Background())
    wg := &sync.WaitGroup{}
    sender := Sender{
        ctx: ctx,
        m: metrics,
        wg: wg,
    }

    go sender.startSending()

    promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
    http.Handle("/metrics", promHandler)
    http.ListenAndServe(":8081", nil)
    cancel()
}

func (s *Sender) startSending() {
    ticker := time.NewTicker(1 * time.Second)
    done := make(chan bool)
    go func(){
        for {
            select {
                case <-done:
                    break
                case <-ticker.C:
                    s.wg.Add(10)
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
                    go s.processUser()
            }
        }
    }()
    time.Sleep(10 * time.Minute)
    ticker.Stop()
    done <- true
    <-s.ctx.Done()
    s.wg.Wait()
}

func (s *Sender) processUser() {
    token := s.Register()
    s.Me(token)
    s.Deposit(token)
    s.Me(token)
    s.Withdraw(token)
    s.Me(token)
    s.Deposit(token)
    s.Me(token)
    s.Withdraw(token)
    s.Me(token)
    s.wg.Done()
}

func (s *Sender) Register() string {
    start := time.Now()
    postBody, _ := json.Marshal(map[string]string{
        "name":  fake.FullName(),
    })
    responseBody := bytes.NewBuffer(postBody)

    req, _ := http.NewRequest("POST", RegisterURL, responseBody)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)

    if err == nil && resp.StatusCode == http.StatusOK {
        s.m.Users.Add(1)
    }
    
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    result := make(map[string]string)

    json.Unmarshal(body, &result)

    s.m.Requests.With(prometheus.Labels{"status": fmt.Sprintf("%d", resp.StatusCode)}).Observe(time.Since(start).Seconds())

    return result["token"]
}

func (s *Sender) Me(token string) {
    start := time.Now()

    req, _ := http.NewRequest("GET", MeURL, nil)

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", token)

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    s.m.Requests.With(prometheus.Labels{"status": fmt.Sprintf("%d", resp.StatusCode)}).Observe(time.Since(start).Seconds())
}

func (s *Sender) Deposit(token string) {
    start := time.Now()
    postBody, _ := json.Marshal(map[string]string{
        "amount":  "16",
        "currency": "usd",
    })
    responseBody := bytes.NewBuffer(postBody)

    req, _ := http.NewRequest("POST", DepositURL, responseBody)

    req.Header.Set("Authorization", token)

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    s.m.Requests.With(prometheus.Labels{"status": fmt.Sprintf("%d", resp.StatusCode)}).Observe(time.Since(start).Seconds())
}

func (s *Sender) Withdraw(token string) {
    start := time.Now()
    postBody, _ := json.Marshal(map[string]string{
        "amount":  "3",
        "currency": "eur",
    })
    responseBody := bytes.NewBuffer(postBody)

    req, _ := http.NewRequest("POST", WithdrawURL, responseBody)

    req.Header.Set("Authorization", token)

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    s.m.Requests.With(prometheus.Labels{"status": fmt.Sprintf("%d", resp.StatusCode)}).Observe(time.Since(start).Seconds())
}
