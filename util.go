package estree

import (
	"bytes"
	"encoding/json"
)

// isNullOrEmptyRawMessage is used when unmarshaling optional fields, to treat
// null and missing values as equivalent.
func isNullOrEmptyRawMessage(m json.RawMessage) bool {
	return len(m) == 0 || bytes.Equal(m, []byte("null"))
}
