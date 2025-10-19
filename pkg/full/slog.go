package errors

import (
	stderrors "errors"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// LogValue returns a slog.Value representation of the receiver.
//
// The returned slog.Value will have the following attributes:
//   - Message
//   - Tags
//   - Attrs
//   - Errors
//   - Stack.
//
// If the receiver is not nil, the returned slog.Value is guaranteed not to be of Kind slog.KindLogValuer.
// If the receiver is nil, the returned slog.Value is guaranteed to be of Kind slog.KindGroup.
//
// Usage must be with slog.Any or slog.Group.
func (receiver *StructuredError) LogValue() slog.Value {
	if receiver == nil {
		return slog.GroupValue(slog.String(messageKey, nilValue))
	}

	length := one

	if len(receiver.Attrs) > zero {
		length++
	}

	if len(receiver.Errors) > zero {
		length++
	}

	if len(receiver.Tags) > zero {
		length++
	}

	if len(receiver.Stack) > zero {
		length++
	}

	values := make([]slog.Attr, zero, length)
	values = append(values, slog.String(messageKey, cmpOr(receiver.Message, nilValue)))

	if len(receiver.Tags) > zero {
		values = append(values, sliceToSlog(tagsKey, receiver.Tags))
	}

	if len(receiver.Attrs) > zero {
		values = append(values, sliceToSlog(attrsKey, receiver.Attrs))
	}

	if len(receiver.Errors) > zero {
		target := normalizerTarget{
			errs: make([]error, zero, len(receiver.Errors)),
		}
		normalizeErrors(zero, &target, receiver.Errors...)

		values = append(values, sliceToSlog(errorsKey, target.errs))
	}

	if len(receiver.Stack) > zero {
		values = append(values, sliceToSlog(stackKey, strings.Split(string(receiver.Stack), newLine)))
	}

	return slog.GroupValue(values...)
}

// LogValue returns a slog.Value representation of the receiver.
//
// The returned slog.Value will have a single attribute with the key
// being the receiver's key and the value being the receiver's value.
//
// If the receiver is nil, the returned slog.Value is guaranteed to be of
// Kind slog.KindGroup with a single attribute of key nilValue and value nilValue.
// If the receiver is not nil, the returned slog.Value is guaranteed not to be of
// Kind slog.KindLogValuer.
//
// Usage must be with slog.Any or slog.Group.
func (receiver *Attr) LogValue() slog.Value {
	return slog.GroupValue(receiver.asSlog())
}

// asSlog is the actual implementation for LogValue.
//
//nolint:forcetypeassert,errcheck // XXXType helpers avoid using reflection
func (receiver *Attr) asSlog() slog.Attr {
	if receiver == nil {
		return slog.String(nilValue, nilValue)
	}

	switch receiver.Type {
	case AnyType:
		return slog.Any(receiver.Key, receiver.Value)
	case ObjectType:
		return sliceToSlog(receiver.Key, receiver.Value.([]Attr))
	case BoolType:
		return slog.Bool(receiver.Key, receiver.Value.(bool))
	case BoolsType:
		return sliceToSlog(receiver.Key, receiver.Value.([]bool))
	case TimeType:
		return slog.Time(receiver.Key, receiver.Value.(time.Time))
	case TimesType:
		return sliceToSlog(receiver.Key, receiver.Value.([]time.Time))
	case DurationType:
		return slog.Duration(receiver.Key, receiver.Value.(time.Duration))
	case DurationsType:
		return sliceToSlog(receiver.Key, receiver.Value.([]time.Duration))
	case IntType:
		return slog.Int(receiver.Key, receiver.Value.(int))
	case IntsType:
		return sliceToSlog(receiver.Key, receiver.Value.([]int))
	case Int64Type:
		return slog.Int64(receiver.Key, receiver.Value.(int64))
	case Int64sType:
		return sliceToSlog(receiver.Key, receiver.Value.([]int64))
	case Uint64Type:
		return slog.Uint64(receiver.Key, receiver.Value.(uint64))
	case Uint64sType:
		return sliceToSlog(receiver.Key, receiver.Value.([]uint64))
	case Float64Type:
		return slog.Float64(receiver.Key, receiver.Value.(float64))
	case Float64sType:
		return sliceToSlog(receiver.Key, receiver.Value.([]float64))
	case StringType:
		return slog.String(receiver.Key, receiver.Value.(string))
	case StringsType:
		return sliceToSlog(receiver.Key, receiver.Value.([]string))
	default:
		return slog.Any(receiver.Key, receiver.Value)
	}
}

// errorToSlog converts an error to a slog.Attr.
//
// If the error is nil, the returned slog.Attr will have a single attribute with the key
// being the given key and the value being nilValue.
//
// If the error is a *StructuredError, the returned slog.Attr will have a single attribute
// with the key being the given key and the value being the *StructuredError's slog.Value
// representation.
//
// If the error is not a *StructuredError, the returned slog.Attr will have a single attribute
// with the key being the given key and the value being the error's Error() method.
//
// Parameters:
//
//	key - the key of the returned slog.Attr
//	err - the error to be converted to a slog.Attr
func errorToSlog(key string, err error) slog.Attr {
	var value *StructuredError
	switch {
	case err == nil:
		return slog.Group(key, slog.String(messageKey, nilValue))
	case stderrors.As(err, &value):
		return slog.Attr{Key: key, Value: value.LogValue()}
	default:
		errStr := strings.TrimSpace(err.Error())

		return slog.Group(key, slog.String(messageKey, cmpOr(errStr, nilValue)))
	}
}

// sliceToSlog converts a slice of any type to a slice of slog.Attr.
// It is needed in order to avoid reflection as much as possible.
func sliceToSlog[T any](key string, slice []T) slog.Attr {
	if len(slice) == zero {
		return slog.Group(key)
	}

	attrs := make([]slog.Attr, zero, len(slice))

	switch values := any(slice).(type) {
	case []Attr:
		for _, attr := range values {
			attrs = append(attrs, attr.asSlog())
		}
	case []error:
		for i, value := range values {
			attrs = append(attrs, errorToSlog(strconv.Itoa(i), value))
		}
	case []bool:
		for i, value := range values {
			attrs = append(attrs, slog.Bool(strconv.Itoa(i), value))
		}
	case []time.Time:
		for i, value := range values {
			attrs = append(attrs, slog.Time(strconv.Itoa(i), value))
		}
	case []time.Duration:
		for i, value := range values {
			attrs = append(attrs, slog.Duration(strconv.Itoa(i), value))
		}
	case []int:
		for i, value := range values {
			attrs = append(attrs, slog.Int(strconv.Itoa(i), value))
		}
	case []int64:
		for i, value := range values {
			attrs = append(attrs, slog.Int64(strconv.Itoa(i), value))
		}
	case []uint64:
		for i, value := range values {
			attrs = append(attrs, slog.Uint64(strconv.Itoa(i), value))
		}
	case []float64:
		for i, value := range values {
			attrs = append(attrs, slog.Float64(strconv.Itoa(i), value))
		}
	case []string:
		for i, value := range values {
			attrs = append(attrs, slog.String(strconv.Itoa(i), strings.TrimSpace(value)))
		}
	default:
		for i, value := range slice {
			attrs = append(attrs, slog.Any(strconv.Itoa(i), value))
		}
	}

	return slog.Attr{Key: key, Value: slog.GroupValue(attrs...)}
}
