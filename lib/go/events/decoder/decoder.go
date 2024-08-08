package decoder

import (
	"fmt"
	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/ccf"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	pkgerrors "github.com/pkg/errors"
)

func GetCadenceEvent(payload []byte) (cadence.Event, error) {
	cadenceValue, err := DecodeCadenceValue(payload)
	if err != nil {
		return cadence.Event{}, err
	}
	return cadenceValue.(cadence.Event), nil
}

func DecodeCadenceValue(payload []byte) (cadence.Value, error) {
	cadenceValue, err := jsoncdc.Decode(nil,
		payload,
		jsoncdc.WithBackwardsCompatibility(),
		jsoncdc.WithAllowUnstructuredStaticTypes(true),
	)
	if err != nil {
		// json decode failed, try ccf decode
		fmt.Printf("json decode failed, try ccf decode, reason: %s\n", err.Error())
		ccfValue, ccfErr := ccf.Decode(nil, payload)
		if ccfErr != nil {
			return cadence.Event{}, pkgerrors.Wrapf(
				ccfErr,
				"failed to decode cadence value through ccf, json decode attempt err: %s",
				err.Error(),
			)
		}
		cadenceValue = ccfValue
	}
	return cadenceValue, nil
}

func ConvertStringMetadata(metadata cadence.Dictionary) map[string]string {
	md := map[string]string{}
	for _, kvp := range metadata.Pairs {
		md[string(kvp.Key.(cadence.String))] = string(kvp.Value.(cadence.String))
	}
	return md
}

func ConvertStringArrayMetadata(metadata cadence.Dictionary) map[string][]string {
	md := map[string][]string{}
	for _, kvp := range metadata.Pairs {
		md[string(kvp.Key.(cadence.String))], _ = convertStringArray(kvp.Value.(cadence.Array))
	}
	return md
}

func convertStringArray(md cadence.Value) ([]string, error) {
	arr := md.(cadence.Array)
	var items []string
	for _, v := range arr.Values {
		converted, err := GetFieldValue(v)
		if err != nil {
			return nil, err
		}
		items = append(items, converted.(string))
	}
	return items, nil
}

func DecodeToEventMap(payload []byte) (map[string]any, error) {
	cadenceValue, err := GetCadenceEvent(payload)
	if err != nil {
		return nil, err
	}
	return ConvertEvent(cadenceValue)
}

func ConvertEvent(evt cadence.Event) (map[string]any, error) {
	output := map[string]any{}
	for k, v := range cadence.FieldsMappedByName(evt) {
		converted, err := GetFieldValue(v)
		if err != nil {
			return nil, err
		}
		output[k] = converted
	}
	return output, nil
}

func ConvertObjectMetadata(value cadence.Composite) (map[string]any, error) {
	structMap := map[string]any{}
	subFields := cadence.FieldsMappedByName(value)
	for key, subField := range subFields {
		val, err := GetFieldValue(subField)
		if err != nil {
			return nil, err
		}
		if val != nil {
			structMap[key] = val
		}
	}
	return structMap, nil
}

// GetFieldValue Convert a cadence value into a any structure for easier consumption in go with options
func GetFieldValue(md cadence.Value) (any, error) {
	switch field := md.(type) {
	case cadence.Optional:
		if field.Value == nil {
			return nil, nil
		}
		return GetFieldValue(field.Value)
	case cadence.Dictionary:
		return convertDict(field)
	case cadence.Array:
		return convertArray(field)
	case cadence.Int:
		return field.Int(), nil
	case cadence.Int8:
		return int8(field), nil
	case cadence.Int16:
		return int16(field), nil
	case cadence.Int32:
		return int32(field), nil
	case cadence.Int64:
		return int64(field), nil
	case cadence.UInt8:
		return uint8(field), nil
	case cadence.UInt16:
		return uint16(field), nil
	case cadence.UInt32:
		return uint32(field), nil
	case cadence.UInt64:
		return uint64(field), nil
	case cadence.Word8:
		return uint8(field), nil
	case cadence.Word16:
		return uint16(field), nil
	case cadence.Word32:
		return uint32(field), nil
	case cadence.Word64:
		return uint64(field), nil
	case cadence.TypeValue:
		return field.StaticType.ID(), nil
	case cadence.String:
		return string(field), nil
	case cadence.UFix64:
		return uint64(field), nil
	case cadence.Fix64:
		return int64(field), nil
	case cadence.Struct:
		return ConvertObjectMetadata(field)
	case cadence.Resource:
		return ConvertObjectMetadata(field)
	case cadence.Bool:
		return bool(field), nil
	case cadence.Bytes:
		return []byte(field), nil
	case cadence.Character:
		return string(field), nil
	case cadence.Function:
		return field.FunctionType.ID(), nil
	case cadence.Address:
		return field.String()[2:], nil
	default:
		return field.String(), nil
	}
}

func convertArray(md cadence.Value) (any, error) {
	arr := md.(cadence.Array)
	var items []any
	for _, v := range arr.Values {
		converted, err := GetFieldValue(v)
		if err != nil {
			return nil, err
		}
		items = append(items, converted)
	}
	return items, nil
}

func convertDict(md cadence.Value) (map[any]any, error) {
	d, ok := md.(cadence.Dictionary)
	if !ok {
		return nil, fmt.Errorf("value is not a dictionary, got %T", md)
	}
	valMap := map[any]any{}
	for _, item := range d.Pairs {
		value, err := GetFieldValue(item.Value)
		if err != nil {
			return nil, err
		}
		key, err := GetFieldValue(item.Key)
		if err != nil {
			return nil, err
		}
		if key == "" {
			return nil, fmt.Errorf("keys cannot be empty")
		}
		if value != nil {
			valMap[key] = value
		}
	}
	return valMap, nil
}
