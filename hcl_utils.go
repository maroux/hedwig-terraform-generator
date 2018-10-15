package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"strconv"

	"fmt"

	"reflect"

	"sort"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/huandu/xstrings"
	"github.com/pkg/errors"
)

// Formats a Terraform file with hclfmt
func hclFmt(filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	res, err := printer.Format(src)
	if err != nil {
		return errors.Wrapf(err, "formatting file %q", filename)
	}

	if bytes.Equal(src, res) {
		return nil
	}

	return ioutil.WriteFile(filename, res, os.ModePerm)
}

// Formats all Terraform files with hclfmt
func hclFmtDir(module string) error {
	moduleInfo, err := ioutil.ReadDir(module)
	if err != nil {
		return err
	}
	for _, file := range moduleInfo {
		if filepath.Ext(file.Name()) != ".tf" {
			continue
		}
		if err := hclFmt(filepath.Join(module, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

// Converts an object of any type into hcl value. Handles arbitrarily nested maps and lists
func hclvalue(v interface{}) string {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	switch rv.Kind() {
	default:
		panic(fmt.Sprintf("Can not handle value '%v' of unknown kind: %v", v, rv.Kind()))
	case reflect.Bool:
		if rv.Bool() {
			return `"true"`
		} else {
			return `"false"`
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", rv.Uint())
	case reflect.Uintptr:
		return fmt.Sprintf("%d", rv.Elem().Uint())
	case reflect.Float32, reflect.Float64:
		fValue := rv.Float()
		iValue := int64(fValue)
		if fValue == float64(iValue) {
			return fmt.Sprintf("%d", iValue)
		} else {
			return fmt.Sprintf("%f", fValue)
		}
	case reflect.Array, reflect.Slice:
		return hcllist(rv)
	case reflect.Map:
		return hclobject(rv)
	case reflect.Ptr, reflect.Interface:
		return hclvalue(rv.Elem())
	case reflect.String:
		return strconv.Quote(rv.String())
	}
}

// Converts a go slice into an HCL list printing values with %v format specifier
// Handles arbitrarily nested maps and lists
func hcllist(rv reflect.Value) string {
	if rv.Len() == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteRune('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteRune('\n')
		b.WriteString(hclvalue(rv.Index(i)))
	}
	b.WriteRune('\n')
	b.WriteRune(']')
	return b.String()
}

// convert to string suitable to be used as a identifier name in HCL
// doesn't handle TF interpolated strings nicely
func hclident(s string) string {
	return xstrings.ToKebabCase(s)
}

type reflectMapKeyList []reflect.Value

func (l reflectMapKeyList) Len() int {
	return len(l)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (l reflectMapKeyList) Less(i, j int) bool {
	return fmt.Sprintf("%v", l[i]) < fmt.Sprintf("%v", l[j])
}

// Swap swaps the elements with indexes i and j.
func (l reflectMapKeyList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// converts a go map into HCL object printing keys and values with %v format specifier
// Handles arbitrarily nested maps and lists
func hclobject(rv reflect.Value) string {
	if rv.Len() == 0 {
		return "{}"
	}
	var b strings.Builder
	b.WriteRune('{')
	keys := reflectMapKeyList(rv.MapKeys())
	sort.Stable(keys)
	for _, key := range keys {
		b.WriteRune('\n')
		strKey := fmt.Sprintf("%v", key)
		quoted := strconv.Quote(strKey)
		// unwrap simple quotes
		if quoted != fmt.Sprintf(`"%s"`, strKey) {
			b.WriteString(quoted)
		} else {
			b.WriteString(strKey)
		}
		b.WriteString(" = ")
		b.WriteString(hclvalue(rv.MapIndex(key)))
	}
	b.WriteRune('\n')
	b.WriteRune('}')
	return b.String()
}
