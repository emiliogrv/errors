package errors

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Error returns the error message as a string.
// Implementation for rhe error built-in interface type for representing an error condition,
// with the nil value representing no error.
//
// The returned slog.Value will have the following attributes:
//   - Message
//   - Tags
//   - Attrs
//   - Errors
//   - Stack.
func (receiver *StructuredError) Error() string {
	var stringsBuilder strings.Builder

	receiver.asString(&stringsBuilder, zero)

	return stringsBuilder.String()
}

// String returns the error message as a string.
// It is equivalent to calling Error().
func (receiver *StructuredError) String() string {
	return receiver.Error()
}

// asString is the actual implementation for Error.
func (receiver *StructuredError) asString(stringsBuilder *strings.Builder, depth int) {
	if receiver == nil {
		valueToString(stringsBuilder, messageKey, nilValue)

		return
	}

	valueToString(stringsBuilder, messageKey, cmpOr(receiver.Message, nilValue))

	if len(receiver.Tags) > zero {
		stringsBuilder.WriteString(comma)
		stringsBuilder.WriteString(newLine)
		sliceToString(stringsBuilder, zero, tagsKey, receiver.Tags)
	}

	if len(receiver.Attrs) > zero {
		stringsBuilder.WriteString(comma)
		stringsBuilder.WriteString(newLine)
		sliceToString(stringsBuilder, depth, attrsKey, receiver.Attrs)
	}

	if len(receiver.Errors) > zero {
		target := normalizerTarget{
			errs: make([]error, zero, len(receiver.Errors)),
		}
		normalizeErrors(zero, &target, receiver.Errors...)

		stringsBuilder.WriteString(comma)
		stringsBuilder.WriteString(newLine)
		tabToString(stringsBuilder, depth)
		sliceToString(stringsBuilder, depth, errorsKey, target.errs)
	}

	if len(receiver.Stack) > zero {
		stringsBuilder.WriteString(comma)
		stringsBuilder.WriteString(newLine)
		valueToString(stringsBuilder, stackKey, string(receiver.Stack))
		stringsBuilder.WriteString(newLine)
	}
}

// String returns the error message as a string.
func (receiver *Attr) String() string {
	var stringsBuilder strings.Builder

	receiver.asString(&stringsBuilder, zero)

	return stringsBuilder.String()
}

// asString is the actual implementation for String.
//
//nolint:forcetypeassert,errcheck // XXXType helpers avoid using reflection
func (receiver *Attr) asString(stringsBuilder *strings.Builder, depth int) {
	if receiver == nil {
		valueToString(stringsBuilder, nilValue, nilValue)

		return
	}

	switch receiver.Type {
	case AnyType:
		valueToString(stringsBuilder, receiver.Key, fmt.Sprintf(verboseFormat, receiver.Value))
	case ObjectType:
		objectToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]Attr))
	case BoolType:
		valueToString(stringsBuilder, receiver.Key, strconv.FormatBool(receiver.Value.(bool)))
	case BoolsType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]bool))
	case TimeType:
		valueToString(stringsBuilder, receiver.Key, receiver.Value.(time.Time).String())
	case TimesType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]time.Time))
	case DurationType:
		valueToString(stringsBuilder, receiver.Key, receiver.Value.(time.Duration).String())
	case DurationsType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]time.Duration))
	case IntType:
		valueToString(stringsBuilder, receiver.Key, strconv.Itoa(receiver.Value.(int)))
	case IntsType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]int))
	case Int64Type:
		valueToString(stringsBuilder, receiver.Key, strconv.FormatInt(receiver.Value.(int64), ten))
	case Int64sType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]int64))
	case Uint64Type:
		valueToString(stringsBuilder, receiver.Key, strconv.FormatUint(receiver.Value.(uint64), ten))
	case Uint64sType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]uint64))
	case Float64Type:
		valueToString(stringsBuilder, receiver.Key, strconv.FormatFloat(receiver.Value.(float64), 'f', -1, sixtyFour))
	case Float64sType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]float64))
	case StringType:
		valueToString(stringsBuilder, receiver.Key, receiver.Value.(string))
	case StringsType:
		sliceToString(stringsBuilder, depth, receiver.Key, receiver.Value.([]string))
	default:
		valueToString(stringsBuilder, receiver.Key, fmt.Sprintf(verboseFormat, receiver.Value))
	}
}

// valueToString writes a key-value pair to the provided strings.Builder.
//
// Parameters:
//
//	stringsBuilder - the strings.Builder to write to
//	key - the key of the key-value pair
//	value - the value of the key-value pair
//
// Returns: A key-value pair is written to the provided strings.Builder.
func valueToString(stringsBuilder *strings.Builder, key, value string) {
	stringsBuilder.WriteString(parenthesisOpen)
	stringsBuilder.WriteString(key)
	stringsBuilder.WriteString(equals)
	stringsBuilder.WriteString(value)
	stringsBuilder.WriteString(parenthesisClose)
}

// errorToString writes an error to the provided strings.Builder.
//
// Parameters:
//
//	stringsBuilder - the strings.Builder to write to
//	depth - the depth to which the error is marshaled
//	err - the error to be written
//
// Returns: An error is written to the provided strings.Builder.
//
// The function writes a key-value pair to the provided strings.Builder.
// If err is nil, the function writes a key-value pair with the key "message" and the value "nil".
// If err is a StructuredError, the function writes a key-value pair with the same fields as the StructuredError.
// If err is not a StructuredError, the function writes a key-value pair with the key "message"
// and the value of the error's Error() method.
func errorToString(stringsBuilder *strings.Builder, depth int, err error) {
	var value *StructuredError
	switch {
	case err == nil:
		valueToString(stringsBuilder, messageKey, nilValue)
	case stderrors.As(err, &value):
		value.asString(stringsBuilder, depth)
	default:
		errStr := strings.TrimSpace(err.Error())
		valueToString(stringsBuilder, messageKey, cmpOr(errStr, nilValue))
	}
}

// objectToString writes an object to the provided strings.Builder.
//
// Parameters:
//
//	stringsBuilder - the strings.Builder to write to
//	depth - the depth to which the object is marshaled
//	key - the key of the key-value pair
//	object - the object to be written
//
// Returns: An object is written to the provided strings.Builder.
//
// The function writes a key-value pair to the provided strings.Builder.
// If object is nil, the function writes a key-value pair with the key "message" and the value "nil".
// If object is a slice of Attr, the function writes a key-value pair with the same fields as the slice of Attr.
func objectToString(stringsBuilder *strings.Builder, depth int, key string, object []Attr) {
	valuesToString(stringsBuilder, depth, key, object, curlyOpen, curlyClose)
}

// sliceToString writes a slice to the provided strings.Builder.
//
// Parameters:
//
//	stringsBuilder - the strings.Builder to write to
//	depth - the depth to which the slice is marshaled
//	key - the key of the key-value pair
//	slice - the slice to be written
//
// Returns: A slice is written to the provided strings.Builder.
//
// The function writes a key-value pair to the provided strings.Builder.
// If slice is empty, the function writes nothing.
// If slice is not empty, the function writes a key-value pair with the same fields as the slice.
func sliceToString[T any](stringsBuilder *strings.Builder, depth int, key string, slice []T) {
	valuesToString(stringsBuilder, depth, key, slice, bracketOpen, bracketClose)
}

// valuesToString writes a slice to the provided strings.Builder.
//
// Parameters:
//
//	stringsBuilder - the strings.Builder to write to
//	depth - the depth to which the slice is marshaled
//	key - the key of the key-value pair
//	slice - the slice to be written
//	opener - the opening string to write
//	closer - the closing string to write
//
// Returns: A slice is written to the provided strings.Builder.
//
// The function writes a key-value pair to the provided strings.Builder.
// If slice is empty, the function writes nothing.
// If slice is not empty, the function writes a key-value pair with the same fields as the slice.
func valuesToString[T any](stringsBuilder *strings.Builder, depth int, key string, slice []T, opener, closer string) {
	stringsBuilder.WriteString(parenthesisOpen)
	stringsBuilder.WriteString(key)
	stringsBuilder.WriteString(equals)
	stringsBuilder.WriteString(opener)

	if len(slice) == zero {
		stringsBuilder.WriteString(closer)

		return
	}

	stringsBuilder.WriteString(newLine)

	depth++

	switch values := any(slice).(type) {
	case []Attr:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			value.asString(stringsBuilder, depth)
		}
	case []error:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			errorToString(stringsBuilder, depth, value)
		}
	case []bool:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(strconv.FormatBool(value))
		}
	case []time.Time:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(value.String())
		}
	case []time.Duration:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(value.String())
		}
	case []int:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(strconv.Itoa(value))
		}
	case []int64:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(strconv.FormatInt(value, ten))
		}
	case []uint64:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(strconv.FormatUint(value, ten))
		}
	case []float64:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(strconv.FormatFloat(value, 'f', -1, sixtyFour))
		}
	case []string:
		for index, value := range values {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			stringsBuilder.WriteString(strings.TrimSpace(value))
		}
	default:
		for index, value := range slice {
			if index > zero {
				stringsBuilder.WriteString(comma)
				stringsBuilder.WriteString(newLine)
			}

			tabToString(stringsBuilder, depth)
			_, _ = fmt.Fprintf(stringsBuilder, verboseFormat, value)
		}
	}

	stringsBuilder.WriteString(newLine)
	tabToString(stringsBuilder, depth-1)
	stringsBuilder.WriteString(closer)
	stringsBuilder.WriteString(parenthesisClose)
}

// tabToString writes depth number of tabs to the provided strings.Builder.
//
// Parameters:
//
//	stringsBuilder - the strings.Builder to write to
//	depth - the number of tabs to write
//
// Returns: depth number of tabs are written to the provided strings.Builder.
func tabToString(stringsBuilder *strings.Builder, depth int) {
	for i := zero; i < depth; i++ {
		stringsBuilder.WriteString(tab)
	}
}
