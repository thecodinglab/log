package file

import (
	"bytes"
	"encoding/json"

	"github.com/thecodinglab/log"
)

type JSONFormatter struct {
}

func (f JSONFormatter) Format(buffer *bytes.Buffer, entry *log.Entry) error {
	return json.NewEncoder(buffer).Encode(entry)
}
