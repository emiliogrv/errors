package errors

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogObjectMarshalerFuncMarshalZerologObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fn         LogObjectMarshalerFunc
		name       string
		wantCalled bool
	}{
		{
			name: "given_func_when_marshal_zerolog_object_then_calls_func",
			fn: func(e *zerolog.Event) {
				e.Str("test", "value")
			},
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.fn.MarshalZerologObject(event)
				event.Msg("test")

				// then
				if test.wantCalled {
					assert.Contains(t, buf.String(), "test")
				}
			},
		)
	}
}

func TestLogArrayMarshalerFuncMarshalZerologArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fn         LogArrayMarshalerFunc
		name       string
		wantCalled bool
	}{
		{
			name: "given_func_when_marshal_zerolog_array_then_calls_func",
			fn: func(a *zerolog.Array) {
				a.Str("value")
			},
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				event.Array("test", test.fn)
				event.Msg("test")

				// then
				if test.wantCalled {
					assert.Contains(t, buf.String(), "value")
				}
			},
		)
	}
}

func TestStructuredErrorMarshalZerologObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		err *StructuredError
		// then
		wantContains []string
	}{
		{
			name:         "given_nil_error_when_marshal_zerolog_object_then_has_nil_message",
			err:          nil,
			wantContains: []string{`"message":"!NILVALUE"`},
		},
		{
			name:         "given_error_with_message_when_marshal_zerolog_object_then_has_message",
			err:          New("test error"),
			wantContains: []string{`"message":"test error"`},
		},
		{
			name:         "given_error_with_tags_when_marshal_zerolog_object_then_has_tags",
			err:          New("test").WithTags("tag1", "tag2"),
			wantContains: []string{`"message":"test"`, `"tags":`},
		},
		{
			name:         "given_error_with_attrs_when_marshal_zerolog_object_then_has_attrs",
			err:          New("test").WithAttrs(String("key", "value")),
			wantContains: []string{`"message":"test"`, `"attrs":`},
		},
		{
			name:         "given_error_with_errors_when_marshal_zerolog_object_then_has_errors",
			err:          New("parent").WithErrors(stderrors.New("child")),
			wantContains: []string{`"message":"parent"`, `"errors":`},
		},
		{
			name:         "given_error_with_stack_when_marshal_zerolog_object_then_has_stack",
			err:          New("test").WithStack([]byte("stack trace")),
			wantContains: []string{`"message":"test"`, `"stack":`},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.err.MarshalZerologObject(event)
				event.Msg("test")

				// then
				got := buf.String()
				for _, want := range test.wantContains {
					assert.Contains(t, got, want)
				}
			},
		)
	}
}

func TestStructuredErrorMarshalZerologObjectFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		err *StructuredError
		// then
		wantKeys []string
	}{
		{
			name:     "given_error_with_message_when_marshal_zerolog_object_then_has_message_key",
			err:      New("test"),
			wantKeys: []string{"message"},
		},
		{
			name:     "given_error_with_tags_when_marshal_zerolog_object_then_has_message_and_tags",
			err:      New("test").WithTags("tag1"),
			wantKeys: []string{"message", "tags"},
		},
		{
			name:     "given_error_with_attrs_when_marshal_zerolog_object_then_has_message_and_attrs",
			err:      New("test").WithAttrs(String("key", "value")),
			wantKeys: []string{"message", "attrs"},
		},
		{
			name:     "given_error_with_errors_when_marshal_zerolog_object_then_has_message_and_errors",
			err:      New("parent").WithErrors(stderrors.New("child")),
			wantKeys: []string{"message", "errors"},
		},
		{
			name:     "given_error_with_stack_when_marshal_zerolog_object_then_has_message_and_stack",
			err:      New("test").WithStack([]byte("stack")),
			wantKeys: []string{"message", "stack"},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.err.MarshalZerologObject(event)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)

				for _, key := range test.wantKeys {
					assert.Contains(t, result, key)
				}
			},
		)
	}
}

func TestAttrMarshalZerologObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantContains []string
	}{
		{
			name:         "given_nil_attr_when_marshal_zerolog_object_then_has_nil_key",
			attr:         nil,
			wantContains: []string{`"!NILVALUE":"!NILVALUE"`},
		},
		{
			name:         "given_string_attr_when_marshal_zerolog_object_then_has_key_value",
			attr:         &Attr{Type: StringType, Key: "name", Value: "test"},
			wantContains: []string{`"name":"test"`},
		},
		{
			name:         "given_int_attr_when_marshal_zerolog_object_then_has_key_value",
			attr:         &Attr{Type: IntType, Key: "count", Value: 42},
			wantContains: []string{`"count":42`},
		},
		{
			name:         "given_bool_attr_when_marshal_zerolog_object_then_has_key_value",
			attr:         &Attr{Type: BoolType, Key: "active", Value: true},
			wantContains: []string{`"active":true`},
		},
		{
			name:         "given_float64_attr_when_marshal_zerolog_object_then_has_key_value",
			attr:         &Attr{Type: Float64Type, Key: "price", Value: 99.99},
			wantContains: []string{`"price":99.99`},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.attr.MarshalZerologObject(event)
				event.Msg("test")

				// then
				got := buf.String()
				for _, want := range test.wantContains {
					assert.Contains(t, got, want)
				}
			},
		)
	}
}

func TestAttrMarshalZerologObjectWithTypes(t *testing.T) {
	fixedTime := time.Date(2023, 10, 15, 12, 30, 0, 0, time.UTC)
	fixedDuration := 5 * time.Second

	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantKey string
	}{
		{
			name:    "given_string_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: StringType, Key: "name", Value: "test"},
			wantKey: "name",
		},
		{
			name:    "given_int_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: IntType, Key: "count", Value: 42},
			wantKey: "count",
		},
		{
			name:    "given_int64_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: Int64Type, Key: "id", Value: int64(123)},
			wantKey: "id",
		},
		{
			name:    "given_uint64_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: Uint64Type, Key: "uid", Value: uint64(456)},
			wantKey: "uid",
		},
		{
			name:    "given_float64_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: Float64Type, Key: "price", Value: 99.99},
			wantKey: "price",
		},
		{
			name:    "given_bool_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: BoolType, Key: "active", Value: true},
			wantKey: "active",
		},
		{
			name:    "given_time_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: TimeType, Key: "created", Value: fixedTime},
			wantKey: "created",
		},
		{
			name:    "given_duration_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: DurationType, Key: "timeout", Value: fixedDuration},
			wantKey: "timeout",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.attr.MarshalZerologObject(event)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)
				assert.Contains(t, result, test.wantKey)
			},
		)
	}
}

func TestAttrMarshalZerologObjectWithSlices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantKey string
	}{
		{
			name:    "given_bools_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: BoolsType, Key: "flags", Value: []bool{true, false}},
			wantKey: "flags",
		},
		{
			name:    "given_ints_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: IntsType, Key: "numbers", Value: []int{1, 2, 3}},
			wantKey: "numbers",
		},
		{
			name:    "given_int64s_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: Int64sType, Key: "ids", Value: []int64{1, 2, 3}},
			wantKey: "ids",
		},
		{
			name:    "given_uint64s_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: Uint64sType, Key: "uids", Value: []uint64{1, 2, 3}},
			wantKey: "uids",
		},
		{
			name:    "given_float64s_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: Float64sType, Key: "prices", Value: []float64{1.1, 2.2}},
			wantKey: "prices",
		},
		{
			name:    "given_strings_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: StringsType, Key: "tags", Value: []string{"tag1", "tag2"}},
			wantKey: "tags",
		},
		{
			name:    "given_times_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: TimesType, Key: "timestamps", Value: []time.Time{time.Now()}},
			wantKey: "timestamps",
		},
		{
			name:    "given_durations_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: DurationsType, Key: "timeouts", Value: []time.Duration{5 * time.Second}},
			wantKey: "timeouts",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.attr.MarshalZerologObject(event)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)
				assert.Contains(t, result, test.wantKey)
			},
		)
	}
}

func TestAttrMarshalZerologObjectWithObjectType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantKey string
	}{
		{
			name:    "given_object_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: ObjectType, Key: "obj", Value: []Attr{String("key", "value")}},
			wantKey: "obj",
		},
		{
			name:    "given_empty_object_attr_when_marshal_zerolog_object_then_has_key",
			attr:    &Attr{Type: ObjectType, Key: "empty", Value: []Attr{}},
			wantKey: "empty",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.attr.MarshalZerologObject(event)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)
				assert.Contains(t, result, test.wantKey)
			},
		)
	}
}

func TestErrorToZerolog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		err error
		// then
		wantContains []string
	}{
		{
			name:         "given_nil_error_when_error_to_zerolog_then_has_nil_message",
			err:          nil,
			wantContains: []string{`"message":"!NILVALUE"`},
		},
		{
			name:         "given_standard_error_when_error_to_zerolog_then_has_message",
			err:          stderrors.New("standard error"),
			wantContains: []string{`"message":"standard error"`},
		},
		{
			name:         "given_structured_error_when_error_to_zerolog_then_has_message",
			err:          New("structured"),
			wantContains: []string{`"message":"structured"`},
		},
		{
			name:         "given_error_with_whitespace_when_error_to_zerolog_then_trims_whitespace",
			err:          stderrors.New("  error  "),
			wantContains: []string{`"message":"error"`},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				errorToZerolog(event, test.err)
				event.Msg("test")

				// then
				got := buf.String()
				for _, want := range test.wantContains {
					assert.Contains(t, got, want)
				}
			},
		)
	}
}

func TestSliceToZerologWithStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		wantKey string
		slice   []string
	}{
		{
			name:    "given_empty_slice_when_slice_to_zerolog_then_has_key",
			key:     "tags",
			slice:   []string{},
			wantKey: "tags",
		},
		{
			name:    "given_single_item_slice_when_slice_to_zerolog_then_has_key",
			key:     "tags",
			slice:   []string{"tag1"},
			wantKey: "tags",
		},
		{
			name:    "given_multiple_items_slice_when_slice_to_zerolog_then_has_key",
			key:     "tags",
			slice:   []string{"tag1", "tag2", "tag3"},
			wantKey: "tags",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				sliceToZerolog(event, test.key, test.slice)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)
				assert.Contains(t, result, test.wantKey)
			},
		)
	}
}

func TestSliceToZerologWithErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		wantKey string
		errs    []error
	}{
		{
			name:    "given_empty_errors_slice_when_slice_to_zerolog_then_has_key",
			key:     "errors",
			errs:    []error{},
			wantKey: "errors",
		},
		{
			name:    "given_single_error_when_slice_to_zerolog_then_has_key",
			key:     "errors",
			errs:    []error{stderrors.New("error1")},
			wantKey: "errors",
		},
		{
			name:    "given_multiple_errors_when_slice_to_zerolog_then_has_key",
			key:     "errors",
			errs:    []error{stderrors.New("error1"), stderrors.New("error2")},
			wantKey: "errors",
		},
		{
			name:    "given_nil_error_when_slice_to_zerolog_then_has_key",
			key:     "errors",
			errs:    []error{nil},
			wantKey: "errors",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				sliceToZerolog(event, test.key, test.errs)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)
				assert.Contains(t, result, test.wantKey)
			},
		)
	}
}

func TestSliceToZerologWithAttrs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		wantKey string
		attrs   []Attr
	}{
		{
			name:    "given_empty_attrs_slice_when_slice_to_zerolog_then_has_key",
			key:     "attrs",
			attrs:   []Attr{},
			wantKey: "attrs",
		},
		{
			name:    "given_single_attr_when_slice_to_zerolog_then_has_key",
			key:     "attrs",
			attrs:   []Attr{String("key", "value")},
			wantKey: "attrs",
		},
		{
			name:    "given_multiple_attrs_when_slice_to_zerolog_then_has_key",
			key:     "attrs",
			attrs:   []Attr{String("key1", "value1"), Int("key2", 42)},
			wantKey: "attrs",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				sliceToZerolog(event, test.key, test.attrs)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)
				assert.Contains(t, result, test.wantKey)
			},
		)
	}
}

func TestStructuredErrorMarshalZerologObjectIntegration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		err *StructuredError
		// then
		wantKeys []string
	}{
		{
			name:     "given_simple_error_when_marshal_zerolog_object_then_has_message",
			err:      New("test"),
			wantKeys: []string{"message"},
		},
		{
			name:     "given_error_with_tags_when_marshal_zerolog_object_then_has_message_and_tags",
			err:      New("test").WithTags("tag1"),
			wantKeys: []string{"message", "tags"},
		},
		{
			name: "given_complex_error_when_marshal_zerolog_object_then_has_all_fields",
			err: New("parent").
				WithTags("api").
				WithAttrs(String("request_id", "123")).
				WithErrors(stderrors.New("child")).
				WithStack([]byte("stack")),
			wantKeys: []string{"message", "tags", "attrs", "errors", "stack"},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var buf bytes.Buffer

				logger := zerolog.New(&buf)
				event := logger.Info()

				// when
				test.err.MarshalZerologObject(event)
				event.Msg("test")

				// then
				var result map[string]any

				err := json.Unmarshal(buf.Bytes(), &result)
				require.NoError(t, err)

				for _, key := range test.wantKeys {
					assert.Contains(t, result, key)
				}
			},
		)
	}
}
