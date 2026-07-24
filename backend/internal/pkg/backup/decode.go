package backup

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

func decodeJSONValue(raw interface{}) (interface{}, error) {
	if raw == nil {
		return nil, nil
	}

	switch value := raw.(type) {
	case map[string]interface{}:
		typed, _ := value["__type"].(string)
		if typed == "bytes" {
			hexStr, _ := value["hex"].(string)
			decoded, err := hex.DecodeString(hexStr)
			if err != nil {
				return nil, err
			}
			return decoded, nil
		}
		return nil, fmt.Errorf("unknown encoded value type %q", typed)
	case string:
		if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
			return parsed, nil
		}
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			return parsed, nil
		}
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			return parsed, nil
		}
		return value, nil
	case float64:
		if value == float64(int64(value)) {
			return int64(value), nil
		}
		return value, nil
	case json.Number:
		if integer, err := value.Int64(); err == nil {
			return integer, nil
		}
		return value.Float64()
	default:
		return value, nil
	}
}
