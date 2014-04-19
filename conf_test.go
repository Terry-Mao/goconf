package goconf

import (
	"testing"
	"time"
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

type TestConfig struct {
	ID     int            `goconf:"core:id"`
	Col    string         `goconf:"core:col"`
	Ignore int            `goconf:"-"`
	Arr    []string       `goconf:"core:arr:,"`
	Arr1   []int          `goconf:"core:arr1:,"`
	Test   time.Duration  `goconf:"core:t_1:time"`
	Buf    int            `goconf:"core:buf:memory"`
	M      map[int]string `goconf:"core:map:,"`
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
		if sleep != 10*time.Second {
			t.Errorf("%s not equals 10*Second", key)
		}
	}
	key = "do"
	if do, err := core.String(key); err != nil {
		t.Errorf("core.String(\"%s\") failed (%s)", key, err.Error())
	} else {
		if do != "hehe" {
			t.Errorf("%s not equals \"hehe\"", key)
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
	test1.Add("id4", "goconf baby")
	save := "./examples/conf_reload.txt"
	if err := conf.Save(save); err != nil {
		t.Errorf("conf.Save(\"%s\") failed (%s)", save, err.Error())
	}

	if _, err := conf.Reload(); err != nil {
		t.Errorf("conf.Reload() failed (%s)", err.Error())
	}
	// test unmarshall
	tf := &TestConfig{}
	if err := conf.Unmarshal(tf); err != nil {
		t.Errorf("c.Unmarshal() failed (%s)", err.Error())
	}
	if tf.ID != 1 {
		t.Errorf("TestConfig ID not equals 1")
	}
	if tf.Col != "goconf" {
		t.Errorf("TestConfig Col not equals \"goconf\"")
	}
	if len(tf.Arr) != 4 {
		t.Errorf("TestConfig Arr length not equals 4")
	}
	if tf.Arr[0] != "1" {
		t.Errorf("TestConfig Arr[0] length not equals \"1\"")
	}
	if tf.Arr[1] != "2" {
		t.Errorf("TestConfig Arr[1] length not equals \"2\"")
	}
	if tf.Arr[2] != "3" {
		t.Errorf("TestConfig Arr[2] length not equals \"3\"")
	}
	if tf.Arr[3] != "come on baby" {
		t.Errorf("TestConfig Arr[3] length not equals \"come on baby\"")
	}
	if len(tf.Arr1) != 3 {
		t.Errorf("TestConfig Arr length not equals 4")
	}
	if tf.Arr1[0] != 1 {
		t.Errorf("TestConfig Arr1[0] length not equals 1")
	}
	if tf.Arr1[1] != 3 {
		t.Errorf("TestConfig Arr[1] length not equals 3")
	}
	if tf.Arr1[2] != 4 {
		t.Errorf("TestConfig Arr1[2] length not equals 4")
	}
	if tf.Test != 2*time.Hour {
		t.Errorf("TestConfig t_1 not equals 2 * time.Hour")
	}
	if tf.Buf != 1*GB {
		t.Errorf("TestConfig t_1 not equals 1 * GB")
	}
	if len(tf.M) != 2 {
		t.Errorf("TestConfig M length not equals 2")
	}
	if v, ok := tf.M[1]; !ok {
		t.Errorf("TestConfig M no key 1")
	} else {
		if v != "str" {
			t.Errorf("TestConfig M[1] not equals \"str\"")
		}
	}
	if v, ok := tf.M[2]; !ok {
		t.Errorf("TestConfig M no key 2")
	} else {
		if v != "str1" {
			t.Errorf("TestConfig M[2] not equals \"str1\"")
		}
	}
}
