package errors

import (
	stderrors "errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStructuredErrorLogValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err      *StructuredError
		name     string
		wantKind slog.Kind
	}{
		{
			name:     "given_nil_error_when_log_value_then_returns_group_value",
			err:      nil,
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_error_with_message_when_log_value_then_returns_group_value",
			err:      New("test error"),
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_error_with_tags_when_log_value_then_returns_group_value",
			err:      New("test").WithTags("tag1", "tag2"),
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_error_with_attrs_when_log_value_then_returns_group_value",
			err:      New("test").WithAttrs(String("key", "value")),
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_error_with_errors_when_log_value_then_returns_group_value",
			err:      New("parent").WithErrors(stderrors.New("child")),
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_error_with_stack_when_log_value_then_returns_group_value",
			err:      New("test").WithStack([]byte("stack trace")),
			wantKind: slog.KindGroup,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.err.LogValue()

				// then
				assert.Equal(t, test.wantKind, got.Kind())
			},
		)
	}
}

func TestStructuredErrorLogValueAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err           *StructuredError
		name          string
		wantAttrCount int
	}{
		{
			name:          "given_nil_error_when_log_value_then_returns_one_attr",
			err:           nil,
			wantAttrCount: 1,
		},
		{
			name:          "given_error_with_only_message_when_log_value_then_returns_one_attr",
			err:           New("test"),
			wantAttrCount: 1,
		},
		{
			name:          "given_error_with_message_and_tags_when_log_value_then_returns_two_attrs",
			err:           New("test").WithTags("tag1"),
			wantAttrCount: 2,
		},
		{
			name:          "given_error_with_message_and_attrs_when_log_value_then_returns_two_attrs",
			err:           New("test").WithAttrs(String("key", "value")),
			wantAttrCount: 2,
		},
		{
			name: "given_error_with_all_fields_when_log_value_then_returns_five_attrs",
			err: New("test").WithTags("tag").WithAttrs(
				String(
					"key",
					"value",
				),
			).WithErrors(stderrors.New("child")).WithStack([]byte("stack")),
			wantAttrCount: 5,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.err.LogValue()

				// then
				attrs := got.Group()
				assert.Len(t, attrs, test.wantAttrCount)
			},
		)
	}
}

func TestAttrLogValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		attr     *Attr
		name     string
		wantKind slog.Kind
	}{
		{
			name:     "given_nil_attr_when_log_value_then_returns_group_value",
			attr:     nil,
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_string_attr_when_log_value_then_returns_group_value",
			attr:     &Attr{Type: StringType, Key: "name", Value: "test"},
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_int_attr_when_log_value_then_returns_group_value",
			attr:     &Attr{Type: IntType, Key: "count", Value: 42},
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_bool_attr_when_log_value_then_returns_group_value",
			attr:     &Attr{Type: BoolType, Key: "active", Value: true},
			wantKind: slog.KindGroup,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.attr.LogValue()

				// then
				assert.Equal(t, test.wantKind, got.Kind())
			},
		)
	}
}

func TestAttrAsSlog(t *testing.T) {
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
			name:    "given_nil_attr_when_as_slog_then_returns_attr_with_nil_key",
			attr:    nil,
			wantKey: nilValue,
		},
		{
			name:    "given_string_attr_when_as_slog_then_returns_string_attr",
			attr:    &Attr{Type: StringType, Key: "name", Value: "test"},
			wantKey: "name",
		},
		{
			name:    "given_int_attr_when_as_slog_then_returns_int_attr",
			attr:    &Attr{Type: IntType, Key: "count", Value: 42},
			wantKey: "count",
		},
		{
			name:    "given_int64_attr_when_as_slog_then_returns_int64_attr",
			attr:    &Attr{Type: Int64Type, Key: "id", Value: int64(123)},
			wantKey: "id",
		},
		{
			name:    "given_uint64_attr_when_as_slog_then_returns_uint64_attr",
			attr:    &Attr{Type: Uint64Type, Key: "uid", Value: uint64(456)},
			wantKey: "uid",
		},
		{
			name:    "given_float64_attr_when_as_slog_then_returns_float64_attr",
			attr:    &Attr{Type: Float64Type, Key: "price", Value: 99.99},
			wantKey: "price",
		},
		{
			name:    "given_bool_attr_when_as_slog_then_returns_bool_attr",
			attr:    &Attr{Type: BoolType, Key: "active", Value: true},
			wantKey: "active",
		},
		{
			name:    "given_time_attr_when_as_slog_then_returns_time_attr",
			attr:    &Attr{Type: TimeType, Key: "created", Value: fixedTime},
			wantKey: "created",
		},
		{
			name:    "given_duration_attr_when_as_slog_then_returns_duration_attr",
			attr:    &Attr{Type: DurationType, Key: "timeout", Value: fixedDuration},
			wantKey: "timeout",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.attr.asSlog()

				// then
				assert.Equal(t, test.wantKey, got.Key)
			},
		)
	}
}

func TestAttrAsSlogWithSlices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantKey string
	}{
		{
			name:    "given_bools_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: BoolsType, Key: "flags", Value: []bool{true, false}},
			wantKey: "flags",
		},
		{
			name:    "given_ints_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: IntsType, Key: "numbers", Value: []int{1, 2, 3}},
			wantKey: "numbers",
		},
		{
			name:    "given_int64s_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: Int64sType, Key: "ids", Value: []int64{1, 2, 3}},
			wantKey: "ids",
		},
		{
			name:    "given_uint64s_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: Uint64sType, Key: "uids", Value: []uint64{1, 2, 3}},
			wantKey: "uids",
		},
		{
			name:    "given_float64s_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: Float64sType, Key: "prices", Value: []float64{1.1, 2.2}},
			wantKey: "prices",
		},
		{
			name:    "given_strings_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: StringsType, Key: "tags", Value: []string{"tag1", "tag2"}},
			wantKey: "tags",
		},
		{
			name:    "given_times_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: TimesType, Key: "timestamps", Value: []time.Time{time.Now()}},
			wantKey: "timestamps",
		},
		{
			name:    "given_durations_attr_when_as_slog_then_returns_group_attr",
			attr:    &Attr{Type: DurationsType, Key: "timeouts", Value: []time.Duration{5 * time.Second}},
			wantKey: "timeouts",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.attr.asSlog()

				// then
				assert.Equal(t, test.wantKey, got.Key)
				assert.Equal(t, slog.KindGroup, got.Value.Kind())
			},
		)
	}
}

func TestAttrAsSlogWithObjectType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantKey  string
		wantKind slog.Kind
	}{
		{
			name:     "given_object_attr_when_as_slog_then_returns_group_attr",
			attr:     &Attr{Type: ObjectType, Key: "obj", Value: []Attr{String("key", "value")}},
			wantKey:  "obj",
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_empty_object_attr_when_as_slog_then_returns_group_attr",
			attr:     &Attr{Type: ObjectType, Key: "empty", Value: []Attr{}},
			wantKey:  "empty",
			wantKind: slog.KindGroup,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.attr.asSlog()

				// then
				assert.Equal(t, test.wantKey, got.Key)
				assert.Equal(t, test.wantKind, got.Value.Kind())
			},
		)
	}
}

func TestAttrAsSlogWithAnyType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		attr *Attr
		// then
		wantKey string
	}{
		{
			name:    "given_any_type_attr_when_as_slog_then_returns_any_attr",
			attr:    &Attr{Type: AnyType, Key: "custom", Value: map[string]string{"key": "value"}},
			wantKey: "custom",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := test.attr.asSlog()

				// then
				assert.Equal(t, test.wantKey, got.Key)
			},
		)
	}
}

func TestErrorToSlog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		key string
		err error
		// then
		wantKey  string
		wantKind slog.Kind
	}{
		{
			name:     "given_nil_error_when_error_to_slog_then_returns_group_with_nil_message",
			key:      "error",
			err:      nil,
			wantKey:  "error",
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_standard_error_when_error_to_slog_then_returns_group_with_message",
			key:      "error",
			err:      stderrors.New("standard error"),
			wantKey:  "error",
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_structured_error_when_error_to_slog_then_returns_group",
			key:      "error",
			err:      New("structured"),
			wantKey:  "error",
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_error_with_whitespace_when_error_to_slog_then_trims_whitespace",
			key:      "error",
			err:      stderrors.New("  error  "),
			wantKey:  "error",
			wantKind: slog.KindGroup,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := errorToSlog(test.key, test.err)

				// then
				assert.Equal(t, test.wantKey, got.Key)
				assert.Equal(t, test.wantKind, got.Value.Kind())
			},
		)
	}
}

func TestSliceToSlog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      string
		wantKey  string
		slice    []string
		wantKind slog.Kind
	}{
		{
			name:     "given_empty_slice_when_slice_to_slog_then_returns_empty_group",
			key:      "tags",
			slice:    []string{},
			wantKey:  "tags",
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_single_item_slice_when_slice_to_slog_then_returns_group_with_item",
			key:      "tags",
			slice:    []string{"tag1"},
			wantKey:  "tags",
			wantKind: slog.KindGroup,
		},
		{
			name:     "given_multiple_items_slice_when_slice_to_slog_then_returns_group_with_all_items",
			key:      "tags",
			slice:    []string{"tag1", "tag2", "tag3"},
			wantKey:  "tags",
			wantKind: slog.KindGroup,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := sliceToSlog(test.key, test.slice)

				// then
				assert.Equal(t, test.wantKey, got.Key)
				assert.Equal(t, test.wantKind, got.Value.Kind())
			},
		)
	}
}

func TestSliceToSlogWithErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		key           string
		wantKey       string
		errs          []error
		wantAttrCount int
	}{
		{
			name:          "given_empty_errors_slice_when_slice_to_slog_then_returns_empty_group",
			key:           "errors",
			errs:          []error{},
			wantKey:       "errors",
			wantAttrCount: 0,
		},
		{
			name:          "given_single_error_when_slice_to_slog_then_returns_group_with_error",
			key:           "errors",
			errs:          []error{stderrors.New("error1")},
			wantKey:       "errors",
			wantAttrCount: 1,
		},
		{
			name:          "given_multiple_errors_when_slice_to_slog_then_returns_group_with_all_errors",
			key:           "errors",
			errs:          []error{stderrors.New("error1"), stderrors.New("error2")},
			wantKey:       "errors",
			wantAttrCount: 2,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := sliceToSlog(test.key, test.errs)

				// then
				assert.Equal(t, test.wantKey, got.Key)
				attrs := got.Value.Group()
				assert.Len(t, attrs, test.wantAttrCount)
			},
		)
	}
}

func TestSliceToSlogWithAttrs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		key           string
		wantKey       string
		attrs         []Attr
		wantAttrCount int
	}{
		{
			name:          "given_empty_attrs_slice_when_slice_to_slog_then_returns_empty_group",
			key:           "attrs",
			attrs:         []Attr{},
			wantKey:       "attrs",
			wantAttrCount: 0,
		},
		{
			name:          "given_single_attr_when_slice_to_slog_then_returns_group_with_attr",
			key:           "attrs",
			attrs:         []Attr{String("key", "value")},
			wantKey:       "attrs",
			wantAttrCount: 1,
		},
		{
			name:          "given_multiple_attrs_when_slice_to_slog_then_returns_group_with_all_attrs",
			key:           "attrs",
			attrs:         []Attr{String("key1", "value1"), Int("key2", 42)},
			wantKey:       "attrs",
			wantAttrCount: 2,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := sliceToSlog(test.key, test.attrs)

				// then
				assert.Equal(t, test.wantKey, got.Key)
				attrs := got.Value.Group()
				assert.Len(t, attrs, test.wantAttrCount)
			},
		)
	}
}

func TestSliceToSlogWithTypedSlices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		key   string
		slice any
		// then
		wantKey   string
		wantCount int
	}{
		{
			name:      "given_bool_slice_when_slice_to_slog_then_returns_group",
			key:       "flags",
			slice:     []bool{true, false},
			wantKey:   "flags",
			wantCount: 2,
		},
		{
			name:      "given_int_slice_when_slice_to_slog_then_returns_group",
			key:       "numbers",
			slice:     []int{1, 2, 3},
			wantKey:   "numbers",
			wantCount: 3,
		},
		{
			name:      "given_int64_slice_when_slice_to_slog_then_returns_group",
			key:       "ids",
			slice:     []int64{1, 2, 3},
			wantKey:   "ids",
			wantCount: 3,
		},
		{
			name:      "given_uint64_slice_when_slice_to_slog_then_returns_group",
			key:       "uids",
			slice:     []uint64{1, 2, 3},
			wantKey:   "uids",
			wantCount: 3,
		},
		{
			name:      "given_float64_slice_when_slice_to_slog_then_returns_group",
			key:       "prices",
			slice:     []float64{1.1, 2.2},
			wantKey:   "prices",
			wantCount: 2,
		},
		{
			name:      "given_time_slice_when_slice_to_slog_then_returns_group",
			key:       "times",
			slice:     []time.Time{time.Now()},
			wantKey:   "times",
			wantCount: 1,
		},
		{
			name:      "given_duration_slice_when_slice_to_slog_then_returns_group",
			key:       "durations",
			slice:     []time.Duration{5 * time.Second},
			wantKey:   "durations",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				var got slog.Attr

				switch v := test.slice.(type) {
				case []bool:
					got = sliceToSlog(test.key, v)
				case []int:
					got = sliceToSlog(test.key, v)
				case []int64:
					got = sliceToSlog(test.key, v)
				case []uint64:
					got = sliceToSlog(test.key, v)
				case []float64:
					got = sliceToSlog(test.key, v)
				case []time.Time:
					got = sliceToSlog(test.key, v)
				case []time.Duration:
					got = sliceToSlog(test.key, v)
				}

				// then
				assert.Equal(t, test.wantKey, got.Key)
				attrs := got.Value.Group()
				assert.Len(t, attrs, test.wantCount)
			},
		)
	}
}
