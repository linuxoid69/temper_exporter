package main

import (
	"log"
	"fmt"
	"time"
	"github.com/zserge/hid"
	"flag"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var (
 tf,to,hf,ho float64
 port int64
 temperature = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "temperx_temperature",
	Help: "The temperature",
 })
 humidity = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "temperx_humidity",
	Help: "The humidity",
 })
)

func init(){
	flag.Float64Var(&tf,"tf", 1.0, "-tf Factor for temperature")
	flag.Float64Var(&to,"to", 0, "-to Offset for temperature")
	flag.Float64Var(&hf,"hf", 1.0, "-hf Factor for humidity")
	flag.Float64Var(&ho,"ho", 0, "-ho Offset for humidity")
	flag.Int64Var(&port,"port", 9112, "-p Listen a port")
}
func main() {
	flag.Parse()
	log.Printf("Run temper exporter on http://0.0.0.0:%v/metrics ", port)
	getMetrics()
	http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func getMetrics() {

	go func ()  {
		for {
			hid_path := "413d:2107:0000:01"
			cmd_raw := []byte{0x01, 0x80, 0x33, 0x01, 0x00, 0x00, 0x00, 0x00}


			hid.UsbWalk(func(device hid.Device) {
				info := device.Info()

				id := fmt.Sprintf("%04x:%04x:%04x:%02x", info.Vendor, info.Product, info.Revision, info.Interface)
				if id != hid_path {
					temperature.Set(0)
					humidity.Set(0)
					return
				}

				if err := device.Open(); err != nil {
					log.Println("Open error: ", err)
					return
				}

				defer device.Close()

				if _, err := device.Write(cmd_raw, 1*time.Second); err != nil {
					log.Println("Output report write failed:", err)
					return
				}

				if buf, err := device.Read(-1, 1*time.Second); err == nil {
					tmp := (float64(buf[2])*256+float64(buf[3]))/100*tf + to
					hum := (float64(buf[4])*256+float64(buf[5]))/100*hf + ho
					temperature.Set(tmp)
					humidity.Set(hum)
				}
			})
			time.Sleep(1 * time.Second)
		}
}()
}
