package errors

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructuredErrorMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		err *StructuredError
		// then
		wantContains []string
		wantErr      bool
	}{
		{
			name:         "given_nil_error_when_marshal_json_then_returns_json_with_nil_message",
			err:          nil,
			wantContains: []string{`"message":"!NILVALUE"`},
			wantErr:      false,
		},
		{
			name:         "given_error_with_message_when_marshal_json_then_returns_json_with_message",
			err:          New("test error"),
			wantContains: []string{`"message":"test error"`},
			wantErr:      false,
		},
		{
			name:         "given_error_with_empty_message_when_marshal_json_then_returns_json_with_nil_value",
			err:          New(""),
			wantContains: []string{`"message":"!NILVALUE"`},
			wantErr:      false,
		},
		{
			name:         "given_error_with_tags_when_marshal_json_then_returns_json_with_tags",
			err:          New("test").WithTags("tag1", "tag2"),
			wantContains: []string{`"message":"test"`, `"tags":[`, `"tag1"`, `"tag2"`},
			wantErr:      false,
		},
		{
			name:         "given_error_with_attrs_when_marshal_json_then_returns_json_with_attrs",
			err:          New("test").WithAttrs(String("key", "value")),
			wantContains: []string{`"message":"test"`, `"attrs":`},
			wantErr:      false,
		},
		{
			name:         "given_error_with_child_errors_when_marshal_json_then_returns_json_with_errors",
			err:          New("parent").WithErrors(stderrors.New("child")),
			wantContains: []string{`"message":"parent"`, `"errors":[`},
			wantErr:      false,
		},
		{
			name:         "given_error_with_stack_when_marshal_json_then_returns_json_with_base64_stack",
			err:          New("test").WithStack([]byte("stack trace")),
			wantContains: []string{`"message":"test"`, `"stack":"`},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got, err := test.err.MarshalJSON()

				// then
				if test.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					gotStr := string(got)
					for _, want := range test.wantContains {
						assert.Contains(t, gotStr, want)
					}
				}
			},
		)
	}
}

func TestStructuredErrorUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		jsonData string
		// then
		wantMessage string
		wantErr     bool
	}{
		{
			name:        "given_json_with_message_when_unmarshal_json_then_sets_message",
			jsonData:    `{"message":"test error"}`,
			wantMessage: "test error",
			wantErr:     false,
		},
		{
			name:        "given_json_with_empty_message_when_unmarshal_json_then_sets_empty_message",
			jsonData:    `{"message":""}`,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:        "given_empty_json_when_unmarshal_json_then_no_error",
			jsonData:    `{}`,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name:        "given_invalid_json_when_unmarshal_json_then_returns_error",
			jsonData:    `{invalid}`,
			wantMessage: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var err StructuredError

				// when
				gotErr := err.UnmarshalJSON([]byte(test.jsonData))

				// then
				if test.wantErr {
					require.Error(t, gotErr)
				} else {
					require.NoError(t, gotErr)
					assert.Equal(t, test.wantMessage, err.Message)
				}
			},
		)
	}
}

func TestStructuredErrorUnmarshalJSONWithFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		jsonData string
		// then
		wantTagsLen   int
		wantAttrsLen  int
		wantErrorsLen int
		wantHasStack  bool
	}{
		{
			name:          "given_json_with_tags_when_unmarshal_json_then_sets_tags",
			jsonData:      `{"message":"test","tags":["tag1","tag2"]}`,
			wantTagsLen:   2,
			wantAttrsLen:  0,
			wantErrorsLen: 0,
			wantHasStack:  false,
		},
		{
			name:          "given_json_with_attrs_when_unmarshal_json_then_sets_attrs",
			jsonData:      `{"message":"test","attrs":[{"Type":1,"Key":"key","Value":"value"}]}`,
			wantTagsLen:   0,
			wantAttrsLen:  1,
			wantErrorsLen: 0,
			wantHasStack:  false,
		},
		{
			name:          "given_json_with_errors_when_unmarshal_json_then_sets_errors",
			jsonData:      `{"message":"parent","errors":[{"message":"child"}]}`,
			wantTagsLen:   0,
			wantAttrsLen:  0,
			wantErrorsLen: 1,
			wantHasStack:  false,
		},
		{
			name:          "given_json_with_stack_when_unmarshal_json_then_sets_stack",
			jsonData:      `{"message":"test","stack":"c3RhY2s="}`,
			wantTagsLen:   0,
			wantAttrsLen:  0,
			wantErrorsLen: 0,
			wantHasStack:  true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var err StructuredError

				// when
				gotErr := err.UnmarshalJSON([]byte(test.jsonData))

				// then
				require.NoError(t, gotErr)
				assert.Len(t, err.Tags, test.wantTagsLen)
				assert.Len(t, err.Attrs, test.wantAttrsLen)
				assert.Len(t, err.Errors, test.wantErrorsLen)

				if test.wantHasStack {
					assert.NotEmpty(t, err.Stack)
				} else {
					assert.Empty(t, err.Stack)
				}
			},
		)
	}
}

func TestStructuredErrorJSONRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  *StructuredError
		name string
	}{
		{
			name: "given_simple_error_when_marshal_unmarshal_then_preserves_message",
			err:  New("test error"),
		},
		{
			name: "given_error_with_tags_when_marshal_unmarshal_then_preserves_tags",
			err:  New("test").WithTags("tag1", "tag2"),
		},
		{
			name: "given_error_with_attrs_when_marshal_unmarshal_then_preserves_attrs",
			err:  New("test").WithAttrs(String("key", "value"), Int("count", 42)),
		},
		{
			name: "given_complex_error_when_marshal_unmarshal_then_preserves_all_fields",
			err: New("parent").
				WithTags("api", "error").
				WithAttrs(String("request_id", "123")).
				WithErrors(stderrors.New("child error")).
				WithStack([]byte("stack trace")),
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				jsonData, err := json.Marshal(test.err)
				require.NoError(t, err)

				var unmarshaled StructuredError

				err = json.Unmarshal(jsonData, &unmarshaled)

				// then
				require.NoError(t, err)
				assert.Equal(t, test.err.Message, unmarshaled.Message)
				assert.Len(t, unmarshaled.Tags, len(test.err.Tags))
				assert.Len(t, unmarshaled.Attrs, len(test.err.Attrs))
				assert.Len(t, unmarshaled.Errors, len(test.err.Errors))
			},
		)
	}
}

func TestValueToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		key   string
		value string
		// then
		want string
	}{
		{
			name:  "given_key_value_when_value_to_json_then_returns_json_key_value",
			key:   "message",
			value: "test",
			want:  `"message":"test"`,
		},
		{
			name:  "given_empty_value_when_value_to_json_then_returns_json_with_empty_value",
			key:   "key",
			value: "",
			want:  `"key":""`,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var bb bytes.Buffer

				// when
				valueToJSON(&bb, test.key, test.value)

				// then
				assert.Equal(t, test.want, bb.String())
			},
		)
	}
}

func TestErrorToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		err error
		// then
		wantContains []string
	}{
		{
			name:         "given_nil_error_when_error_to_json_then_returns_json_with_nil_message",
			err:          nil,
			wantContains: []string{`{"message":"!NILVALUE"}`},
		},
		{
			name:         "given_standard_error_when_error_to_json_then_returns_json_with_message",
			err:          stderrors.New("standard error"),
			wantContains: []string{`{"message":"standard error"}`},
		},
		{
			name:         "given_structured_error_when_error_to_json_then_returns_full_json",
			err:          New("structured"),
			wantContains: []string{`"message":"structured"`},
		},
		{
			name:         "given_error_with_whitespace_when_error_to_json_then_trims_whitespace",
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
				var bb bytes.Buffer

				// when
				errorToJSON(&bb, test.err)

				// then
				got := bb.String()
				for _, want := range test.wantContains {
					assert.Contains(t, got, want)
				}
			},
		)
	}
}

func TestSliceToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		key   string
		slice []string
		// then
		wantContains []string
	}{
		{
			name:         "given_empty_slice_when_slice_to_json_then_returns_empty_array",
			key:          "tags",
			slice:        []string{},
			wantContains: []string{`"tags":[]`},
		},
		{
			name:         "given_single_item_slice_when_slice_to_json_then_returns_array_with_item",
			key:          "tags",
			slice:        []string{"tag1"},
			wantContains: []string{`"tags":`, `"tag1"`},
		},
		{
			name:         "given_multiple_items_slice_when_slice_to_json_then_returns_array_with_all_items",
			key:          "tags",
			slice:        []string{"tag1", "tag2", "tag3"},
			wantContains: []string{`"tags":`, `"tag1"`, `"tag2"`, `"tag3"`},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var bb bytes.Buffer

				// when
				sliceToJSON(&bb, test.key, test.slice)

				// then
				got := bb.String()
				for _, want := range test.wantContains {
					assert.Contains(t, got, want)
				}
			},
		)
	}
}

func TestSliceToJSONWithErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		key  string
		errs []error
		// then
		wantContains []string
	}{
		{
			name:         "given_empty_errors_slice_when_slice_to_json_then_returns_empty_array",
			key:          "errors",
			errs:         []error{},
			wantContains: []string{`"errors":[]`},
		},
		{
			name:         "given_single_error_when_slice_to_json_then_returns_array_with_error",
			key:          "errors",
			errs:         []error{stderrors.New("error1")},
			wantContains: []string{`"errors":`, `"message":"error1"`},
		},
		{
			name:         "given_multiple_errors_when_slice_to_json_then_returns_array_with_all_errors",
			key:          "errors",
			errs:         []error{stderrors.New("error1"), stderrors.New("error2")},
			wantContains: []string{`"errors":`, `"message":"error1"`, `"message":"error2"`},
		},
		{
			name:         "given_nil_error_when_slice_to_json_then_includes_nil_message",
			key:          "errors",
			errs:         []error{nil},
			wantContains: []string{`"errors":`, `"message":"!NILVALUE"`},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				var bb bytes.Buffer

				// when
				sliceToJSON(&bb, test.key, test.errs)

				// then
				got := bb.String()
				for _, want := range test.wantContains {
					assert.Contains(t, got, want)
				}
			},
		)
	}
}
