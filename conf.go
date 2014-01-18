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
)

const (
	crlf    = '\n'
	commet  = "#"
	splite  = "="
	include = "include"
)

var (
	ErrNotFoundSplite = errors.New(fmt.Sprintf("not found splite:\"%s\"", splite))
	ErrNotFoundKey    = errors.New("not found the config key")
	ErrDuplicateFile  = errors.New("duplicate config file parsed")
	ErrIncludeFile    = errors.New("include file format error")

	includeLen = len(include)
)

// Config is the key-value configuration object.
type Config struct {
	data map[string]string
	file string
}

// New return a new Config which parse the specified file.
func New(file string) (*Config, error) {
	c := &Config{data: map[string]string{}, file: file}
	// open config file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	rd := bufio.NewReader(f)
	files := []string{file}
	fileMap := map[string]bool{file: true}
	for {
		line, err := rd.ReadString(crlf)
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
				return nil, err
			}

			defer f.Close()
			rd = bufio.NewReader(f)
			continue
		} else if err != nil {
			return nil, err
		}

		// trim space
		line = strings.TrimSpace(line)
		// ignore blank line
		if line == "" {
			continue
		}

		// ignore commet line
		if strings.HasPrefix(line, commet) {
			continue
		}

		// handle include
		if strings.HasPrefix(line, include) {
			if len(line) > includeLen {
				// add other config files
				newFiles, err := includeFiles(strings.TrimSpace(line[includeLen+1:]), fileMap)
				if err != nil {
					return nil, err
				}

				files = append(files, newFiles...)
				continue
			} else {
				return nil, ErrIncludeFile
			}
		}

		// get the spliter index
		idx := strings.Index(line, splite)
		if idx <= 0 {
			return nil, ErrNotFoundSplite
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

	return c, nil
}

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

// String get config string value
func (c *Config) String(key string) (string, error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	} else {
		return "", ErrNotFoundKey
	}
}

// Int get config int value
func (c *Config) Int(key string) (int, error) {
	if v, ok := c.data[key]; ok {
		return strconv.Atoi(v)
	} else {
		return 0, ErrNotFoundKey
	}
}
