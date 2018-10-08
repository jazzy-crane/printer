// +build windows,!manual

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jazzy-crane/printer"
)

var jobNotify = []uint16{
	printer.JOB_NOTIFY_FIELD_PRINTER_NAME,
	printer.JOB_NOTIFY_FIELD_MACHINE_NAME,
	printer.JOB_NOTIFY_FIELD_PORT_NAME,
	printer.JOB_NOTIFY_FIELD_USER_NAME,
	printer.JOB_NOTIFY_FIELD_STATUS,
}

func main() {
	pnames, err := printer.ReadNames()
	if err != nil {
		log.Println("Error", err)
		os.Exit(1)
	}

	multiplexed := make(chan *printer.PrinterNotifyInfo)
	done := make(chan struct{})

	notifyOptions := &printer.PRINTER_NOTIFY_OPTIONS{
		Version: 2,
		Flags:   0,
		Count:   1,
		PTypes: &printer.PRINTER_NOTIFY_OPTIONS_TYPE{
			Type:    uint16(printer.JOB_NOTIFY_TYPE),
			Count:   uint32(len(jobNotify)),
			PFields: &jobNotify[0],
		},
	}

	for _, pname := range pnames {
		fmt.Println("Opening printer", pname)
		p, err := printer.Open(pname)
		if err != nil {
			log.Println("Error printer.Open", pname, err)
			os.Exit(1)
		}

		n, err := p.GetPrinterNotifications(done, printer.PRINTER_CHANGE_ALL, 0, notifyOptions)
		if err != nil {
			log.Println("Error printer::GetPrinterNotifications", pname, err)
			os.Exit(1)
		}

		go func(notifications <-chan *printer.PrinterNotifyInfo) {
			for n := range notifications {
				multiplexed <- n
			}
			fmt.Println("Assume cleanup complete")
		}(n)
	}

	timeout := time.After(time.Minute)

loop:
	for {
		select {
		case <-timeout:
			close(done)
			time.Sleep(time.Second)
			break loop
		case notification := <-multiplexed:
			fmt.Printf("\n%s\n", notification)
		}
	}
}
