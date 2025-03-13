package athena

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

const (
	// TimestampLayout is the Go time layout string for an Athena `timestamp`.
	TimestampLayout             = "2006-01-02 15:04:05.999"
	TimestampWithTimeZoneLayout = "2006-01-02 15:04:05.999 MST"
	DateLayout                  = "2006-01-02"
)

const nullStringResultModeGzipDL string = "\\N"

func convertRow(columns []types.ColumnInfo, in []types.Datum, ret []driver.Value) error {
	for i, val := range in {
		coerced, err := convertValue(*columns[i].Type, val.VarCharValue)
		if err != nil {
			return err
		}

		ret[i] = coerced
	}

	return nil
}

func convertRowFromTableInfo(columns []types.Column, in []string, ret []driver.Value) error {
	for i, val := range in {
		var coerced any
		var err error
		if val == nullStringResultModeGzipDL {
			var nullVal *string
			coerced, err = convertValue(*columns[i].Type, nullVal)
		} else {
			coerced, err = convertValue(*columns[i].Type, &val)
		}
		if err != nil {
			return err
		}

		ret[i] = coerced
	}

	return nil
}

func convertRowFromCsv(columns []types.ColumnInfo, in []downloadField, ret []driver.Value) error {
	for i, df := range in {
		var coerced any
		var err error
		if df.isNil {
			var nullVal *string
			coerced, err = convertValue(*columns[i].Type, nullVal)
		} else {
			coerced, err = convertValue(*columns[i].Type, &df.val)
		}
		if err != nil {
			return err
		}

		ret[i] = coerced
	}

	return nil
}

func convertValue(athenaType string, rawValue *string) (any, error) {
	if rawValue == nil {
		return nil, nil
	}

	if len(athenaType) > 7 && athenaType[:7] == "decimal" {
		athenaType = "decimal"
	}

	val := *rawValue
	switch athenaType {
	case "tinyint":
		return strconv.ParseInt(val, 10, 8)
	case "smallint":
		return strconv.ParseInt(val, 10, 16)
	case "integer":
		return strconv.ParseInt(val, 10, 32)
	case "bigint":
		return strconv.ParseInt(val, 10, 64)
	case "boolean":
		switch val {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
		return nil, fmt.Errorf("cannot parse '%s' as boolean", val)
	case "float":
		return strconv.ParseFloat(val, 32)
	case "double", "decimal":
		return strconv.ParseFloat(val, 64)
	case "varchar",
		"timestamp", "timestamp with time zone",
		"date",
		"time", "time with time zone":
		return val, nil
	case "array":
		// Child element type is assumed to be any, because types.ColumnInfo does not contain the child element type.
		var array []any
		if err := json.Unmarshal([]byte(val), &array); err != nil {
			return nil, err
		}
		return array, nil
	default:
		return []byte(val), nil
	}
}
