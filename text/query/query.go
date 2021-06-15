package query

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

var (
	ErrInvalidInteger         = errors.New("invalid integer")
	ErrInvalidFloatNumber     = errors.New("invalid float number")
	ErrInvalidBoolean         = errors.New("invalid boolean")
	ErrRequiredArgumentNotSet = errors.New("required argument not set")
)

type Query = map[string][]string

func getArgument(q Query, key string, required bool) (value string, err error) {
	if vs := q[key]; len(vs) > 0 {
		value = vs[0]
	} else if required {
		err = ErrRequiredArgumentNotSet
	}
	return
}

func parseInt64(q Query, key string, required bool, dft int64) (int64, error) {
	if value, err := getArgument(q, key, required); err != nil {
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

func parseUint64(q Query, key string, required bool, dft uint64) (uint64, error) {
	if value, err := getArgument(q, key, required); err != nil {
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

func parseFloat64(q Query, key string, required bool, dft float64) (float64, error) {
	if value, err := getArgument(q, key, required); err != nil {
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

func Int(q Query, key string, dft int) (int, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int(i), nil
}

func RequiredInt(q Query, key string) (int, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func Int8(q Query, key string, dft int8) (int8, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int8(i), nil
}

func RequiredInt8(q Query, key string) (int8, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

func Int16(q Query, key string, dft int16) (int16, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int16(i), nil
}

func RequiredInt16(q Query, key string) (int16, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

func Int32(q Query, key string, dft int32) (int32, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int32(i), nil
}

func RequiredInt32(q Query, key string) (int32, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func Int64(q Query, key string, dft int64) (int64, error) {
	return parseInt64(q, key, false, dft)
}

func RequiredInt64(q Query, key string) (int64, error) {
	return parseInt64(q, key, true, 0)
}

func Uint(q Query, key string, dft uint) (uint, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint(i), nil
}

func RequiredUint(q Query, key string) (uint, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func Uint8(q Query, key string, dft uint8) (uint8, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint8(i), nil
}

func RequiredUint8(q Query, key string) (uint8, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

func Uint16(q Query, key string, dft uint16) (uint16, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint16(i), nil
}

func RequiredUint16(q Query, key string) (uint16, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

func Uint32(q Query, key string, dft uint32) (uint32, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint32(i), nil
}

func RequiredUint32(q Query, key string) (uint32, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func Uint64(q Query, key string, dft uint64) (uint64, error) {
	return parseUint64(q, key, false, dft)
}

func RequiredUint64(q Query, key string) (uint64, error) {
	return parseUint64(q, key, true, 0)
}

func Float32(q Query, key string, dft float32) (float32, error) {
	f, err := parseFloat64(q, key, false, float64(dft))
	if err != nil {
		return dft, err
	}
	return float32(f), err
}

func RequiredFloat32(q Query, key string) (float32, error) {
	f, err := parseFloat64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return float32(f), err
}

func Float64(q Query, key string, dft float64) (float64, error) {
	return parseFloat64(q, key, false, dft)
}

func RequiredFloat64(q Query, key string) (float64, error) {
	return parseFloat64(q, key, true, 0)
}

func Bool(q Query, key string, dft bool) (bool, error) {
	if value, err := getArgument(q, key, false); err != nil {
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

func RequiredBool(q Query, key string) (bool, error) {
	if value, err := getArgument(q, key, true); err != nil {
		return false, err
	} else {
		x, err := strconv.ParseBool(value)
		if err != nil {
			return false, ErrInvalidBoolean
		}
		return x, nil
	}
}

func String(q Query, key string, dft string) string {
	if value, _ := getArgument(q, key, false); value == "" {
		return dft
	} else {
		return value
	}
}

func RequiredString(q Query, key string) (string, error) {
	if value, err := getArgument(q, key, true); err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func JSON(q Query, key string, ptr interface{}) error {
	if value, _ := getArgument(q, key, false); value == "" {
		return nil
	} else {
		return json.Unmarshal([]byte(value), ptr)
	}
}

func RequiredDuration(q Query, key string) (time.Duration, error) {
	if value, err := getArgument(q, key, true); err != nil {
		return 0, err
	} else {
		return time.ParseDuration(value)
	}
}

func Duration(q Query, key string, dft time.Duration) (time.Duration, error) {
	if value, _ := getArgument(q, key, false); value == "" {
		return dft, nil
	} else {
		return time.ParseDuration(value)
	}
}
