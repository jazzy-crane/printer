// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Windows printing.
package printer

import (
	"errors"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

//go:generate go run mksyscall_windows.go -output zapi.go printer.go

type DOC_INFO_1 struct {
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
}

type PRINTER_INFO_5 struct {
	PrinterName              *uint16
	PortName                 *uint16
	Attributes               uint32
	DeviceNotSelectedTimeout uint32
	TransmissionRetryTimeout uint32
}

type DRIVER_INFO_8 struct {
	Version                  uint32
	Name                     *uint16
	Environment              *uint16
	DriverPath               *uint16
	DataFile                 *uint16
	ConfigFile               *uint16
	HelpFile                 *uint16
	DependentFiles           *uint16
	MonitorName              *uint16
	DefaultDataType          *uint16
	PreviousNames            *uint16
	DriverDate               syscall.Filetime
	DriverVersion            uint64
	MfgName                  *uint16
	OEMUrl                   *uint16
	HardwareID               *uint16
	Provider                 *uint16
	PrintProcessor           *uint16
	VendorSetup              *uint16
	ColorProfiles            *uint16
	InfPath                  *uint16
	PrinterDriverAttributes  uint32
	CoreDriverDependencies   *uint16
	MinInboxDriverVerDate    syscall.Filetime
	MinInboxDriverVerVersion uint32
}

type JOB_INFO_1 struct {
	JobID        uint32
	PrinterName  *uint16
	MachineName  *uint16
	UserName     *uint16
	Document     *uint16
	DataType     *uint16
	Status       *uint16
	StatusCode   uint32
	Priority     uint32
	Position     uint32
	TotalPages   uint32
	PagesPrinted uint32
	Submitted    syscall.Systemtime
}

type PRINTER_NOTIFY_OPTIONS_TYPE struct {
	Type      uint16
	Reserved0 uint16
	Reserved1 uint32
	Reserved2 uint32
	Count     uint32
	PFields   *uint16
}

type PRINTER_NOTIFY_OPTIONS struct {
	Version uint32
	Flags   uint32
	Count   uint32
	PTypes  *PRINTER_NOTIFY_OPTIONS_TYPE
}

type NOTIFY_INFO struct {
	Datasz  uint32
	Dataptr unsafe.Pointer
}

type PRINTER_NOTIFY_INFO_DATA struct {
	Type       uint16
	Field      uint16
	Reserved   uint32
	ID         uint32
	NotifyInfo NOTIFY_INFO
}

type PRINTER_NOTIFY_INFO struct {
	Version uint32
	Flags   uint32
	Count   uint32
	PData   [0xff]PRINTER_NOTIFY_INFO_DATA
}

const (
	PRINTER_ENUM_LOCAL       = 2
	PRINTER_ENUM_CONNECTIONS = 4

	PRINTER_DRIVER_XPS = 0x00000002

	PRINTER_CHANGE_ADD_PRINTER               = 0x00000001
	PRINTER_CHANGE_SET_PRINTER               = 0x00000002
	PRINTER_CHANGE_DELETE_PRINTER            = 0x00000004
	PRINTER_CHANGE_FAILED_CONNECTION_PRINTER = 0x00000008
	PRINTER_CHANGE_PRINTER                   = 0x000000FF
	PRINTER_CHANGE_ADD_JOB                   = 0x00000100
	PRINTER_CHANGE_SET_JOB                   = 0x00000200
	PRINTER_CHANGE_DELETE_JOB                = 0x00000400
	PRINTER_CHANGE_WRITE_JOB                 = 0x00000800
	PRINTER_CHANGE_JOB                       = 0x0000FF00
	PRINTER_CHANGE_ADD_FORM                  = 0x00010000
	PRINTER_CHANGE_SET_FORM                  = 0x00020000
	PRINTER_CHANGE_DELETE_FORM               = 0x00040000
	PRINTER_CHANGE_FORM                      = 0x00070000
	PRINTER_CHANGE_ADD_PORT                  = 0x00100000
	PRINTER_CHANGE_CONFIGURE_PORT            = 0x00200000
	PRINTER_CHANGE_DELETE_PORT               = 0x00400000
	PRINTER_CHANGE_PORT                      = 0x00700000
	PRINTER_CHANGE_ADD_PRINT_PROCESSOR       = 0x01000000
	PRINTER_CHANGE_DELETE_PRINT_PROCESSOR    = 0x04000000
	PRINTER_CHANGE_PRINT_PROCESSOR           = 0x07000000
	PRINTER_CHANGE_SERVER                    = 0x08000000
	PRINTER_CHANGE_ADD_PRINTER_DRIVER        = 0x10000000
	PRINTER_CHANGE_SET_PRINTER_DRIVER        = 0x20000000
	PRINTER_CHANGE_DELETE_PRINTER_DRIVER     = 0x40000000
	PRINTER_CHANGE_PRINTER_DRIVER            = 0x70000000
	PRINTER_CHANGE_TIMEOUT                   = 0x80000000
	PRINTER_CHANGE_ALL                       = 0x7F77FFFF

	JOB_NOTIFY_FIELD_PRINTER_NAME        = 0x00
	JOB_NOTIFY_FIELD_MACHINE_NAME        = 0x01
	JOB_NOTIFY_FIELD_PORT_NAME           = 0x02
	JOB_NOTIFY_FIELD_USER_NAME           = 0x03
	JOB_NOTIFY_FIELD_NOTIFY_NAME         = 0x04
	JOB_NOTIFY_FIELD_DATATYPE            = 0x05
	JOB_NOTIFY_FIELD_PRINT_PROCESSOR     = 0x06
	JOB_NOTIFY_FIELD_PARAMETERS          = 0x07
	JOB_NOTIFY_FIELD_DRIVER_NAME         = 0x08
	JOB_NOTIFY_FIELD_DEVMODE             = 0x09
	JOB_NOTIFY_FIELD_STATUS              = 0x0A
	JOB_NOTIFY_FIELD_STATUS_STRING       = 0x0B
	JOB_NOTIFY_FIELD_SECURITY_DESCRIPTOR = 0x0C
	JOB_NOTIFY_FIELD_DOCUMENT            = 0x0D
	JOB_NOTIFY_FIELD_PRIORITY            = 0x0E
	JOB_NOTIFY_FIELD_POSITION            = 0x0F
	JOB_NOTIFY_FIELD_SUBMITTED           = 0x10
	JOB_NOTIFY_FIELD_START_TIME          = 0x11
	JOB_NOTIFY_FIELD_UNTIL_TIME          = 0x12
	JOB_NOTIFY_FIELD_TIME                = 0x13
	JOB_NOTIFY_FIELD_TOTAL_PAGES         = 0x14
	JOB_NOTIFY_FIELD_PAGES_PRINTED       = 0x15
	JOB_NOTIFY_FIELD_TOTAL_BYTES         = 0x16
	JOB_NOTIFY_FIELD_BYTES_PRINTED       = 0x17
	JOB_NOTIFY_FIELD_REMOTE_JOB_ID       = 0x18

	PRINTER_NOTIFY_TYPE = 0 // TODO: Implement support for this
	JOB_NOTIFY_TYPE     = 1

	PRINTER_NOTIFY_INFO_DISCARDED  = 1
	PRINTER_NOTIFY_OPTIONS_REFRESH = 1

	JOB_STATUS_PAUSED            = 0x00000001 // Job is paused
	JOB_STATUS_ERROR             = 0x00000002 // An error is associated with the job
	JOB_STATUS_DELETING          = 0x00000004 // Job is being deleted
	JOB_STATUS_SPOOLING          = 0x00000008 // Job is spooling
	JOB_STATUS_PRINTING          = 0x00000010 // Job is printing
	JOB_STATUS_OFFLINE           = 0x00000020 // Printer is offline
	JOB_STATUS_PAPEROUT          = 0x00000040 // Printer is out of paper
	JOB_STATUS_PRINTED           = 0x00000080 // Job has printed
	JOB_STATUS_DELETED           = 0x00000100 // Job has been deleted
	JOB_STATUS_BLOCKED_DEVQ      = 0x00000200 // Printer driver cannot print the job
	JOB_STATUS_USER_INTERVENTION = 0x00000400 // User action required
	JOB_STATUS_RESTART           = 0x00000800 // Job has been restarted
	JOB_STATUS_COMPLETE          = 0x00001000 // Job has been delivered to the printer
	JOB_STATUS_RETAINED          = 0x00002000 // Job has been retained in the print queue
	JOB_STATUS_RENDERING_LOCALLY = 0x00004000 // Job rendering locally on the client
)

var ErrNoNotification = errors.New("no notification information")

//sys	GetDefaultPrinter(buf *uint16, bufN *uint32) (err error) = winspool.GetDefaultPrinterW
//sys	ClosePrinter(h syscall.Handle) (err error) = winspool.ClosePrinter
//sys	OpenPrinter(name *uint16, h *syscall.Handle, defaults uintptr) (err error) = winspool.OpenPrinterW
//sys	StartDocPrinter(h syscall.Handle, level uint32, docinfo *DOC_INFO_1) (err error) = winspool.StartDocPrinterW
//sys	EndDocPrinter(h syscall.Handle) (err error) = winspool.EndDocPrinter
//sys	WritePrinter(h syscall.Handle, buf *byte, bufN uint32, written *uint32) (err error) = winspool.WritePrinter
//sys	StartPagePrinter(h syscall.Handle) (err error) = winspool.StartPagePrinter
//sys	EndPagePrinter(h syscall.Handle) (err error) = winspool.EndPagePrinter
//sys	EnumPrinters(flags uint32, name *uint16, level uint32, buf *byte, bufN uint32, needed *uint32, returned *uint32) (err error) = winspool.EnumPrintersW
//sys	GetPrinterDriver(h syscall.Handle, env *uint16, level uint32, di *byte, n uint32, needed *uint32) (err error) = winspool.GetPrinterDriverW
//sys	EnumJobs(h syscall.Handle, firstJob uint32, noJobs uint32, level uint32, buf *byte, bufN uint32, bytesNeeded *uint32, jobsReturned *uint32) (err error) = winspool.EnumJobsW
//sys   FindFirstPrinterChangeNotification(h syscall.Handle, filter uint32, options uint32, notifyOptions *PRINTER_NOTIFY_OPTIONS) (rtn syscall.Handle, err error) = winspool.FindFirstPrinterChangeNotification
//sys   FindNextPrinterChangeNotification(h syscall.Handle, cause *uint16, options *PRINTER_NOTIFY_OPTIONS, info **PRINTER_NOTIFY_INFO) (err error) = winspool.FindNextPrinterChangeNotification
//sys   FindClosePrinterChangeNotification(h syscall.Handle) (err error) = winspool.FindClosePrinterChangeNotification
//sys   FreePrinterNotifyInfo(info *PRINTER_NOTIFY_INFO) (err error) = winspool.FreePrinterNotifyInfo

func Default() (string, error) {
	b := make([]uint16, 3)
	n := uint32(len(b))
	err := GetDefaultPrinter(&b[0], &n)
	if err != nil {
		if err != syscall.ERROR_INSUFFICIENT_BUFFER {
			return "", err
		}
		b = make([]uint16, n)
		err = GetDefaultPrinter(&b[0], &n)
		if err != nil {
			return "", err
		}
	}
	return syscall.UTF16ToString(b), nil
}

// ReadNames return printer names on the system
func ReadNames() ([]string, error) {
	const flags = PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS
	var needed, returned uint32
	buf := make([]byte, 1)
	err := EnumPrinters(flags, nil, 5, &buf[0], uint32(len(buf)), &needed, &returned)
	if err != nil {
		if err != syscall.ERROR_INSUFFICIENT_BUFFER {
			return nil, err
		}
		buf = make([]byte, needed)
		err = EnumPrinters(flags, nil, 5, &buf[0], uint32(len(buf)), &needed, &returned)
		if err != nil {
			return nil, err
		}
	}
	ps := (*[1024]PRINTER_INFO_5)(unsafe.Pointer(&buf[0]))[:returned]
	names := make([]string, 0, returned)
	for _, p := range ps {
		v := (*[1024]uint16)(unsafe.Pointer(p.PrinterName))[:]
		names = append(names, syscall.UTF16ToString(v))
	}
	return names, nil
}

type Printer struct {
	h                    syscall.Handle
	notificationsRunning bool
}

func Open(name string) (*Printer, error) {
	var p Printer
	// TODO: implement pDefault parameter
	nameutf16, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}

	err = OpenPrinter(nameutf16, &p.h, 0)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// DriverInfo stores information about printer driver.
type DriverInfo struct {
	Name        string
	Environment string
	DriverPath  string
	Attributes  uint32
}

// JobInfo stores information about a print job.
type JobInfo struct {
	JobID           uint32
	UserMachineName string
	UserName        string
	DocumentName    string
	DataType        string
	Status          string
	StatusCode      uint32
	Priority        uint32
	Position        uint32
	TotalPages      uint32
	PagesPrinted    uint32
	Submitted       time.Time
}

// Jobs returns information about all print jobs on this printer
func (p *Printer) Jobs() ([]JobInfo, error) {
	var bytesNeeded, jobsReturned uint32
	buf := make([]byte, 1)
	for {
		err := EnumJobs(p.h, 0, 255, 1, &buf[0], uint32(len(buf)), &bytesNeeded, &jobsReturned)
		if err == nil {
			break
		}
		if err != syscall.ERROR_INSUFFICIENT_BUFFER {
			return nil, err
		}
		if bytesNeeded <= uint32(len(buf)) {
			return nil, err
		}
		buf = make([]byte, bytesNeeded)
	}
	if jobsReturned <= 0 {
		return nil, nil
	}
	pjs := make([]JobInfo, 0, jobsReturned)
	ji := (*[2048]JOB_INFO_1)(unsafe.Pointer(&buf[0]))[:jobsReturned]
	for _, j := range ji {
		pji := JobInfo{
			JobID:        j.JobID,
			StatusCode:   j.StatusCode,
			Priority:     j.Priority,
			Position:     j.Position,
			TotalPages:   j.TotalPages,
			PagesPrinted: j.PagesPrinted,
		}
		if j.MachineName != nil {
			pji.UserMachineName = syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(j.MachineName))[:])
		}
		if j.UserName != nil {
			pji.UserName = syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(j.UserName))[:])
		}
		if j.Document != nil {
			pji.DocumentName = syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(j.Document))[:])
		}
		if j.DataType != nil {
			pji.DataType = syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(j.DataType))[:])
		}
		if j.Status != nil {
			pji.Status = syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(j.Status))[:])
		}
		if strings.TrimSpace(pji.Status) == "" {
			if pji.StatusCode == 0 {
				pji.Status += "Queue Paused, "
			}
			if pji.StatusCode&JOB_STATUS_PRINTING != 0 {
				pji.Status += "Printing, "
			}
			if pji.StatusCode&JOB_STATUS_PAUSED != 0 {
				pji.Status += "Paused, "
			}
			if pji.StatusCode&JOB_STATUS_ERROR != 0 {
				pji.Status += "Error, "
			}
			if pji.StatusCode&JOB_STATUS_DELETING != 0 {
				pji.Status += "Deleting, "
			}
			if pji.StatusCode&JOB_STATUS_SPOOLING != 0 {
				pji.Status += "Spooling, "
			}
			if pji.StatusCode&JOB_STATUS_OFFLINE != 0 {
				pji.Status += "Printer Offline, "
			}
			if pji.StatusCode&JOB_STATUS_PAPEROUT != 0 {
				pji.Status += "Out of Paper, "
			}
			if pji.StatusCode&JOB_STATUS_PRINTED != 0 {
				pji.Status += "Printed, "
			}
			if pji.StatusCode&JOB_STATUS_DELETED != 0 {
				pji.Status += "Deleted, "
			}
			if pji.StatusCode&JOB_STATUS_BLOCKED_DEVQ != 0 {
				pji.Status += "Driver Error, "
			}
			if pji.StatusCode&JOB_STATUS_USER_INTERVENTION != 0 {
				pji.Status += "User Action Required, "
			}
			if pji.StatusCode&JOB_STATUS_RESTART != 0 {
				pji.Status += "Restarted, "
			}
			if pji.StatusCode&JOB_STATUS_COMPLETE != 0 {
				pji.Status += "Sent to Printer, "
			}
			if pji.StatusCode&JOB_STATUS_RETAINED != 0 {
				pji.Status += "Retained, "
			}
			if pji.StatusCode&JOB_STATUS_RENDERING_LOCALLY != 0 {
				pji.Status += "Rendering on Client, "
			}
			pji.Status = strings.TrimRight(pji.Status, ", ")
		}
		pji.Submitted = time.Date(
			int(j.Submitted.Year),
			time.Month(int(j.Submitted.Month)),
			int(j.Submitted.Day),
			int(j.Submitted.Hour),
			int(j.Submitted.Minute),
			int(j.Submitted.Second),
			int(1000*j.Submitted.Milliseconds),
			time.Local,
		).UTC()
		pjs = append(pjs, pji)
	}
	return pjs, nil
}

// DriverInfo returns information about printer p driver.
func (p *Printer) DriverInfo() (*DriverInfo, error) {
	var needed uint32
	b := make([]byte, 1024*10)
	for {
		err := GetPrinterDriver(p.h, nil, 8, &b[0], uint32(len(b)), &needed)
		if err == nil {
			break
		}
		if err != syscall.ERROR_INSUFFICIENT_BUFFER {
			return nil, err
		}
		if needed <= uint32(len(b)) {
			return nil, err
		}
		b = make([]byte, needed)
	}
	di := (*DRIVER_INFO_8)(unsafe.Pointer(&b[0]))
	return &DriverInfo{
		Attributes:  di.PrinterDriverAttributes,
		Name:        syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(di.Name))[:]),
		DriverPath:  syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(di.DriverPath))[:]),
		Environment: syscall.UTF16ToString((*[2048]uint16)(unsafe.Pointer(di.Environment))[:]),
	}, nil
}

func (p *Printer) StartDocument(name, datatype string) error {
	d := DOC_INFO_1{
		DocName:    &(syscall.StringToUTF16(name))[0],
		OutputFile: nil,
		Datatype:   &(syscall.StringToUTF16(datatype))[0],
	}
	return StartDocPrinter(p.h, 1, &d)
}

// StartRawDocument calls StartDocument and passes either "RAW" or "XPS_PASS"
// as a document type, depending if printer driver is XPS-based or not.
func (p *Printer) StartRawDocument(name string) error {
	di, err := p.DriverInfo()
	if err != nil {
		return err
	}
	// See https://support.microsoft.com/en-us/help/2779300/v4-print-drivers-using-raw-mode-to-send-pcl-postscript-directly-to-the
	// for details.
	datatype := "RAW"
	if di.Attributes&PRINTER_DRIVER_XPS != 0 {
		datatype = "XPS_PASS"
	}
	return p.StartDocument(name, datatype)
}

func (p *Printer) Write(b []byte) (int, error) {
	var written uint32
	err := WritePrinter(p.h, &b[0], uint32(len(b)), &written)
	if err != nil {
		return 0, err
	}
	return int(written), nil
}

func (p *Printer) EndDocument() error {
	return EndDocPrinter(p.h)
}

func (p *Printer) StartPage() error {
	return StartPagePrinter(p.h)
}

func (p *Printer) EndPage() error {
	return EndPagePrinter(p.h)
}

func (p *Printer) Close() error {
	return ClosePrinter(p.h)
}

type NotifyInfoData struct {
	Type  uint16 // one of PRINTER_NOTIFY_TYPE or JOB_NOTIFY_TYPE
	Field uint16 // JOB_NOTIFY_FIELD_* or PRINTER_NOTIFY_FIELD_* depending on the above
	ID    uint32 // if JOB_NOTIFY_TYPE, this is the print job ID

	Value interface{}
}

type NotifyInfo struct {
	Version int
	Flags   uint
	Cause   uint
	Data    []*NotifyInfoData
}

func (pnid *PRINTER_NOTIFY_INFO_DATA) ToNotifyInfoData() *NotifyInfoData {
	p := &NotifyInfoData{
		Type:  pnid.Type,
		Field: pnid.Field,
		ID:    pnid.ID,
	}

	if pnid.Type == JOB_NOTIFY_TYPE {
		switch pnid.Field {
		case JOB_NOTIFY_FIELD_PRINTER_NAME,
			JOB_NOTIFY_FIELD_MACHINE_NAME,
			JOB_NOTIFY_FIELD_PORT_NAME,
			JOB_NOTIFY_FIELD_USER_NAME,
			JOB_NOTIFY_FIELD_NOTIFY_NAME,
			JOB_NOTIFY_FIELD_DATATYPE,
			JOB_NOTIFY_FIELD_PRINT_PROCESSOR,
			JOB_NOTIFY_FIELD_PARAMETERS,
			JOB_NOTIFY_FIELD_DRIVER_NAME,
			JOB_NOTIFY_FIELD_STATUS_STRING,
			JOB_NOTIFY_FIELD_DOCUMENT:
			ps := (*[0xffff]uint16)(pnid.NotifyInfo.Dataptr)[:pnid.NotifyInfo.Datasz/2]
			p.Value = syscall.UTF16ToString(ps)
		case JOB_NOTIFY_FIELD_STATUS,
			JOB_NOTIFY_FIELD_PRIORITY,
			JOB_NOTIFY_FIELD_POSITION,
			JOB_NOTIFY_FIELD_START_TIME,
			JOB_NOTIFY_FIELD_UNTIL_TIME,
			JOB_NOTIFY_FIELD_TIME,
			JOB_NOTIFY_FIELD_TOTAL_PAGES,
			JOB_NOTIFY_FIELD_PAGES_PRINTED,
			JOB_NOTIFY_FIELD_TOTAL_BYTES,
			JOB_NOTIFY_FIELD_BYTES_PRINTED:
			p.Value = int(pnid.NotifyInfo.Datasz)
		case JOB_NOTIFY_FIELD_DEVMODE:
			// TODO pnid.NotifyInfo.Dataptr is a pointer to a DEVMODE structure that contains device-initialization and environment data for the printer driver.
		case JOB_NOTIFY_FIELD_SECURITY_DESCRIPTOR:
			// TODO Not supported according to https://docs.microsoft.com/en-us/windows/desktop/printdocs/printer-notify-info-data , though does seem to be something there
		case JOB_NOTIFY_FIELD_SUBMITTED:
			// TODO pnid.NotifyInfo.Dataptr is a pointer to a SYSTEMTIME structure that specifies the time when the job was submitted.
		default:
		}
	}

	return p
}

func (pni *PRINTER_NOTIFY_INFO) ToNotifyInfo() *NotifyInfo {
	p := &NotifyInfo{
		Version: int(pni.Version),
		Flags:   uint(pni.Flags),
		Data:    make([]*NotifyInfoData, pni.Count),
	}

	for i := 0; i < int(pni.Count); i++ {
		p.Data[i] = pni.PData[i].ToNotifyInfoData()
	}

	return p
}

type ChangeNotificationHandle struct {
	h syscall.Handle
}

func (p *Printer) ChangeNotifications(filter uint32, options uint32, printerNotifyOptions *PRINTER_NOTIFY_OPTIONS) (*ChangeNotificationHandle, error) {
	h, err := FindFirstPrinterChangeNotification(p.h, filter, options, printerNotifyOptions)
	if err != nil {
		return nil, err
	}

	return &ChangeNotificationHandle{
		h: h,
	}, nil
}

func (c *ChangeNotificationHandle) Next(printerNotifyOptions *PRINTER_NOTIFY_OPTIONS) (*NotifyInfo, error) {
	var cause uint16
	var notifyInfo *PRINTER_NOTIFY_INFO

	err := FindNextPrinterChangeNotification(c.h, &cause, nil, &notifyInfo)
	if err != nil {
		return nil, err
	}

	if notifyInfo != nil && (notifyInfo.Flags&PRINTER_NOTIFY_INFO_DISCARDED) == PRINTER_NOTIFY_INFO_DISCARDED {
		/* If the PRINTER_NOTIFY_INFO_DISCARDED bit is set in the Flags member of the PRINTER_NOTIFY_INFO structure,
		an overflow or error occurred, and notifications may have been lost.
		In this case, no additional notifications will be sent until you make a second
		FindNextPrinterChangeNotification call that specifies PRINTER_NOTIFY_OPTIONS_REFRESH.
		*/
		_ = FreePrinterNotifyInfo(notifyInfo)
		notifyInfo = nil

		pno := &PRINTER_NOTIFY_OPTIONS{
			Version: 2,
			Flags:   PRINTER_NOTIFY_OPTIONS_REFRESH,
			Count:   0,
			PTypes:  nil,
		}

		err := FindNextPrinterChangeNotification(c.h, &cause, pno, &notifyInfo)
		if err != nil {
			return nil, err
		}
	}

	if notifyInfo != nil {
		pni := notifyInfo.ToNotifyInfo()
		pni.Cause = uint(cause)

		_ = FreePrinterNotifyInfo(notifyInfo)

		return pni, nil
	} else {
		return nil, ErrNoNotification
	}
}

func (c *ChangeNotificationHandle) Wait(milliseconds uint32) (uint32, error) {
	return syscall.WaitForSingleObject(c.h, milliseconds)
}

func (c *ChangeNotificationHandle) Close() error {
	return FindClosePrinterChangeNotification(c.h)
}
