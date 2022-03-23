package evli

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Source describe source of the event emitted, must be used for retrieve data
type Source struct {
	payload   Event
	emittedAd time.Time
	meta      Meta
	mu        sync.Mutex
	// schema is binding field of structured value
	schema map[string]reflect.Value
}

type Meta = map[string]interface{}

// EventNameMod defined which format of name you want
type EventNameMod int

const (
	// ShortName give struct name as "Struct" format
	ShortName EventNameMod = iota
	// StructName give struct name as "pkg.Struct" format
	LongName
	// PackageName give struct name as "folder/file/pkg.Struct" format
	PackageName
)

// Payload allow to retrieve data to from the event.
// must be a pointer if is struct, otherwise in func case, must be a declared value
func (s *Source) Payload(e interface{}) error {
	vl := reflect.ValueOf(e)
	if vl.Kind() == reflect.Ptr {
		vl = vl.Elem()
	}
	s.mu.Lock()
	for i := 0; i < vl.NumField(); i++ {
		name := vl.Type().Field(i).Name
		if tag, ok := vl.Type().Field(i).Tag.Lookup("field"); ok {
			name = tag
		}

		if val, ok := s.schema[name]; ok {
			associate(val, vl.Field(i))
		}
	}
	s.mu.Unlock()
	return nil
}

// EventName get name of the event
func (s *Source) EventName(mod EventNameMod) string {
	var name string
	s.mu.Lock()
	switch mod {
	case ShortName:
		name = reflect.TypeOf(s.payload).Name()
	case LongName:
		name = reflect.TypeOf(s.payload).String()
	case PackageName:
		name = fmt.Sprintf("%s.%s", reflect.TypeOf(s.payload).PkgPath(), reflect.TypeOf(s.payload).Name())
	}
	s.mu.Unlock()
	return name
}

// EmittedAd get the time when event were emitted
func (s *Source) EmittedAt() time.Time {
	s.mu.Lock()
	t := s.emittedAd
	s.mu.Unlock()
	return t
}

// Meta get metadata of source given by the emitter
func (s *Source) Meta() map[string]interface{} {
	s.mu.Lock()
	meta := s.meta
	s.mu.Unlock()
	return meta
}

// associate val value in target value
func associate(val, target reflect.Value) {
	if !target.CanSet() {
		return
	}
	if target.Kind() == val.Kind() {
		target.Set(val)
	}
}

// schema get schema of event source to easily associate in Subscriber
func schema(v interface{}) map[string]reflect.Value {
	vals := make(map[string]reflect.Value)
	vt := reflect.ValueOf(v)
	for i := 0; i < vt.Type().NumField(); i++ {
		f := vt.Field(i)
		name := vt.Type().Field(i).Name
		vals[name] = f
	}
	return vals
}
