package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrInvalidInteger     = errors.New("invalid integer")
	ErrInvalidFloatNumber = errors.New("invalid float number")
	ErrInvalidBoolean     = errors.New("invalid boolean")
)

func getArgument(r *http.Request, key string, required bool) (value string, err error) {
	if required {
		if r.Form == nil {
			r.ParseMultipartForm(32 << 20)
		}
		if vs := r.Form[key]; len(vs) > 0 {
			value = vs[0]
		} else {
			err = ErrMissingRequiredArgument
		}
	} else {
		value = r.FormValue(key)
	}
	return
}

func parseInt64(r *http.Request, key string, required bool, dft int64) (int64, error) {
	if value, err := getArgument(r, key, required); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return dft, ErrInvalidInteger
		}
		return x, nil
	}
}

func parseUint64(r *http.Request, key string, required bool, dft uint64) (uint64, error) {
	if value, err := getArgument(r, key, required); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return dft, ErrInvalidInteger
		}
		return x, nil
	}
}

func parseFloat64(r *http.Request, key string, required bool, dft float64) (float64, error) {
	if value, err := getArgument(r, key, required); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return dft, ErrInvalidFloatNumber
		}
		return x, nil
	}
}

func ParseInt(r *http.Request, key string, dft int) (int, error) {
	i, err := parseInt64(r, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int(i), nil
}

func ParseRequiredInt(r *http.Request, key string) (int, error) {
	i, err := parseInt64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func ParseInt8(r *http.Request, key string, dft int8) (int8, error) {
	i, err := parseInt64(r, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int8(i), nil
}

func ParseRequiredInt8(r *http.Request, key string) (int8, error) {
	i, err := parseInt64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

func ParseInt16(r *http.Request, key string, dft int16) (int16, error) {
	i, err := parseInt64(r, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int16(i), nil
}

func ParseRequiredInt16(r *http.Request, key string) (int16, error) {
	i, err := parseInt64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

func ParseInt32(r *http.Request, key string, dft int32) (int32, error) {
	i, err := parseInt64(r, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int32(i), nil
}

func ParseRequiredInt32(r *http.Request, key string) (int32, error) {
	i, err := parseInt64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func ParseInt64(r *http.Request, key string, dft int64) (int64, error) {
	return parseInt64(r, key, false, dft)
}

func ParseRequiredInt64(r *http.Request, key string) (int64, error) {
	return parseInt64(r, key, true, 0)
}

func ParseUint(r *http.Request, key string, dft uint) (uint, error) {
	i, err := parseUint64(r, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint(i), nil
}

func ParseRequiredUint(r *http.Request, key string) (uint, error) {
	i, err := parseUint64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func ParseUint8(r *http.Request, key string, dft uint8) (uint8, error) {
	i, err := parseUint64(r, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint8(i), nil
}

func ParseRequiredUint8(r *http.Request, key string) (uint8, error) {
	i, err := parseUint64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

func ParseUint16(r *http.Request, key string, dft uint16) (uint16, error) {
	i, err := parseUint64(r, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint16(i), nil
}

func ParseRequiredUint16(r *http.Request, key string) (uint16, error) {
	i, err := parseUint64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

func ParseUint32(r *http.Request, key string, dft uint32) (uint32, error) {
	i, err := parseUint64(r, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint32(i), nil
}

func ParseRequiredUint32(r *http.Request, key string) (uint32, error) {
	i, err := parseUint64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func ParseUint64(r *http.Request, key string, dft uint64) (uint64, error) {
	return parseUint64(r, key, false, dft)
}

func ParseRequiredUint64(r *http.Request, key string) (uint64, error) {
	return parseUint64(r, key, true, 0)
}

func ParseFloat32(r *http.Request, key string, dft float32) (float32, error) {
	f, err := parseFloat64(r, key, false, float64(dft))
	if err != nil {
		return dft, err
	}
	return float32(f), err
}

func ParseRequiredFloat32(r *http.Request, key string) (float32, error) {
	f, err := parseFloat64(r, key, true, 0)
	if err != nil {
		return 0, err
	}
	return float32(f), err
}

func ParseFloat64(r *http.Request, key string, dft float64) (float64, error) {
	return parseFloat64(r, key, false, dft)
}

func ParseRequiredFloat64(r *http.Request, key string) (float64, error) {
	return parseFloat64(r, key, true, 0)
}

func ParseBool(r *http.Request, key string, dft bool) (bool, error) {
	if value, err := getArgument(r, key, false); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseBool(value)
		if err != nil {
			return dft, ErrInvalidBoolean
		}
		return x, nil
	}
}

func ParseRequiredBool(r *http.Request, key string) (bool, error) {
	if value, err := getArgument(r, key, true); err != nil {
		return false, err
	} else {
		x, err := strconv.ParseBool(value)
		if err != nil {
			return false, ErrInvalidBoolean
		}
		return x, nil
	}
}

func ParseString(r *http.Request, key string, dft string) string {
	if value, _ := getArgument(r, key, false); value == "" {
		return dft
	} else {
		return value
	}
}

func ParseRequiredString(r *http.Request, key string) (string, error) {
	if value, err := getArgument(r, key, true); err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func ParseIntegers(r *http.Request, key string) ([]int64, error) {
	if value, _ := getArgument(r, key, false); value == "" {
		return []int64{}, nil
	} else {
		strs := strings.Split(value, ",")
		values := make([]int64, 0, len(strs))
		for _, s := range strs {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			i, e := strconv.ParseInt(s, 10, 64)
			if e != nil {
				return nil, e
			}
			values = append(values, i)
		}
		return values, nil
	}
}

func ParseJSON(r *http.Request, key string, ptr interface{}) error {
	if value, _ := getArgument(r, key, false); value == "" {
		return nil
	} else {
		return json.Unmarshal([]byte(value), ptr)
	}
}
