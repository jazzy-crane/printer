// +build windows,manual

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jazzy-crane/printer"
)

func main() {
	pnames, err := printer.ReadNames()
	if err != nil {
		log.Println("Error", err)
		os.Exit(1)
	}

	multiplexed := make(chan *printer.NotifyInfo)

	notifyOptions := &printer.PRINTER_NOTIFY_OPTIONS{
		Version: 2,
		Flags:   0,
		Count:   1,
		PTypes: &printer.PRINTER_NOTIFY_OPTIONS_TYPE{
			Type:    uint16(printer.JOB_NOTIFY_TYPE),
			Count:   uint32(len(printer.JobNotifyAll)),
			PFields: &printer.JobNotifyAll[0],
		},
	}

	for _, pname := range pnames {
		p, err := printer.Open(pname)
		if err != nil {
			log.Println("Error printer.Open", pname, err)
			os.Exit(1)
		}

		notifications, err := p.ChangeNotifications(printer.PRINTER_CHANGE_ALL, 0, notifyOptions)
		if err != nil {
			log.Println("Error ChangeNotifications", err)
			os.Exit(1)
		}

		go func(p *printer.Printer, pcnh *printer.ChangeNotificationHandle) {
			defer pcnh.Close()
			defer p.Close()

			for {
				pni, err := pcnh.Next(nil)
				if err != nil {
					if err != printer.ErrNoNotification {
						log.Println("Unexpected error from ChangeNotificationHandle::Next", err)
					}
					continue
				}
				multiplexed <- pni
			}
		}(p, notifications)
	}

	for notification := range multiplexed {
		fmt.Printf("\n%s\n", notification)
	}
}
