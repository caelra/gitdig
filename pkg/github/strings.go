// package github

// import (
// 	"bytes"
// 	"fmt"
// 	"reflect"
// 	"time"
// )

// type Timestamp struct {
// 	time.Time
// }

// var timestampType = reflect.TypeOf(Timestamp{})

// // source: https://github.com/google/go-github/blob/master/github/strings.go

// // Stringify attempts to create a reasonable string representation of types in
// // the GitHub library. It does things like resolve pointers to their values
// // and omits struct fields with nil values.
// func Stringify(message interface{}) string {
// 	var buf bytes.Buffer
// 	v := reflect.ValueOf(message)
// 	stringifyValue(&buf, v)
// 	return buf.String()
// }

// // stringifyValue was heavily inspired by the goprotobuf library.
// // 2nd heavily inspired on the referenced source with some improvements

// func stringifyValue(w *bytes.Buffer, val reflect.Value) {
// 	if val.Kind() == reflect.Ptr && val.IsNil() {
// 		w.WriteString("<nil>")
// 		return
// 	}

// 	v := reflect.Indirect(val)

// 	switch v.Kind() {
// 	case reflect.String:
// 		fmt.Fprintf(w, `"%s"`, v.String())
// 	case reflect.Slice:
// 		w.WriteByte('[')
// 		for i := 0; i < v.Len(); i++ {
// 			if i > 0 {
// 				w.WriteByte(' ')
// 			}
// 			stringifyValue(w, v.Index(i))
// 		}

// 		w.WriteByte(']')
// 	case reflect.Struct:
// 		if v.Type().Name() != "" {
// 			w.WriteString(v.Type().String())
// 		}

// 		if v.Type() == timestampType {
// 			fmt.Fprintf(w, "{%s}", v.Interface())
// 			return
// 		}

// 		w.WriteByte('{')

// 		for i := 0; i < v.NumField(); i++ {
// 			fv := v.Field(i)
// 			if isEmptyValue(fv) {
// 				continue
// 			}
// 			if i > 0 {
// 				w.WriteString(", ")
// 			}

// 			w.WriteString(v.Type().Field(i).Name)
// 			w.WriteByte(':')
// 			stringifyValue(w, fv)
// 		}

// 		w.WriteByte('}')
// 	default:
// 		if v.CanInterface() {
// 			fmt.Fprint(w, v.Interface())
// 		}
// 	}
// }

package github

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

type Timestamp struct {
	time.Time
}

var timestampType = reflect.TypeOf(Timestamp{})

// Stringify attempts to create a reasonable string representation of types in
// the GitHub library. It does things like resolve pointers to their values
// and omits struct fields with nil values.
func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

// stringifyValue was heavily inspired by the goprotobuf library.

func stringifyValue(w *bytes.Buffer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i))
		}

		w.Write([]byte{']'})
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}

		// special handling of Timestamp values
		if v.Type() == timestampType {
			fmt.Fprintf(w, "{%s}", v.Interface())
			return
		}

		w.Write([]byte{'{'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Map && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(", "))
			} else {
				sep = true
			}

			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv)
		}

		w.Write([]byte{'}'})
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
