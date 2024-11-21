package report_processor

import (
	"io"
)

type ReportProcessor interface {
	Process(r io.Reader, rt string) error
}
