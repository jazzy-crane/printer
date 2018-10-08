package printer

import "fmt"

var JobNotifyAll = []uint16{
	JOB_NOTIFY_FIELD_PRINTER_NAME,
	JOB_NOTIFY_FIELD_MACHINE_NAME,
	JOB_NOTIFY_FIELD_PORT_NAME,
	JOB_NOTIFY_FIELD_USER_NAME,
	JOB_NOTIFY_FIELD_NOTIFY_NAME,
	JOB_NOTIFY_FIELD_DATATYPE,
	JOB_NOTIFY_FIELD_PRINT_PROCESSOR,
	JOB_NOTIFY_FIELD_PARAMETERS,
	JOB_NOTIFY_FIELD_DRIVER_NAME,
	JOB_NOTIFY_FIELD_DEVMODE,
	JOB_NOTIFY_FIELD_STATUS,
	JOB_NOTIFY_FIELD_STATUS_STRING,
	JOB_NOTIFY_FIELD_SECURITY_DESCRIPTOR,
	JOB_NOTIFY_FIELD_DOCUMENT,
	JOB_NOTIFY_FIELD_PRIORITY,
	JOB_NOTIFY_FIELD_POSITION,
	JOB_NOTIFY_FIELD_SUBMITTED,
	JOB_NOTIFY_FIELD_START_TIME,
	JOB_NOTIFY_FIELD_UNTIL_TIME,
	JOB_NOTIFY_FIELD_TIME,
	JOB_NOTIFY_FIELD_TOTAL_PAGES,
	JOB_NOTIFY_FIELD_PAGES_PRINTED,
	JOB_NOTIFY_FIELD_TOTAL_BYTES,
	JOB_NOTIFY_FIELD_BYTES_PRINTED,
	JOB_NOTIFY_FIELD_REMOTE_JOB_ID,
}

func JobNotifyFieldToString(field uint16) string {
	switch field {
	case JOB_NOTIFY_FIELD_PRINTER_NAME:
		return "Printer name"
	case JOB_NOTIFY_FIELD_MACHINE_NAME:
		return "Machine name"
	case JOB_NOTIFY_FIELD_PORT_NAME:
		return "Port name"
	case JOB_NOTIFY_FIELD_USER_NAME:
		return "User name"
	case JOB_NOTIFY_FIELD_NOTIFY_NAME:
		return "Notify name"
	case JOB_NOTIFY_FIELD_DATATYPE:
		return "Datatype"
	case JOB_NOTIFY_FIELD_PRINT_PROCESSOR:
		return "Print processor"
	case JOB_NOTIFY_FIELD_PARAMETERS:
		return "Parameters"
	case JOB_NOTIFY_FIELD_DRIVER_NAME:
		return "Driver name"
	case JOB_NOTIFY_FIELD_DEVMODE:
		return "Devmode"
	case JOB_NOTIFY_FIELD_STATUS:
		return "Status"
	case JOB_NOTIFY_FIELD_STATUS_STRING:
		return "Status(string)"
	case JOB_NOTIFY_FIELD_SECURITY_DESCRIPTOR:
		return "Security descriptor"
	case JOB_NOTIFY_FIELD_DOCUMENT:
		return "Document"
	case JOB_NOTIFY_FIELD_PRIORITY:
		return "Priority"
	case JOB_NOTIFY_FIELD_POSITION:
		return "Position"
	case JOB_NOTIFY_FIELD_SUBMITTED:
		return "Submitted time"
	case JOB_NOTIFY_FIELD_START_TIME:
		return "Start time"
	case JOB_NOTIFY_FIELD_UNTIL_TIME:
		return "Until time"
	case JOB_NOTIFY_FIELD_TIME:
		return "Time since start"
	case JOB_NOTIFY_FIELD_TOTAL_PAGES:
		return "Total pages"
	case JOB_NOTIFY_FIELD_PAGES_PRINTED:
		return "Pages printed"
	case JOB_NOTIFY_FIELD_TOTAL_BYTES:
		return "Total bytes"
	case JOB_NOTIFY_FIELD_BYTES_PRINTED:
		return "Bytes printed"
	case JOB_NOTIFY_FIELD_REMOTE_JOB_ID:
		return "Remote job id"
	}

	return "<UNKNOWN>"
}

func (pnid *PrinterNotifyInfoData) String() string {
	if pnid.Type == JOB_NOTIFY_TYPE {
		return fmt.Sprintf("Job #%d %s: %v", pnid.ID, JobNotifyFieldToString(pnid.Field), pnid.Value)
	} else if pnid.Type == PRINTER_NOTIFY_TYPE {
		return fmt.Sprintf("Printer Field %d Value %v", pnid.Field, pnid.Value)
	}

	return fmt.Sprintf("%#v\n", pnid)
}
