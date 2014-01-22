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
	includeLen = len(Include)
	// memory unit
	Byte = int64(1)
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
	// time unit
	Nanosecond = int64(time.Nanosecond)
	Second     = int64(time.Second)
	Minute     = int64(time.Minute)
	Hour       = int64(time.Hour)
)

var (
	ErrNotFoundspliter = errors.New("not found spliter")
	ErrNotFoundKey     = errors.New("not found the config key")
	ErrDuplicateFile   = errors.New("duplicate config file parsed")
	ErrIncludeFile     = errors.New("include file format error")
	ErrBooleanValue    = errors.New("boolean string error")
)

// Config is the key-value configuration object.
type Config struct {
	data    map[string]string
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
	c.data = map[string]string{}
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
	for {
		line, err := rd.ReadString(CRLF)
		if err == io.EOF {
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
				return ErrIncludeFile
			}
		}
		// get the spliter index
		idx := strings.Index(line, c.Spliter)
		if idx <= 0 {
			return ErrNotFoundspliter
		}
		// get the key and value
		key := strings.TrimSpace(line[:idx])
		value := ""
		if len(line) > idx {
			value = strings.TrimSpace(line[idx+1:])
		}
		// store the key-value config
		c.data[key] = value
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
				return nil, ErrDuplicateFile
			}
			files = append(files, file)
			// save parse file
			fileMap[file] = true
		}
	}
	return files, nil
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
	for k, v := range c.data {
		if _, err := f.WriteString(fmt.Sprintf("%s = %s%c", k, v, CRLF)); err != nil {
			return err
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

// Add add a new key-value configuration.
func (c *Config) Add(k, v string) {
	c.data[k] = v
}

// Remove remove the specified key configuration.
func (c *Config) Remove(k string) {
	delete(c.data, k)
}

// String get config string value.
func (c *Config) String(key string) (string, error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	} else {
		return "", ErrNotFoundKey
	}
}

// Int get config int value.
func (c *Config) Int(key string) (int64, error) {
	if v, ok := c.data[key]; ok {
		return strconv.ParseInt(v, 10, 64)
	} else {
		return 0, ErrNotFoundKey
	}
}

// Uint get config uint value.
func (c *Config) Uint(key string) (uint64, error) {
	if v, ok := c.data[key]; ok {
		return strconv.ParseUint(v, 10, 64)
	} else {
		return 0, ErrNotFoundKey
	}
}

// Float get config float value.
func (c *Config) Float(key string) (float64, error) {
	if v, ok := c.data[key]; ok {
		return strconv.ParseFloat(v, 64)
	} else {
		return 0, ErrNotFoundKey
	}
}

// Bool get config boolean value.
//
// "yes", "1", "y", "true", "enable" means true.
//
// "no", "0", "n", "false", "disable" means false.
func (c *Config) Bool(key string) (bool, error) {
	if v, ok := c.data[key]; ok {
		v = strings.ToLower(v)
		if v == "true" || v == "yes" || v == "1" || v == "y" || v == "enable" {
			return true, nil
		} else if v == "false" || v == "no" || v == "0" || v == "n" || v == "disable" {
			return false, nil
		} else {
			return false, ErrBooleanValue
		}
	} else {
		return false, ErrNotFoundKey
	}
}

// Byte get config byte number value.
//
// 1kb = 1k = 1024.
//
// 1mb = 1m = 1024 * 1024.
//
// 1gb = 1g = 1024 * 1024 * 1024.
func (c *Config) MemSize(key string) (int64, error) {
	if v, ok := c.data[key]; ok {
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
		return 0, ErrNotFoundKey
	}
}

// Second get config second value.
//
// 1s = 1sec = 1.
//
// 1m = 1min = 60.
//
// 1h = 1hour = 60 * 60.
func (c *Config) Duration(key string) (int64, error) {
	if v, ok := c.data[key]; ok {
		unit := Nanosecond
		subIdx := len(v)
		if strings.HasSuffix(v, "s") {
			unit = Second
			subIdx = subIdx - 1
		} else if strings.HasSuffix(v, "sec") {
			unit = Second
			subIdx = subIdx - 3
		} else if strings.HasSuffix(v, "m") {
			unit = Minute
			subIdx = subIdx - 1
		} else if strings.HasSuffix(v, "min") {
			unit = Minute
			subIdx = subIdx - 3
		} else if strings.HasSuffix(v, "h") {
			unit = Hour
			subIdx = subIdx - 1
		} else if strings.HasSuffix(v, "hour") {
			unit = Hour
			subIdx = subIdx - 4
		}
		b, err := strconv.ParseInt(v[:subIdx], 10, 64)
		if err != nil {
			return 0, err
		}
		return b * unit, nil
	} else {
		return 0, ErrNotFoundKey
	}
}
