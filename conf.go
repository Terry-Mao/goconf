package goconf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	Byte = int64(1)
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

var (
	ErrNoKey   = errors.New("not found the config key")
	ErrDupFile = errors.New("duplicate config file parsed")
	ErrFormat  = errors.New("config file format error")
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
	for {
		line, err := rd.ReadString(CRLF)
		if err == io.EOF && len(line) == 0 {
			// parse file finish
			// all files parsed, break
			if len(files) <= 1 {
				break
			}
			// get the next file
			files = files[1:]
			f, err = os.Open(files[0])
			if err != nil {
				return err
			}
			defer f.Close()
			rd = bufio.NewReader(f)
			// reset section
			section = ""
			continue
		} else if err != nil {
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
				return ErrFormat
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
				return ErrFormat
			}
		}
		// get the spliter index
		idx := strings.Index(line, c.Spliter)
		if idx <= 0 {
			return ErrFormat
		}
		// get the key and value
		key := strings.TrimSpace(line[:idx])
		value := ""
		if len(line) > idx {
			value = strings.TrimSpace(line[idx+1:])
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
				return nil, ErrDupFile
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
			if _, err := f.WriteString(fmt.Sprintf("%s = %s%c", k, v, CRLF)); err != nil {
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

// String get config string value.
func (s *Section) String(key string) (string, error) {
	if v, ok := s.data[key]; ok {
		return v, nil
	} else {
		return "", ErrNoKey
	}
}

// Strings get config []string value.
func (s *Section) Strings(key, delim string) ([]string, error) {
	if v, ok := s.data[key]; ok {
		return strings.Split(v, delim), nil
	} else {
		return nil, ErrNoKey
	}
}

// Int get config int value.
func (s *Section) Int(key string) (int64, error) {
	if v, ok := s.data[key]; ok {
		return strconv.ParseInt(v, 10, 64)
	} else {
		return 0, ErrNoKey
	}
}

// Uint get config uint value.
func (s *Section) Uint(key string) (uint64, error) {
	if v, ok := s.data[key]; ok {
		return strconv.ParseUint(v, 10, 64)
	} else {
		return 0, ErrNoKey
	}
}

// Float get config float value.
func (s *Section) Float(key string) (float64, error) {
	if v, ok := s.data[key]; ok {
		return strconv.ParseFloat(v, 64)
	} else {
		return 0, ErrNoKey
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
		if v == "true" || v == "yes" || v == "1" || v == "y" || v == "enable" {
			return true, nil
		} else if v == "false" || v == "no" || v == "0" || v == "n" || v == "disable" {
			return false, nil
		} else {
			return false, nil
		}
	} else {
		return false, ErrNoKey
	}
}

// Byte get config byte number value.
//
// 1kb = 1k = 1024.
//
// 1mb = 1m = 1024 * 1024.
//
// 1gb = 1g = 1024 * 1024 * 1024.
func (s *Section) MemSize(key string) (int64, error) {
	if v, ok := s.data[key]; ok {
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
		return b * unit, nil
	} else {
		return 0, ErrNoKey
	}
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
		unit := time.Nanosecond
		subIdx := len(v)
		if strings.HasSuffix(v, "s") {
			unit = time.Second
			subIdx = subIdx - 1
		} else if strings.HasSuffix(v, "sec") {
			unit = time.Second
			subIdx = subIdx - 3
		} else if strings.HasSuffix(v, "m") {
			unit = time.Minute
			subIdx = subIdx - 1
		} else if strings.HasSuffix(v, "min") {
			unit = time.Minute
			subIdx = subIdx - 3
		} else if strings.HasSuffix(v, "h") {
			unit = time.Hour
			subIdx = subIdx - 1
		} else if strings.HasSuffix(v, "hour") {
			unit = time.Hour
			subIdx = subIdx - 4
		}
		b, err := strconv.ParseInt(v[:subIdx], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(b) * unit, nil
	} else {
		return 0, ErrNoKey
	}
}

// Keys return all the section keys.
func (s *Section) Keys() []string {
	keys := []string{}
	for k, _ := range s.data {
		keys = append(keys, k)
	}
	return keys
}
