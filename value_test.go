package athena

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertValue(t *testing.T) {
	toPtr := func(s string) *string {
		return &s
	}

	tests := map[string]struct {
		athenaType string
		rawValue   *string
		want       any
	}{
		"null": {
			athenaType: "integer",
			rawValue:   nil,
			want:       nil,
		},
		"tinyint": {
			athenaType: "tinyint",
			rawValue:   toPtr("123"),
			want:       int64(123),
		},
		"smallint": {
			athenaType: "smallint",
			rawValue:   toPtr("123"),
			want:       int64(123),
		},
		"integer": {
			athenaType: "integer",
			rawValue:   toPtr("2147483647"),
			want:       int64(2147483647),
		},
		"bigint": {
			athenaType: "bigint",
			rawValue:   toPtr("9223372036854775807"),
			want:       int64(9223372036854775807),
		},
		"boolean true": {
			athenaType: "boolean",
			rawValue:   toPtr("true"),
			want:       true,
		},
		"boolean false": {
			athenaType: "boolean",
			rawValue:   toPtr("false"),
			want:       false,
		},
		"float": {
			athenaType: "float",
			rawValue:   toPtr("1.75"),
			want:       float64(1.75),
		},
		"double": {
			athenaType: "double",
			rawValue:   toPtr("1.75"),
			want:       float64(1.75),
		},
		"decimal": {
			athenaType: "decimal",
			rawValue:   toPtr("123.45"),
			want:       float64(123.45),
		},
		"varchar": {
			athenaType: "varchar",
			rawValue:   toPtr("hello world"),
			want:       "hello world",
		},
		"timestamp": {
			athenaType: "timestamp",
			rawValue:   toPtr("2023-01-15 12:34:56.789"),
			want:       "2023-01-15 12:34:56.789",
		},
		"timestamp with time zone": {
			athenaType: "timestamp with time zone",
			rawValue:   toPtr("2023-01-15 12:34:56.789 JST"),
			want:       "2023-01-15 12:34:56.789 JST",
		},
		"date": {
			athenaType: "date",
			rawValue:   toPtr("2023-01-15"),
			want:       "2023-01-15",
		},
		"unknown type": {
			athenaType: "unknown",
			rawValue:   toPtr("unknown"),
			want:       []byte("unknown"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := convertValue(tc.athenaType, tc.rawValue)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			assert.Equal(t, got, tc.want)
		})
	}
}
