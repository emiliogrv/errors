package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		errs        []error
		wantErrsLen int
		wantNil     bool
		wantJoined  bool
	}{
		{
			name:        "given_no_errors_when_join_then_returns_nil",
			errs:        []error{},
			wantNil:     true,
			wantErrsLen: 0,
			wantJoined:  false,
		},
		{
			name:        "given_all_nil_errors_when_join_then_returns_nil",
			errs:        []error{nil, nil, nil},
			wantNil:     true,
			wantErrsLen: 0,
			wantJoined:  false,
		},
		{
			name:        "given_single_non_nil_error_when_join_then_returns_joined_error",
			errs:        []error{stderrors.New("error1")},
			wantNil:     false,
			wantErrsLen: 1,
			wantJoined:  true,
		},
		{
			name:        "given_multiple_non_nil_errors_when_join_then_returns_joined_error_with_all",
			errs:        []error{stderrors.New("error1"), stderrors.New("error2"), stderrors.New("error3")},
			wantNil:     false,
			wantErrsLen: 3,
			wantJoined:  true,
		},
		{
			name:        "given_mixed_nil_and_non_nil_errors_when_join_then_returns_joined_error_without_nils",
			errs:        []error{stderrors.New("error1"), nil, stderrors.New("error2"), nil},
			wantNil:     false,
			wantErrsLen: 2,
			wantJoined:  true,
		},
		{
			name:        "given_structured_errors_when_join_then_returns_joined_error",
			errs:        []error{New("structured1"), New("structured2")},
			wantNil:     false,
			wantErrsLen: 2,
			wantJoined:  true,
		},
		{
			name:        "given_mixed_error_types_when_join_then_returns_joined_error_with_all",
			errs:        []error{stderrors.New("standard"), New("structured"), nil},
			wantNil:     false,
			wantErrsLen: 2,
			wantJoined:  true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Join(test.errs...)

				// then
				if test.wantNil {
					require.NoError(t, got)
				} else {
					require.Error(t, got)

					structErr := &StructuredError{}
					ok := stderrors.As(got, &structErr)
					assert.True(t, ok)
					assert.Equal(t, test.wantJoined, structErr.joined)
					assert.Len(t, structErr.Errors, test.wantErrsLen)
				}
			},
		)
	}
}

func TestJoinUnwrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		errs []error
		// then
		wantUnwrapLen int
	}{
		{
			name:          "given_joined_error_when_unwrap_then_returns_all_errors",
			errs:          []error{stderrors.New("err1"), stderrors.New("err2")},
			wantUnwrapLen: 2,
		},
		{
			name:          "given_joined_error_with_single_error_when_unwrap_then_returns_single_error",
			errs:          []error{stderrors.New("only")},
			wantUnwrapLen: 1,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// given
				joined := Join(test.errs...)

				// when
				unwrapper, ok := joined.(MultiUnwrapper)

				// then
				assert.True(t, ok)

				unwrapped := unwrapper.Unwrap()
				assert.Len(t, unwrapped, test.wantUnwrapLen)
			},
		)
	}
}

func TestJoinIf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		errs        []error
		wantErrsLen int
		wantNil     bool
		wantJoined  bool
	}{
		{
			name:        "given_empty_slice_when_join_if_then_returns_nil",
			errs:        []error{},
			wantNil:     true,
			wantErrsLen: 0,
			wantJoined:  false,
		},
		{
			name:        "given_first_error_nil_when_join_if_then_returns_nil",
			errs:        []error{nil, stderrors.New("error2"), stderrors.New("error3")},
			wantNil:     true,
			wantErrsLen: 0,
			wantJoined:  false,
		},
		{
			name:        "given_first_error_non_nil_when_join_if_then_returns_joined_error",
			errs:        []error{stderrors.New("error1"), stderrors.New("error2")},
			wantNil:     false,
			wantErrsLen: 2,
			wantJoined:  true,
		},
		{
			name:        "given_first_error_non_nil_and_rest_nil_when_join_if_then_returns_joined_error_with_first",
			errs:        []error{stderrors.New("error1"), nil, nil},
			wantNil:     false,
			wantErrsLen: 1,
			wantJoined:  true,
		},
		{
			name:        "given_first_error_non_nil_and_mixed_rest_when_join_if_then_returns_joined_error_without_nils",
			errs:        []error{stderrors.New("error1"), nil, stderrors.New("error2")},
			wantNil:     false,
			wantErrsLen: 2,
			wantJoined:  true,
		},
		{
			name:        "given_single_non_nil_error_when_join_if_then_returns_joined_error",
			errs:        []error{stderrors.New("only")},
			wantNil:     false,
			wantErrsLen: 1,
			wantJoined:  true,
		},
		{
			name:        "given_single_nil_error_when_join_if_then_returns_nil",
			errs:        []error{nil},
			wantNil:     true,
			wantErrsLen: 0,
			wantJoined:  false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := JoinIf(test.errs...)

				// then
				if test.wantNil {
					require.NoError(t, got)
				} else {
					require.Error(t, got)

					structErr := &StructuredError{}
					ok := stderrors.As(got, &structErr)
					assert.True(t, ok)
					assert.Equal(t, test.wantJoined, structErr.joined)
					assert.Len(t, structErr.Errors, test.wantErrsLen)
				}
			},
		)
	}
}

func TestJoinIfBehaviorDifference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		errs []error
		// then
		wantJoinNil   bool
		wantJoinIfNil bool
	}{
		{
			name:          "given_first_nil_rest_non_nil_when_comparing_join_and_join_if_then_different_results",
			errs:          []error{nil, stderrors.New("error2")},
			wantJoinNil:   false, // Join returns error because there's a non-nil error
			wantJoinIfNil: true,  // JoinIf returns nil because first is nil
		},
		{
			name:          "given_all_non_nil_when_comparing_join_and_join_if_then_same_results",
			errs:          []error{stderrors.New("error1"), stderrors.New("error2")},
			wantJoinNil:   false,
			wantJoinIfNil: false,
		},
		{
			name:          "given_all_nil_when_comparing_join_and_join_if_then_same_results",
			errs:          []error{nil, nil},
			wantJoinNil:   true,
			wantJoinIfNil: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				gotJoin := Join(test.errs...)
				gotJoinIf := JoinIf(test.errs...)

				// then
				if test.wantJoinNil {
					require.NoError(t, gotJoin)
				} else {
					require.Error(t, gotJoin)
				}

				if test.wantJoinIfNil {
					require.NoError(t, gotJoinIf)
				} else {
					require.Error(t, gotJoinIf)
				}
			},
		)
	}
}

func TestJoinWithStructuredErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		errs []error
		// then
		wantErrsLen int
	}{
		{
			name: "given_structured_errors_with_attrs_when_join_then_preserves_attrs",
			errs: []error{
				New("error1").WithAttrs(String("key1", "value1")),
				New("error2").WithAttrs(String("key2", "value2")),
			},
			wantErrsLen: 2,
		},
		{
			name: "given_structured_errors_with_nested_errors_when_join_then_preserves_nested",
			errs: []error{
				New("parent1").WithErrors(stderrors.New("child1")),
				New("parent2").WithErrors(stderrors.New("child2")),
			},
			wantErrsLen: 2,
		},
		{
			name: "given_structured_errors_with_tags_when_join_then_preserves_tags",
			errs: []error{
				New("error1").WithTags("tag1"),
				New("error2").WithTags("tag2"),
			},
			wantErrsLen: 2,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				got := Join(test.errs...)

				// then
				require.Error(t, got)

				structErr := &StructuredError{}
				ok := stderrors.As(got, &structErr)
				assert.True(t, ok)
				assert.Len(t, structErr.Errors, test.wantErrsLen)

				// Verify original errors are preserved
				for i, err := range test.errs {
					assert.Equal(t, err, structErr.Errors[i])
				}
			},
		)
	}
}
