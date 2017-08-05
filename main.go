package main

import (
	"log"

	"sync"

	"time"

	"github.com/PagerDuty/godspeed"
	"github.com/mdlayher/apcupsd"
)

func main() {
	g, err := godspeed.NewDefault()
	if err != nil {
		log.Fatalf("Failed to make a datadog client: %v\n", err)
	}

	defer g.Conn.Close()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		for {
			time.Sleep(time.Minute)
			apcClient, err := apcupsd.Dial("tcp", "127.0.0.1:3551")
			if err != nil {
				log.Printf("Failed to dial APC: %v\n", err)
				continue
			}


			status, err := apcClient.Status()
			if err != nil {
				log.Printf("Failed to get APC status: %v\n", err)
				continue
			}

			apcClient.Close()

			err = g.Gauge("ups.load", status.LoadPercent, nil)
			if err != nil {
				log.Printf("Failed to write to datadog: %v\n", err)
			}

			err = g.Gauge("ups.charge", status.BatteryChargePercent, nil)
			if err != nil {
				log.Printf("Failed to write to datadog: %v\n", err)
			}
		}
	}()

	wg.Wait()
}
