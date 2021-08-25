package goquery

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stoewer/go-strcase"
)

func RowsToCSV(rows Rows, toCamelCase bool, dateFormat string) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("column error: %v", err)
	}

	ct, err := rows.ColumnTypes()
	if err != nil {
		return "", fmt.Errorf("column type error: %v", err)
	}

	types := make([]reflect.Type, len(ct))
	for i, tp := range ct {
		types[i] = tp
	}

	values := make([]interface{}, len(ct))
	var builder strings.Builder

	for i, col := range columns {
		if i > 0 {
			builder.WriteString(",")
		}
		if toCamelCase {
			col = strcase.LowerCamelCase(col)
		}
		builder.WriteString(`"` + col + `"`)
	}
	builder.WriteString("\n")

	for rows.Next() {
		for i := range values {
			values[i] = reflect.New(types[i]).Interface()
		}
		err = rows.Scan(values...)
		if err != nil {
			return "", fmt.Errorf("failed to scan values: %v", err)
		}
		for i, v := range values {
			if i > 0 {
				builder.WriteString(",")
			}
			var valstring string
			switch v.(type) {
			case string, *string:
				valstring = fmt.Sprintf("\"%s\"", reflect.ValueOf(v).Elem().String())
			case int, *int, int32, *int32, int64, *int64:
				valstring = fmt.Sprintf("%d", reflect.ValueOf(v).Elem().Int())
			case float32, *float32, float64, *float64:
				valstring = fmt.Sprintf("%f", reflect.ValueOf(v).Elem().Float())
			case time.Time, *time.Time:
				if dateFormat != "" {
					valstring = reflect.ValueOf(v).Elem().Interface().(time.Time).Format(dateFormat)
				} else {
					valstring = reflect.ValueOf(v).Elem().Interface().(time.Time).String()
				}

			case uuid.UUID:
				valstring = v.(*uuid.UUID).String()
			default:
				return "", fmt.Errorf("unsupported csv conversion type: %v", reflect.ValueOf(v).Elem().Type())
			}

			builder.WriteString(fmt.Sprintf("%v", valstring))
		}
		builder.WriteString("\n")
	}
	return builder.String(), nil
}
