package errors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value any
		name  string
		key   string
		want  Attr
	}{
		{
			name:  "given_string_value_when_any_then_returns_attr_with_any_type",
			key:   "test_key",
			value: "test_value",
			want: Attr{
				Type:  AnyType,
				Key:   "test_key",
				Value: "test_value",
			},
		},
		{
			name:  "given_int_value_when_any_then_returns_attr_with_any_type",
			key:   "number",
			value: 42,
			want: Attr{
				Type:  AnyType,
				Key:   "number",
				Value: 42,
			},
		},
		{
			name:  "given_nil_value_when_any_then_returns_attr_with_any_type",
			key:   "nil_key",
			value: nil,
			want: Attr{
				Type:  AnyType,
				Key:   "nil_key",
				Value: nil,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Any(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []Attr
	}{
		{
			name:  "given_empty_attrs_when_object_then_returns_attr_with_object_type",
			key:   "empty",
			value: []Attr{},
			want: Attr{
				Type:  ObjectType,
				Key:   "empty",
				Value: []Attr{},
			},
		},
		{
			name: "given_multiple_attrs_when_object_then_returns_attr_with_object_type",
			key:  "multiple",
			value: []Attr{
				{Type: StringType, Key: "name", Value: "test"},
				{Type: IntType, Key: "age", Value: 30},
			},
			want: Attr{
				Type: ObjectType,
				Key:  "multiple",
				Value: []Attr{
					{Type: StringType, Key: "name", Value: "test"},
					{Type: IntType, Key: "age", Value: 30},
				},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Object(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		want  Attr
		value bool
	}{
		{
			name:  "given_true_value_when_bool_then_returns_attr_with_bool_type",
			key:   "is_active",
			value: true,
			want: Attr{
				Type:  BoolType,
				Key:   "is_active",
				Value: true,
			},
		},
		{
			name:  "given_false_value_when_bool_then_returns_attr_with_bool_type",
			key:   "is_disabled",
			value: false,
			want: Attr{
				Type:  BoolType,
				Key:   "is_disabled",
				Value: false,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Bool(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestBools(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []bool
	}{
		{
			name:  "given_empty_slice_when_bools_then_returns_attr_with_bools_type",
			key:   "flags",
			value: []bool{},
			want: Attr{
				Type:  BoolsType,
				Key:   "flags",
				Value: []bool{},
			},
		},
		{
			name:  "given_multiple_bools_when_bools_then_returns_attr_with_bools_type",
			key:   "multiple",
			value: []bool{true, false, true},
			want: Attr{
				Type:  BoolsType,
				Key:   "multiple",
				Value: []bool{true, false, true},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Bools(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestTime(t *testing.T) {
	fixedTime := time.Date(2023, 10, 15, 12, 30, 0, 0, time.UTC)

	t.Parallel()

	tests := []struct {
		value time.Time
		name  string
		key   string
		want  Attr
	}{
		{
			name:  "given_time_value_when_time_then_returns_attr_with_time_type",
			key:   "created_at",
			value: fixedTime,
			want: Attr{
				Type:  TimeType,
				Key:   "created_at",
				Value: fixedTime,
			},
		},
		{
			name:  "given_zero_time_when_time_then_returns_attr_with_time_type",
			key:   "zero_time",
			value: time.Time{},
			want: Attr{
				Type:  TimeType,
				Key:   "zero_time",
				Value: time.Time{},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Time(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestTimes(t *testing.T) {
	time1 := time.Date(2023, 10, 15, 12, 30, 0, 0, time.UTC)
	time2 := time.Date(2023, 10, 16, 12, 30, 0, 0, time.UTC)

	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []time.Time
	}{
		{
			name:  "given_empty_slice_when_times_then_returns_attr_with_times_type",
			key:   "timestamps",
			value: []time.Time{},
			want: Attr{
				Type:  TimesType,
				Key:   "timestamps",
				Value: []time.Time{},
			},
		},
		{
			name:  "given_multiple_times_when_times_then_returns_attr_with_times_type",
			key:   "multiple",
			value: []time.Time{time1, time2},
			want: Attr{
				Type:  TimesType,
				Key:   "multiple",
				Value: []time.Time{time1, time2},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Times(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		want  Attr
		value time.Duration
	}{
		{
			name:  "given_seconds_duration_when_duration_then_returns_attr_with_duration_type",
			key:   "timeout",
			value: 5 * time.Second,
			want: Attr{
				Type:  DurationType,
				Key:   "timeout",
				Value: 5 * time.Second,
			},
		},
		{
			name:  "given_zero_duration_when_duration_then_returns_attr_with_duration_type",
			key:   "zero",
			value: 0,
			want: Attr{
				Type:  DurationType,
				Key:   "zero",
				Value: time.Duration(0),
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Duration(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestDurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []time.Duration
	}{
		{
			name:  "given_empty_slice_when_durations_then_returns_attr_with_durations_type",
			key:   "durations",
			value: []time.Duration{},
			want: Attr{
				Type:  DurationsType,
				Key:   "durations",
				Value: []time.Duration{},
			},
		},
		{
			name:  "given_multiple_durations_when_durations_then_returns_attr_with_durations_type",
			key:   "multiple",
			value: []time.Duration{5 * time.Second, 100 * time.Millisecond},
			want: Attr{
				Type:  DurationsType,
				Key:   "multiple",
				Value: []time.Duration{5 * time.Second, 100 * time.Millisecond},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Durations(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		want  Attr
		value int
	}{
		{
			name:  "given_positive_int_when_int_then_returns_attr_with_int_type",
			key:   "count",
			value: 42,
			want: Attr{
				Type:  IntType,
				Key:   "count",
				Value: 42,
			},
		},
		{
			name:  "given_negative_int_when_int_then_returns_attr_with_int_type",
			key:   "negative",
			value: -10,
			want: Attr{
				Type:  IntType,
				Key:   "negative",
				Value: -10,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Int(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestInts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []int
	}{
		{
			name:  "given_empty_slice_when_ints_then_returns_attr_with_ints_type",
			key:   "numbers",
			value: []int{},
			want: Attr{
				Type:  IntsType,
				Key:   "numbers",
				Value: []int{},
			},
		},
		{
			name:  "given_multiple_ints_when_ints_then_returns_attr_with_ints_type",
			key:   "multiple",
			value: []int{1, 2, 3, 4, 5},
			want: Attr{
				Type:  IntsType,
				Key:   "multiple",
				Value: []int{1, 2, 3, 4, 5},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Ints(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		want  Attr
		value int64
	}{
		{
			name:  "given_positive_int64_when_int64_then_returns_attr_with_int64_type",
			key:   "id",
			value: 9223372036854775807,
			want: Attr{
				Type:  Int64Type,
				Key:   "id",
				Value: int64(9223372036854775807),
			},
		},
		{
			name:  "given_zero_int64_when_int64_then_returns_attr_with_int64_type",
			key:   "zero",
			value: 0,
			want: Attr{
				Type:  Int64Type,
				Key:   "zero",
				Value: int64(0),
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Int64(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestInt64s(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []int64
	}{
		{
			name:  "given_empty_slice_when_int64s_then_returns_attr_with_int64s_type",
			key:   "ids",
			value: []int64{},
			want: Attr{
				Type:  Int64sType,
				Key:   "ids",
				Value: []int64{},
			},
		},
		{
			name:  "given_multiple_int64s_when_int64s_then_returns_attr_with_int64s_type",
			key:   "multiple",
			value: []int64{1, 2, 3, 4, 5},
			want: Attr{
				Type:  Int64sType,
				Key:   "multiple",
				Value: []int64{1, 2, 3, 4, 5},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Int64s(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		want  Attr
		value uint64
	}{
		{
			name:  "given_max_uint64_when_uint64_then_returns_attr_with_uint64_type",
			key:   "max",
			value: 18446744073709551615,
			want: Attr{
				Type:  Uint64Type,
				Key:   "max",
				Value: uint64(18446744073709551615),
			},
		},
		{
			name:  "given_zero_uint64_when_uint64_then_returns_attr_with_uint64_type",
			key:   "zero",
			value: 0,
			want: Attr{
				Type:  Uint64Type,
				Key:   "zero",
				Value: uint64(0),
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Uint64(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestUint64s(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []uint64
	}{
		{
			name:  "given_empty_slice_when_uint64s_then_returns_attr_with_uint64s_type",
			key:   "ids",
			value: []uint64{},
			want: Attr{
				Type:  Uint64sType,
				Key:   "ids",
				Value: []uint64{},
			},
		},
		{
			name:  "given_multiple_uint64s_when_uint64s_then_returns_attr_with_uint64s_type",
			key:   "multiple",
			value: []uint64{1, 2, 3, 4, 5},
			want: Attr{
				Type:  Uint64sType,
				Key:   "multiple",
				Value: []uint64{1, 2, 3, 4, 5},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Uint64s(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		want  Attr
		value float64
	}{
		{
			name:  "given_positive_float64_when_float64_then_returns_attr_with_float64_type",
			key:   "price",
			value: 99.99,
			want: Attr{
				Type:  Float64Type,
				Key:   "price",
				Value: 99.99,
			},
		},
		{
			name:  "given_negative_float64_when_float64_then_returns_attr_with_float64_type",
			key:   "temperature",
			value: -15.5,
			want: Attr{
				Type:  Float64Type,
				Key:   "temperature",
				Value: -15.5,
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Float64(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestFloat64s(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []float64
	}{
		{
			name:  "given_empty_slice_when_float64s_then_returns_attr_with_float64s_type",
			key:   "prices",
			value: []float64{},
			want: Attr{
				Type:  Float64sType,
				Key:   "prices",
				Value: []float64{},
			},
		},
		{
			name:  "given_multiple_float64s_when_float64s_then_returns_attr_with_float64s_type",
			key:   "multiple",
			value: []float64{1.1, 2.2, 3.3},
			want: Attr{
				Type:  Float64sType,
				Key:   "multiple",
				Value: []float64{1.1, 2.2, 3.3},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Float64s(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		key   string
		value string
		want  Attr
	}{
		{
			name:  "given_non_empty_string_when_string_then_returns_attr_with_string_type",
			key:   "message",
			value: "hello world",
			want: Attr{
				Type:  StringType,
				Key:   "message",
				Value: "hello world",
			},
		},
		{
			name:  "given_empty_string_when_string_then_returns_attr_with_string_type",
			key:   "empty",
			value: "",
			want: Attr{
				Type:  StringType,
				Key:   "empty",
				Value: "",
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := String(test.key, test.value)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
				assert.Equal(t, test.want.Value, got.Value)
			},
		)
	}
}

func TestStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  Attr
		name  string
		key   string
		value []string
	}{
		{
			name:  "given_empty_slice_when_strings_then_returns_attr_with_strings_type",
			key:   "tags",
			value: []string{},
			want: Attr{
				Type:  StringsType,
				Key:   "tags",
				Value: []string{},
			},
		},
		{
			name:  "given_multiple_strings_when_strings_then_returns_attr_with_strings_type",
			key:   "multiple",
			value: []string{"one", "two", "three"},
			want: Attr{
				Type:  StringsType,
				Key:   "multiple",
				Value: []string{"one", "two", "three"},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Strings(test.key, test.value...)

				// then
				assert.Equal(t, test.want.Type, got.Type)
				assert.Equal(t, test.want.Key, got.Key)
			},
		)
	}
}
