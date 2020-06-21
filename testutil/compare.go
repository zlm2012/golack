package testutil

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Compare is similar to reflect.DeepEqual, but this prints out more useful information with error message as below:
// [team_profile_reorder > Profile > Fields > Element index at 0 > ID] Expected string is not set. Expected: Xf06054AAA. Actual: INVALID_ID.
// [file_comment_edited > Comment > Content] Expected string is not set. Expected: comment content. Actual: different comment text.
func Compare(hierarchy []string, expected reflect.Value, actual reflect.Value, t *testing.T) {
	if expected.Kind() != actual.Kind() {
		t.Errorf("%s Expected type is %s, but is %s.", hierarchyStr(hierarchy), expected.String(), actual.String())
		return
	}

	// Check zero value
	if !expected.IsValid() {
		if actual.IsValid() {
			fmt.Printf("%s is expected to be zero value, but is not: %#v", hierarchyStr(hierarchy), actual.String())
		}
		return
	}

	switch expected.Kind() {
	case reflect.String:
		if expected.String() != actual.String() {
			t.Errorf("%s Expected %s is not set. Expected: %s. Actual: %s.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.String(), actual.String())
		}

	case reflect.Bool:
		if expected.Bool() != actual.Bool() {
			t.Errorf("%s Expected %s is not set. Expected: %t. Actual: %t",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Bool(), actual.Bool())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if expected.Uint() != actual.Uint() {
			t.Errorf(
				"%s Expected %s is not set. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Uint(), actual.Uint())
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if expected.Int() != actual.Int() {
			t.Errorf(
				"%s Expected %s is not set. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Int(), actual.Int())
		}

	case reflect.Ptr, reflect.Interface:
		Compare(hierarchy, expected.Elem(), actual.Elem(), t)

	case reflect.Struct:
		for i := 0; i < expected.NumField(); i++ {
			tmp := make([]string, len(hierarchy))
			copy(tmp, hierarchy)
			tmp = append(tmp, expected.Type().Field(i).Name)
			Compare(tmp, expected.Field(i), actual.Field(i), t)
		}

	case reflect.Array, reflect.Slice:
		if expected.Len() != actual.Len() {
			t.Errorf("%s Element %s size differs. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Len(), actual.Len())
			return
		}

		tmp := make([]string, len(hierarchy))
		copy(tmp, hierarchy)
		for i := 0; i < expected.Len(); i++ {
			tmp = append(tmp, fmt.Sprintf("Element index at %d", i))
			Compare(tmp, expected.Index(i), actual.Index(i), t)
		}

	case reflect.Map:
		if expected.Len() != actual.Len() {
			t.Errorf("%s Element %s size differs. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Len(), actual.Len())
			return
		}

		if expected.Pointer() == actual.Pointer() {
			return
		}

		for _, k := range expected.MapKeys() {
			val1 := expected.MapIndex(k)
			val2 := expected.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() {
				t.Errorf("%s Expected map element is not given: %s.", hierarchyStr(hierarchy), k.String())
				return
			}

			tmp := make([]string, len(hierarchy))
			copy(tmp, hierarchy)
			tmp = append(tmp, k.String())
			Compare(tmp, val1, val2, t)
		}

	default:
		t.Errorf("Uncontrollable Kind %s is given: %#v", expected.Kind().String(), expected)

	}
}

func hierarchyStr(stack []string) string {
	return fmt.Sprintf("[%s]", strings.Join(stack, " > "))
}
