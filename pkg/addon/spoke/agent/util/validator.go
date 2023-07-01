package util

import (
	"math"

	"github.com/spf13/cast"

	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/models"
)

func checkValueInRange(valueType string, reading interface{}) bool {
	isValid := false

	if valueType == models.ValueTypeString || valueType == models.ValueTypeBool || valueType == models.ValueTypeObject {
		return true
	}

	if valueType == models.ValueTypeInt8 || valueType == models.ValueTypeInt16 ||
		valueType == models.ValueTypeInt32 || valueType == models.ValueTypeInt64 {
		val := cast.ToInt64(reading)
		isValid = checkIntValueRange(valueType, val)
	}

	if valueType == models.ValueTypeUint8 || valueType == models.ValueTypeUint16 ||
		valueType == models.ValueTypeUint32 || valueType == models.ValueTypeUint64 {
		val := cast.ToUint64(reading)
		isValid = checkUintValueRange(valueType, val)
	}

	if valueType == models.ValueTypeFloat32 || valueType == models.ValueTypeFloat64 {
		val := cast.ToFloat64(reading)
		isValid = checkFloatValueRange(valueType, val)
	}

	return isValid
}

func checkUintValueRange(valueType string, val uint64) bool {
	var isValid = false
	switch valueType {
	case models.ValueTypeUint8:
		if val <= math.MaxUint8 {
			isValid = true
		}
	case models.ValueTypeUint16:
		if val <= math.MaxUint16 {
			isValid = true
		}
	case models.ValueTypeUint32:
		if val <= math.MaxUint32 {
			isValid = true
		}
	case models.ValueTypeUint64:
		maxiMum := uint64(math.MaxUint64)
		if val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(valueType string, val int64) bool {
	var isValid = false
	switch valueType {
	case models.ValueTypeInt8:
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case models.ValueTypeInt16:
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case models.ValueTypeInt32:
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case models.ValueTypeInt64:
		isValid = true
	}
	return isValid
}

func checkFloatValueRange(valueType string, val float64) bool {
	var isValid = false
	switch valueType {
	case models.ValueTypeFloat32:
		if !math.IsNaN(val) && math.Abs(val) <= math.MaxFloat32 {
			isValid = true
		}
	case models.ValueTypeFloat64:
		if !math.IsNaN(val) && !math.IsInf(val, 0) {
			isValid = true
		}
	}
	return isValid
}
