package goconf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// formatter
	CRLF       = '\n'
	Commet     = "#"
	Spliter    = " "
	Include    = "include"
	SectionS   = "["
	SectionE   = "]"
	includeLen = len(Include)
	// memory unit
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

// Section is the key-value data object.
type Section struct {
	data map[string]string
	Name string
}

// Config is the key-value configuration object.
type Config struct {
	data    map[string]*Section
	file    string
	Commet  string
	Spliter string
}

// New return a new default Config object (commet = '#', spliter = ' ').
func New() *Config {
	return &Config{Commet: Commet, Spliter: Spliter}
}

// Parse parse the specified config file.
func (c *Config) Parse(file string) error {
	c.data = map[string]*Section{}
	c.file = file
	// open config file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	files := []string{file}
	fileMap := map[string]bool{file: true}
	section := ""
	filename := file
	lineNum := 0
	for {
		lineNum++
		line, err := rd.ReadString(CRLF)
		if err == io.EOF && len(line) == 0 {
			// parse file finish
			// all files parsed, break
			if len(files) <= 1 {
				break
			}
			// get the next file
			files = files[1:]
			// reset
			section = ""
			lineNum = 0
			filename = files[0]
			f, err = os.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close()
			rd = bufio.NewReader(f)
			continue
		} else if err != nil && err != io.EOF {
			return err
		}
		// trim space
		line = strings.TrimSpace(line)
		// ignore blank line
		if line == "" {
			continue
		}
		// ignore commet line
		if strings.HasPrefix(line, c.Commet) {
			continue
		}
		// get secion
		if strings.HasPrefix(line, SectionS) {
			if !strings.HasSuffix(line, SectionE) {
				return errors.New(fmt.Sprintf("no end section: %s in %s:%d", SectionE, filename, lineNum))
			}
			section = line[1 : len(line)-1]
			continue
		}
		// handle include
		if strings.HasPrefix(line, Include) {
			if len(line) > includeLen {
				// add other config files
				newFiles, err := includeFiles(strings.TrimSpace(line[includeLen+1:]), fileMap)
				if err != nil {
					return err
				}

				files = append(files, newFiles...)
				continue
			} else {
				return errors.New(fmt.Sprintf("no include file in %s:%d", filename, lineNum))
			}
		}
		// get the spliter index
		// xinhuang327: Add support for empty value
		key := line
		value := ""
		idx := strings.Index(line, c.Spliter)
		if idx > 0 {
			// get the key and value
			key = strings.TrimSpace(line[:idx])
			if len(line) > idx {
				value = strings.TrimSpace(line[idx+1:])
			}
		}
		// store the key-value config
		s, ok := c.data[section]
		if !ok {
			s = &Section{data: map[string]string{}, Name: section}
			c.data[section] = s
		}
		s.data[key] = value
	}
	return nil
}

// parse config file include other files
func includeFiles(path string, fileMap map[string]bool) ([]string, error) {
	files := []string{}
	// match pattern
	pattern := filepath.Base(path)
	dirName := filepath.Dir(path)
	// get child files
	dir, err := os.Open(dirName)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fis, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		// skip dir
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		if ok, err := filepath.Match(pattern, name); err != nil {
			return nil, err
		} else if ok {
			file := filepath.Join(dirName, name)
			if _, exist := fileMap[file]; exist {
				return nil, errors.New(fmt.Sprintf("include duplicate file: %s", file))
			}
			files = append(files, file)
			// save parse file
			fileMap[file] = true
		}
	}
	return files, nil
}

// Get get a config section by key.
func (c *Config) Get(section string) *Section {
	s, ok := c.data[section]
	if ok {
		return s
	} else {
		return nil
	}
}

// Add add a new config section, if exist the section key then return the
// existing one.
func (c *Config) Add(section string) *Section {
	s, ok := c.data[section]
	if !ok {
		s = &Section{data: map[string]string{}, Name: section}
		c.data[section] = s
	}
	return s
}

// Remove remove the specified section.
func (c *Config) Remove(section string) {
	delete(c.data, section)
}

// Sections return all the config sections.
func (c *Config) Sections() []string {
	sections := []string{}
	for k, _ := range c.data {
		sections = append(sections, k)
	}
	return sections
}

// Save save current configuration to specified file, if file is "" then rewrite the original file.
//
// This method will ignore all the comment and include instruction if original file has.
func (c *Config) Save(file string) error {
	if file == "" {
		file = c.file
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for section, data := range c.data {
		if _, err := f.WriteString(fmt.Sprintf("[%s]%c", section, CRLF)); err != nil {
			return err
		}
		for k, v := range data.data {
			if _, err := f.WriteString(fmt.Sprintf("%s%s%s%c", k, c.Spliter, v, CRLF)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Reload reload the config file and return a new Config.
func (c *Config) Reload() (*Config, error) {
	nc := &Config{Commet: c.Commet, Spliter: c.Spliter}
	if err := nc.Parse(c.file); err != nil {
		return nil, err
	}
	return nc, nil
}

// Add add a new key-value configuration for the section.
func (s *Section) Add(k, v string) {
	s.data[k] = v
}

// Remove remove the specified key configuration for the section.
func (s *Section) Remove(k string) {
	delete(s.data, k)
}

// An NoKeyError describes a goconf key that was not found in the section.
type NoKeyError struct {
	Key     string
	Section string
}

func (e *NoKeyError) Error() string {
	return fmt.Sprintf("key: \"%s\" not found in [%s]", e.Key, e.Section)
}

// String get config string value.
func (s *Section) String(key string) (string, error) {
	if v, ok := s.data[key]; ok {
		return v, nil
	} else {
		return "", &NoKeyError{Key: key, Section: s.Name}
	}
}

// Strings get config []string value.
func (s *Section) Strings(key, delim string) ([]string, error) {
	if v, ok := s.data[key]; ok {
		return strings.Split(v, delim), nil
	} else {
		return nil, &NoKeyError{Key: key, Section: s.Name}
	}
}

// Int get config int value.
func (s *Section) Int(key string) (int64, error) {
	if v, ok := s.data[key]; ok {
		return strconv.ParseInt(v, 10, 64)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.Name}
	}
}

// Uint get config uint value.
func (s *Section) Uint(key string) (uint64, error) {
	if v, ok := s.data[key]; ok {
		return strconv.ParseUint(v, 10, 64)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.Name}
	}
}

// Float get config float value.
func (s *Section) Float(key string) (float64, error) {
	if v, ok := s.data[key]; ok {
		return strconv.ParseFloat(v, 64)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.Name}
	}
}

// Bool get config boolean value.
//
// "yes", "1", "y", "true", "enable" means true.
//
// "no", "0", "n", "false", "disable" means false.
//
// if the specified value unknown then return false.
func (s *Section) Bool(key string) (bool, error) {
	if v, ok := s.data[key]; ok {
		v = strings.ToLower(v)
		return parseBool(v), nil
	} else {
		return false, &NoKeyError{Key: key, Section: s.Name}
	}
}

func parseBool(v string) bool {
	if v == "true" || v == "yes" || v == "1" || v == "y" || v == "enable" {
		return true
	} else if v == "false" || v == "no" || v == "0" || v == "n" || v == "disable" {
		return false
	} else {
		return false
	}
}

// Byte get config byte number value.
//
// 1kb = 1k = 1024.
//
// 1mb = 1m = 1024 * 1024.
//
// 1gb = 1g = 1024 * 1024 * 1024.
func (s *Section) MemSize(key string) (int, error) {
	if v, ok := s.data[key]; ok {
		return parseMemory(v)
	} else {
		return 0, &NoKeyError{Key: key, Section: s.Name}
	}
}

func parseMemory(v string) (int, error) {
	unit := Byte
	subIdx := len(v)
	if strings.HasSuffix(v, "k") {
		unit = KB
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "kb") {
		unit = KB
		subIdx = subIdx - 2
	} else if strings.HasSuffix(v, "m") {
		unit = MB
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "mb") {
		unit = MB
		subIdx = subIdx - 2
	} else if strings.HasSuffix(v, "g") {
		unit = GB
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "gb") {
		unit = GB
		subIdx = subIdx - 2
	}
	b, err := strconv.ParseInt(v[:subIdx], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(b) * unit, nil
}

// Duration get config second value.
//
// 1s = 1sec = 1.
//
// 1m = 1min = 60.
//
// 1h = 1hour = 60 * 60.
func (s *Section) Duration(key string) (time.Duration, error) {
	if v, ok := s.data[key]; ok {
		if t, err := parseTime(v); err != nil {
			return 0, err
		} else {
			return time.Duration(t), nil
		}
	} else {
		return 0, &NoKeyError{Key: key, Section: s.Name}
	}
}

func parseTime(v string) (int64, error) {
	unit := int64(time.Nanosecond)
	subIdx := len(v)
	if strings.HasSuffix(v, "ms") {
		unit = int64(time.Millisecond)
		subIdx = subIdx - 2
	} else if strings.HasSuffix(v, "s") {
		unit = int64(time.Second)
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "sec") {
		unit = int64(time.Second)
		subIdx = subIdx - 3
	} else if strings.HasSuffix(v, "m") {
		unit = int64(time.Minute)
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "min") {
		unit = int64(time.Minute)
		subIdx = subIdx - 3
	} else if strings.HasSuffix(v, "h") {
		unit = int64(time.Hour)
		subIdx = subIdx - 1
	} else if strings.HasSuffix(v, "hour") {
		unit = int64(time.Hour)
		subIdx = subIdx - 4
	}
	b, err := strconv.ParseInt(v[:subIdx], 10, 64)
	if err != nil {
		return 0, err
	}
	return b * unit, nil
}

// Keys return all the section keys.
func (s *Section) Keys() []string {
	keys := []string{}
	for k, _ := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "goconf: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "goconf: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "goconf: Unmarshal(nil " + e.Type.String() + ")"
}

// Unmarshal parses the goconf struct and stores the result in the value
// pointed to by v.
//
// Struct values encode as goconf objects. Each exported struct field
// becomes a member of the object unless
//   - the field's tag is "-", or
//   - the field is empty and its tag specifies the "omitempty" option.
// The empty values are false, 0, any
// nil pointer or interface value, and any array, slice, map, or string of
// length zero. The object's section and key string is the struct field name
// but can be specified in the struct field's tag value. The "goconf" key in
// the struct field's tag value is the key name, followed by an optional comma
// and options. Examples:
//
//   // Field is ignored by this package.
//   Field int `goconf:"-"`
//
//   // Field appears in goconf section "base" as key "myName".
//   Field int `goconf:"base:myName"`
//
//   // Field appears in goconf section "base" as key "myName", the value split
//   // by delimiter ",".
//   Field []string `goconf:"base:myName:,"`
//
//   // Field appears in goconf section "base" as key "myName", the value split
//   // by delimiter "," and key-value is splited by "=".
//   Field map[int]string `goconf:"base:myName:,"`
//
//   // Field appears in goconf section "base" as key "myName", the value
//   // conver to time.Duration. When has extra tag "time", then goconf can
//   // parse such "1h", "1s" config values.
//   //
//   // Note the extra tag "time" only effect the int64 (time.Duration is int64)
//   Field time.Duration `goconf:"base:myName:time"`
//
//   // Field appears in goconf section "base" as key "myName", when has extra
//   // tag, then goconf can parse like "1gb", "1mb" config values.
//   //
//   // Note the extra tag "time" only effect the int (memory size is int).
//   Field int `goconf:"base:myName:memory"`
//
func (c *Config) Unmarshal(v interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr || vv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	rv := vv.Elem()
	rt := rv.Type()
	n := rv.NumField()
	// enum every struct field
	for i := 0; i < n; i++ {
		vf := rv.Field(i)
		tf := rt.Field(i)
		tag := tf.Tag.Get("goconf")
		// if tag empty or "-" ignore
		if tag == "-" || tag == "" || tag == "omitempty" {
			continue
		}
		tagArr := strings.SplitN(tag, ":", 3)
		if len(tagArr) < 2 {
			return errors.New(fmt.Sprintf("error tag: %s, must be section:field:delim(optional)", tag))
		}
		section := tagArr[0]
		key := tagArr[1]
		s := c.Get(section)
		if s == nil {
			// no config section
			continue
		}
		value, ok := s.data[key]
		if !ok {
			// no confit key
			continue
		}
		switch vf.Kind() {
		case reflect.String:
			vf.SetString(value)
		case reflect.Bool:
			vf.SetBool(parseBool(value))
		case reflect.Float32:
			if tmp, err := strconv.ParseFloat(value, 32); err != nil {
				return err
			} else {
				vf.SetFloat(tmp)
			}
		case reflect.Float64:
			if tmp, err := strconv.ParseFloat(value, 64); err != nil {
				return err
			} else {
				vf.SetFloat(tmp)
			}
		case reflect.Int:
			if len(tagArr) == 3 {
				format := tagArr[2]
				// parse memory size
				if format == "memory" {
					if tmp, err := parseMemory(value); err != nil {
						return err
					} else {
						vf.SetInt(int64(tmp))
					}
				} else {
					return errors.New(fmt.Sprintf("unknown tag: %s in struct field: %s (support tags: \"memory\")", format, tf.Name))
				}
			} else {
				if tmp, err := strconv.ParseInt(value, 10, 32); err != nil {
					return err
				} else {
					vf.SetInt(tmp)
				}
			}
		case reflect.Int8:
			if tmp, err := strconv.ParseInt(value, 10, 8); err != nil {
				return err
			} else {
				vf.SetInt(tmp)
			}
		case reflect.Int16:
			if tmp, err := strconv.ParseInt(value, 10, 16); err != nil {
				return err
			} else {
				vf.SetInt(tmp)
			}
		case reflect.Int32:
			if tmp, err := strconv.ParseInt(value, 10, 32); err != nil {
				return err
			} else {
				vf.SetInt(tmp)
			}
		case reflect.Int64:
			if len(tagArr) == 3 {
				format := tagArr[2]
				// parse time
				if format == "time" {
					if tmp, err := parseTime(value); err != nil {
						return err
					} else {
						vf.SetInt(tmp)
					}
				} else {
					return errors.New(fmt.Sprintf("unknown tag: %s in struct field: %s (support tags: \"time\")", format, tf.Name))
				}
			} else {
				if tmp, err := strconv.ParseInt(value, 10, 64); err != nil {
					return err
				} else {
					vf.SetInt(tmp)
				}
			}
		case reflect.Uint:
			if tmp, err := strconv.ParseUint(value, 10, 32); err != nil {
				return err
			} else {
				vf.SetUint(tmp)
			}
		case reflect.Uint8:
			if tmp, err := strconv.ParseUint(value, 10, 8); err != nil {
				return err
			} else {
				vf.SetUint(tmp)
			}
		case reflect.Uint16:
			if tmp, err := strconv.ParseUint(value, 10, 16); err != nil {
				return err
			} else {
				vf.SetUint(tmp)
			}
		case reflect.Uint32:
			if tmp, err := strconv.ParseUint(value, 10, 32); err != nil {
				return err
			} else {
				vf.SetUint(tmp)
			}
		case reflect.Uint64:
			if tmp, err := strconv.ParseUint(value, 10, 64); err != nil {
				return err
			} else {
				vf.SetUint(tmp)
			}
		case reflect.Slice:
			delim := ","
			if len(tagArr) > 2 {
				delim = tagArr[2]
			}
			strs := strings.Split(value, delim)
			sli := reflect.MakeSlice(tf.Type, 0, len(strs))
			for _, str := range strs {
				vv, err := getValue(tf.Type.Elem().String(), str)
				if err != nil {
					return err
				}
				sli = reflect.Append(sli, vv)
			}
			vf.Set(sli)
		case reflect.Map:
			delim := ","
			if len(tagArr) > 2 {
				delim = tagArr[2]
			}
			strs := strings.Split(value, delim)
			m := reflect.MakeMap(tf.Type)
			for _, str := range strs {
				mapStrs := strings.SplitN(str, "=", 2)
				if len(mapStrs) < 2 {
					return errors.New(fmt.Sprintf("error map: %s, must be split by \"=\"", str))
				}
				vk, err := getValue(tf.Type.Key().String(), mapStrs[0])
				if err != nil {
					return err
				}
				vv, err := getValue(tf.Type.Elem().String(), mapStrs[1])
				if err != nil {
					return err
				}
				m.SetMapIndex(vk, vv)
			}
			vf.Set(m)
		default:
			return errors.New(fmt.Sprintf("cannot unmarshall unsuported kind: %s into struct field: %s", vf.Kind().String(), tf.Name))
		}
	}
	return nil
}

// getValue parse String to the type "t" reflect.Value.
func getValue(t, v string) (reflect.Value, error) {
	var vv reflect.Value
	switch t {
	case "bool":
		d := parseBool(v)
		vv = reflect.ValueOf(d)
	case "int":
		d, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(int(d))
	case "int8":
		d, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(int8(d))
	case "int16":
		d, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(int16(d))
	case "int32":
		d, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(int32(d))
	case "int64":
		d, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(int64(d))
	case "uint":
		d, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(uint(d))
	case "uint8":
		d, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(uint8(d))
	case "uint16":
		d, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(uint16(d))
	case "uint32":
		d, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(uint32(d))
	case "uint64":
		d, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(uint64(d))
	case "float32":
		d, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(float32(d))
	case "float64":
		d, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return vv, err
		}
		vv = reflect.ValueOf(float64(d))
	case "string":
		vv = reflect.ValueOf(v)
	default:
		return vv, errors.New(fmt.Sprintf("unkown type: %s", t))
	}
	return vv, nil
}
