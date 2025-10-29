package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompatibilityWithStdErrors tests that this package can be used as a drop-in
// replacement for the standard errors package by comparing string outputs.
// Note: The structured error format produces different string output than standard errors,
// but maintains full API compatibility for errors.Is, errors.As, and errors.Unwrap.
func TestCompatibilityWithStdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		buildStdError      func() error
		buildCustomError   func() error
		name               string
		note               string
		expectedTextInBoth []string
		expectedTextInStd  []string
		expectEqual        bool
	}{
		{
			name: "given_simple_error_when_comparing_new_then_strings_differ_due_to_format",
			buildStdError: func() error {
				return stderrors.New("simple error")
			},
			buildCustomError: func() error {
				return New("simple error")
			},
			expectEqual:        false,
			note:               "Custom errors use structured format: (message=simple error)",
			expectedTextInBoth: []string{"simple error"},
		},
		{
			name: "given_empty_error_when_comparing_new_then_strings_differ",
			buildStdError: func() error {
				return stderrors.New("")
			},
			buildCustomError: func() error {
				return New("")
			},
			expectEqual: false,
			note:        "Empty message becomes (message=!NILVALUE) in custom format",
		},
		{
			name: "given_wrapped_error_when_using_fmt_errorf_then_strings_differ",
			buildStdError: func() error {
				base := stderrors.New("base error")

				return fmt.Errorf("wrapped: %w", base)
			},
			buildCustomError: func() error {
				base := New("base error")

				return fmt.Errorf("wrapped: %w", base)
			},
			expectEqual:        false,
			note:               "Custom error wrapped in fmt.Errorf shows structured format",
			expectedTextInBoth: []string{"wrapped:", "base error"},
		},
		{
			name: "given_multiple_wrapped_errors_when_using_fmt_errorf_then_strings_differ",
			buildStdError: func() error {
				err1 := stderrors.New("error 1")
				err2 := stderrors.New("error 2")

				return fmt.Errorf("%w, %w", err1, err2)
			},
			buildCustomError: func() error {
				err1 := New("error 1")
				err2 := New("error 2")

				return fmt.Errorf("%w, %w", err1, err2)
			},
			expectEqual:        false,
			note:               "Custom errors in fmt.Errorf show structured format",
			expectedTextInBoth: []string{"error 1", "error 2"},
		},
		{
			name: "given_joined_errors_when_using_join_then_strings_differ",
			buildStdError: func() error {
				err1 := stderrors.New("error 1")
				err2 := stderrors.New("error 2")
				err3 := stderrors.New("error 3")

				return stderrors.Join(err1, err2, err3)
			},
			buildCustomError: func() error {
				err1 := New("error 1")
				err2 := New("error 2")
				err3 := New("error 3")

				return Join(err1, err2, err3)
			},
			expectEqual:        false,
			note:               "Join produces structured format with errors array",
			expectedTextInBoth: []string{"error 1", "error 2", "error 3"},
		},
		{
			name: "given_joined_errors_with_nil_when_using_join_then_strings_differ",
			buildStdError: func() error {
				err1 := stderrors.New("error 1")
				err2 := stderrors.New("error 2")

				return stderrors.Join(err1, nil, err2)
			},
			buildCustomError: func() error {
				err1 := New("error 1")
				err2 := New("error 2")

				return Join(err1, nil, err2)
			},
			expectEqual:        false,
			note:               "Nil errors are filtered but format differs",
			expectedTextInBoth: []string{"error 1", "error 2"},
		},
		{
			name: "given_nested_joined_errors_when_using_join_then_strings_differ",
			buildStdError: func() error {
				inner := stderrors.Join(
					stderrors.New("inner 1"),
					stderrors.New("inner 2"),
				)

				return stderrors.Join(
					stderrors.New("outer"),
					inner,
				)
			},
			buildCustomError: func() error {
				inner := Join(
					New("inner 1"),
					New("inner 2"),
				)

				return Join(
					New("outer"),
					inner,
				)
			},
			expectEqual:        false,
			note:               "Nested joins show structured format hierarchy",
			expectedTextInBoth: []string{"outer", "inner 1", "inner 2"},
		},
		{
			name: "given_wrapped_custom_error_when_using_with_errors_then_output_differs",
			buildStdError: func() error {
				base := stderrors.New("base error")

				return fmt.Errorf("wrapper: %w", base)
			},
			buildCustomError: func() error {
				base := New("base error")

				return New("wrapper").WithErrors(base)
			},
			expectEqual:        false,
			note:               "WithErrors produces different format than fmt.Errorf",
			expectedTextInBoth: []string{"wrapper", "base error"},
			expectedTextInStd:  []string{"wrapper:"},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				stdErr := test.buildStdError()
				customErr := test.buildCustomError()

				stdStr := ""
				if stdErr != nil {
					stdStr = stdErr.Error()
				}

				customStr := ""
				if customErr != nil {
					customStr = customErr.Error()
				}

				// then - verify format difference
				if test.expectEqual {
					assert.Equal(t, stdStr, customStr, "Error strings should match. Note: %s", test.note)
				} else {
					assert.NotEqual(t, stdStr, customStr, "Error strings should differ. Note: %s", test.note)
				}

				// then - verify expected text appears in both
				for _, expectedText := range test.expectedTextInBoth {
					assert.Contains(t, stdStr, expectedText, "Standard error should contain: %s", expectedText)
					assert.Contains(t, customStr, expectedText, "Custom error should contain: %s", expectedText)
				}

				// then - verify text that only appears in standard errors
				for _, expectedText := range test.expectedTextInStd {
					assert.Contains(t, stdStr, expectedText, "Standard error should contain: %s", expectedText)
				}
			},
		)
	}
}

// TestCompatibilityIs tests that errors.Is works correctly with both
// standard and custom errors.
func TestCompatibilityIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		buildErr    func() error
		buildTarget func() error
		name        string
		wantMatch   bool
	}{
		{
			name: "given_custom_error_when_comparing_with_is_then_matches_itself",
			buildErr: func() error {
				return New("test error")
			},
			buildTarget: func() error {
				return New("test error")
			},
			wantMatch: false, // Different instances
		},
		{
			name: "given_wrapped_custom_error_when_using_is_then_finds_target",
			buildErr: func() error {
				target := New("target error")

				return fmt.Errorf("wrapped: %w", target)
			},
			buildTarget: func() error {
				// This will be a different instance, so it won't match
				return New("target error")
			},
			wantMatch: false,
		},
		{
			name: "given_joined_errors_when_using_is_then_finds_target",
			buildErr: func() error {
				err1 := New("error 1")
				err2 := New("error 2")

				return Join(err1, err2)
			},
			buildTarget: func() error {
				return New("error 1")
			},
			wantMatch: false, // Different instances
		},
		{
			name: "given_std_error_wrapped_in_custom_when_using_is_then_finds_std_target",
			buildErr: func() error {
				stdErr := stderrors.New("std error")

				return New("wrapper").WithErrors(stdErr)
			},
			buildTarget: func() error {
				return stderrors.New("std error")
			},
			wantMatch: false, // Different instances
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				err := test.buildErr()
				target := test.buildTarget()
				got := stderrors.Is(err, target)

				// then
				assert.Equal(t, test.wantMatch, got)
			},
		)
	}
}

// TestCompatibilityIsWithSameInstance tests errors.Is with same error instances.
func TestCompatibilityIsWithSameInstance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		buildErr  func() (error, error)
		name      string
		wantMatch bool
	}{
		{
			name: "given_same_custom_error_instance_when_using_is_then_matches",
			buildErr: func() (error, error) {
				target := New("target error")

				return target, target
			},
			wantMatch: true,
		},
		{
			name: "given_wrapped_same_instance_when_using_is_then_matches",
			buildErr: func() (error, error) {
				target := New("target error")
				wrapped := fmt.Errorf("wrapped: %w", target)

				return wrapped, target
			},
			wantMatch: true,
		},
		{
			name: "given_joined_with_same_instance_when_using_is_then_matches",
			buildErr: func() (error, error) {
				target := New("target error")
				joined := Join(New("other"), target)

				return joined, target
			},
			wantMatch: true,
		},
		{
			name: "given_custom_wrapped_same_instance_when_using_is_then_matches",
			buildErr: func() (error, error) {
				target := New("target error")
				wrapped := New("wrapper").WithErrors(target)

				return wrapped, target
			},
			wantMatch: true,
		},
		{
			name: "given_std_error_wrapped_same_instance_when_using_is_then_matches",
			buildErr: func() (error, error) {
				target := stderrors.New("std error")
				wrapped := New("wrapper").WithErrors(target)

				return wrapped, target
			},
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				err, target := test.buildErr()
				got := stderrors.Is(err, target)

				// then
				assert.Equal(t, test.wantMatch, got)
			},
		)
	}
}

// TestCompatibilityAs tests that errors.As works correctly with both
// standard and custom errors.
func TestCompatibilityAs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		buildErr            func() error
		name                string
		wantStructuredError bool
	}{
		{
			name: "given_custom_error_when_using_as_then_extracts_structured_error",
			buildErr: func() error {
				return New("test error")
			},
			wantStructuredError: true,
		},
		{
			name: "given_wrapped_custom_error_when_using_as_then_extracts_structured_error",
			buildErr: func() error {
				return fmt.Errorf("wrapped: %w", New("base error"))
			},
			wantStructuredError: true,
		},
		{
			name: "given_joined_custom_errors_when_using_as_then_extracts_structured_error",
			buildErr: func() error {
				return Join(New("error 1"), New("error 2"))
			},
			wantStructuredError: true,
		},
		{
			name: "given_std_error_when_using_as_then_does_not_extract_structured_error",
			buildErr: func() error {
				return stderrors.New("std error")
			},
			wantStructuredError: false,
		},
		{
			name: "given_std_error_wrapped_in_custom_when_using_as_then_extracts_structured_error",
			buildErr: func() error {
				return New("wrapper").WithErrors(stderrors.New("std error"))
			},
			wantStructuredError: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				err := test.buildErr()

				var structuredErr *StructuredError

				got := stderrors.As(err, &structuredErr)

				// then
				assert.Equal(t, test.wantStructuredError, got)

				if test.wantStructuredError {
					assert.NotNil(t, structuredErr)
				} else {
					assert.Nil(t, structuredErr)
				}
			},
		)
	}
}

// TestCompatibilityUnwrap tests that errors.Unwrap works correctly.
func TestCompatibilityUnwrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		buildErr      func() error
		name          string
		wantUnwrapNil bool
	}{
		{
			name: "given_simple_custom_error_when_unwrapping_then_returns_nil",
			buildErr: func() error {
				return New("simple error")
			},
			wantUnwrapNil: true,
		},
		{
			name: "given_wrapped_custom_error_when_unwrapping_then_returns_base",
			buildErr: func() error {
				return fmt.Errorf("wrapped: %w", New("base error"))
			},
			wantUnwrapNil: false,
		},
		{
			name: "given_custom_with_errors_when_unwrapping_then_returns_nil",
			buildErr: func() error {
				// Note: Unwrap() returns []error, not error, so standard Unwrap returns nil
				return New("wrapper").WithErrors(New("base error"))
			},
			wantUnwrapNil: true,
		},
		{
			name: "given_std_error_when_unwrapping_then_returns_nil",
			buildErr: func() error {
				return stderrors.New("std error")
			},
			wantUnwrapNil: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				err := test.buildErr()
				unwrapped := stderrors.Unwrap(err)

				// then
				if test.wantUnwrapNil {
					assert.NoError(t, unwrapped)
				} else {
					assert.Error(t, unwrapped)
				}
			},
		)
	}
}

// TestCompatibilityFmtErrorfBehavior tests that fmt.Errorf works correctly with
// both standard and custom errors. Text alongside %w is preserved in fmt.Errorf.
func TestCompatibilityFmtErrorfBehavior(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		buildStdErr    func() error
		buildCustomErr func() error
		// then
		searchText         string
		wantStdContains    bool
		wantCustomContains bool
	}{
		{
			name: "given_text_before_wrapped_error_when_using_fmt_errorf_then_text_is_preserved",
			buildStdErr: func() error {
				base := stderrors.New("base error")

				return fmt.Errorf("text example %w", base)
			},
			buildCustomErr: func() error {
				base := New("base error")

				return fmt.Errorf("text example %w", base)
			},
			searchText:         "text example",
			wantStdContains:    true,
			wantCustomContains: true,
		},
		{
			name: "given_text_between_wrapped_errors_when_using_fmt_errorf_then_text_is_preserved",
			buildStdErr: func() error {
				err1 := stderrors.New("error 1")
				err2 := stderrors.New("error 2")

				return fmt.Errorf("%w text between %w", err1, err2)
			},
			buildCustomErr: func() error {
				err1 := New("error 1")
				err2 := New("error 2")

				return fmt.Errorf("%w text between %w", err1, err2)
			},
			searchText:         "text between",
			wantStdContains:    true,
			wantCustomContains: true,
		},
		{
			name: "given_text_after_wrapped_error_when_using_fmt_errorf_then_text_is_preserved",
			buildStdErr: func() error {
				base := stderrors.New("base error")

				return fmt.Errorf("%w text after", base)
			},
			buildCustomErr: func() error {
				base := New("base error")

				return fmt.Errorf("%w text after", base)
			},
			searchText:         "text after",
			wantStdContains:    true,
			wantCustomContains: true,
		},
		{
			name: "given_only_text_when_using_fmt_errorf_then_text_is_preserved_in_both",
			buildStdErr: func() error {
				return stderrors.New("only text no wrapping")
			},
			buildCustomErr: func() error {
				return New("only text no wrapping")
			},
			searchText:         "only text no wrapping",
			wantStdContains:    true,
			wantCustomContains: true,
		},
		{
			name: "given_text_with_format_verbs_when_using_fmt_errorf_then_text_is_preserved_in_both",
			buildStdErr: func() error {
				return fmt.Errorf("text with %s and %d", "string", 42)
			},
			buildCustomErr: func() error {
				return fmt.Errorf("text with %s and %d", "string", 42)
			},
			searchText:         "text with string and 42",
			wantStdContains:    true,
			wantCustomContains: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				stdErr := test.buildStdErr()
				customErr := test.buildCustomErr()

				stdStr := stdErr.Error()
				customStr := customErr.Error()

				// then - verify both preserve text in fmt.Errorf
				if test.wantStdContains {
					assert.Contains(t, stdStr, test.searchText, "Standard error should contain text")
				} else {
					assert.NotContains(t, stdStr, test.searchText, "Standard error should not contain text")
				}

				if test.wantCustomContains {
					assert.Contains(t, customStr, test.searchText, "Custom error should contain text")
				} else {
					assert.NotContains(t, customStr, test.searchText, "Custom error should not contain text")
				}
			},
		)
	}
}

// TestCompatibilityTextLossInNestedErrors tests the known limitation where
// fmt.Errorf wrapper text is lost when the wrapped error is added to WithErrors().
// This happens because the marshaling extracts the inner StructuredError directly.
func TestCompatibilityTextLossInNestedErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// given
		buildErr func() error
		// then
		searchText       string
		wantContainsText bool
	}{
		{
			name: "given_fmt_errorf_wrapped_custom_error_when_nested_in_with_errors_then_wrapper_text_is_lost",
			buildErr: func() error {
				base := New("base error")
				wrapped := fmt.Errorf("text example %w", base)

				return New("outer").WithErrors(wrapped)
			},
			searchText:       "text example",
			wantContainsText: false,
		},
		{
			name: "given_fmt_errorf_wrapped_custom_error_when_not_nested_then_wrapper_text_is_preserved",
			buildErr: func() error {
				base := New("base error")

				return fmt.Errorf("text example %w", base)
			},
			searchText:       "text example",
			wantContainsText: true,
		},
		{
			name: "given_multiple_fmt_errorf_wrapped_errors_when_nested_in_with_errors_then_all_wrapper_text_is_lost",
			buildErr: func() error {
				err1 := New("error 1")
				wrapped1 := fmt.Errorf("wrapper 1 %w", err1)
				err2 := New("error 2")
				wrapped2 := fmt.Errorf("wrapper 2 %w", err2)

				return New("outer").WithErrors(wrapped1, wrapped2)
			},
			searchText:       "wrapper",
			wantContainsText: false,
		},
		{
			name: "given_fmt_errorf_wrapped_custom_error_when_joined_then_wrapper_text_is_lost",
			buildErr: func() error {
				base := New("base error")
				wrapped := fmt.Errorf("text example %w", base)

				return Join(wrapped, New("another error"))
			},
			searchText:       "text example",
			wantContainsText: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				err := test.buildErr()
				errStr := err.Error()

				// then
				if test.wantContainsText {
					assert.Contains(t, errStr, test.searchText)
				} else {
					assert.NotContains(t, errStr, test.searchText)
				}
			},
		)
	}
}

// TestCompatibilityNilHandling tests that nil errors are handled consistently.
func TestCompatibilityNilHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		buildErr func() error
		name     string
		wantNil  bool
	}{
		{
			name: "given_nil_errors_when_joining_then_returns_nil",
			buildErr: func() error {
				return Join(nil, nil, nil)
			},
			wantNil: true,
		},
		{
			name: "given_mix_of_nil_and_errors_when_joining_then_returns_non_nil",
			buildErr: func() error {
				return Join(nil, New("error"), nil)
			},
			wantNil: false,
		},
		{
			name: "given_empty_message_when_creating_error_then_returns_non_nil",
			buildErr: func() error {
				return New("")
			},
			wantNil: false,
		},
		{
			name: "given_empty_message_when_creating_error_then_returns_non_nil",
			buildErr: func() error {
				return stderrors.New("")
			},
			wantNil: false,
		},
		{
			name: "given_nil_wrapped_error_when_using_with_errors_then_returns_non_nil",
			buildErr: func() error {
				return New("wrapper").WithErrors(nil)
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				// when
				err := test.buildErr()

				// then
				if test.wantNil {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			},
		)
	}
}
