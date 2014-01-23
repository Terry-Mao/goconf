package goconf

import (
	"testing"
)

var (
	conf *Config
)

func init() {
	file := "./examples/conf_test.txt"
	conf = New()
	if err := conf.Parse(file); err != nil {
		panic(err)
	}
}

func TestSection(t *testing.T) {
	section := "core"
	core := conf.Get(section)
	if core == nil {
		t.Errorf("not found section:\"%s\"", section)
	}
	section = "test"
	test := conf.Get(section)
	if core == nil {
		t.Errorf("not found section:\"%s\"", section)
	}
	section = "test1"
	test1 := conf.Get(section)
	if core == nil {
		t.Errorf("not found section:\"%s\"", section)
	}
	key := "id"
	if id, err := core.Int(key); err != nil {
		t.Errorf("core.Int(\"%s\") failed (%s)", key, err.Error())
	} else {
		if id != 1 {
			t.Errorf("%s not equals 1", key)
		}
	}
	key = "col"
	if col, err := core.String(key); err != nil {
		t.Errorf("core.String(\"%s\") failed (%s)", key, err.Error())
	} else {
		if col != "goconf" {
			t.Errorf("%s not equals \"goconf\"", key)
		}
	}
	key = "f"
	if f, err := core.Float(key); err != nil {
		t.Errorf("core.Float(\"%s\") failed (%s)", key, err.Error())
	} else {
		if f != 1.23 {
			t.Errorf("%s not equals 1.23", key)
		}
	}
	key = "b"
	if b, err := core.Bool(key); err != nil {
		t.Errorf("core.Bool(\"%s\") failed (%s)", key, err.Error())
	} else {
		if !b {
			t.Errorf("%s not equals true", key)
		}
	}
	key = "buf"
	if buf, err := core.MemSize(key); err != nil {
		t.Errorf("core.MemSize(\"%s\") failed (%s)", key, err.Error())
	} else {
		if buf != 1*1024*1024*1024 {
			t.Errorf("%s not equals 1*1024*1024*1024", key)
		}
	}
	key = "sleep"
	if sleep, err := core.Duration(key); err != nil {
		t.Errorf("core.Duration(\"%s\") failed (%s)", key, err.Error())
	} else {
		if sleep != 10*Second {
			t.Errorf("%s not equals 10*Second", key)
		}
	}
	key = "id2"
	if id2, err := test.Int(key); err != nil {
		t.Errorf("test.Int(\"%s\") failed (%s)", key, err.Error())
	} else {
		if id2 != 2 {
			t.Errorf("%s not equals 2", key)
		}
	}
	key = "id3"
	if id3, err := test1.Bool(key); err != nil {
		t.Errorf("test.Bool(\"%s\") failed (%s)", key, err.Error())
	} else {
		if !id3 {
			t.Errorf("%s not equals false", key)
		}
	}
}

func TestSave(t *testing.T) {
}

func TestReload(t *testing.T) {
}
