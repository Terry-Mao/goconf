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
		t.FailNow()
	}
	section = "test"
	test := conf.Get(section)
	if test == nil {
		t.Errorf("not found section:\"%s\"", section)
		t.FailNow()
	}
	section = "test1"
	test1 := conf.Get(section)
	if test1 == nil {
		t.Errorf("not found section:\"%s\"", section)
		t.FailNow()
	}
	key := "id"
	if id, err := core.Int(key); err != nil {
		t.Errorf("core.Int(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if id != 1 {
			t.Errorf("%s not equals 1", key)
			t.FailNow()
		}
	}
	key = "col"
	if col, err := core.String(key); err != nil {
		t.Errorf("core.String(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if col != "goconf" {
			t.Errorf("%s not equals \"goconf\"", key)
			t.FailNow()
		}
	}
	key = "f"
	if f, err := core.Float(key); err != nil {
		t.Errorf("core.Float(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if f != 1.23 {
			t.Errorf("%s not equals 1.23", key)
			t.FailNow()
		}
	}
	key = "b"
	if b, err := core.Bool(key); err != nil {
		t.Errorf("core.Bool(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if !b {
			t.Errorf("%s not equals true", key)
			t.FailNow()
		}
	}
	key = "buf"
	if buf, err := core.MemSize(key); err != nil {
		t.Errorf("core.MemSize(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if buf != 1*1024*1024*1024 {
			t.Errorf("%s not equals 1*1024*1024*1024", key)
			t.FailNow()
		}
	}
	key = "sleep"
	if sleep, err := core.Duration(key); err != nil {
		t.Errorf("core.Duration(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if sleep != 10*time.Second {
			t.Errorf("%s not equals 10*Second", key)
			t.FailNow()
		}
	}
	key = "do"
	if do, err := core.String(key); err != nil {
		t.Errorf("core.String(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if do != "hehe" {
			t.Errorf("%s not equals \"hehe\"", key)
			t.FailNow()
		}
	}
	key = "id2"
	if id2, err := test.Int(key); err != nil {
		t.Errorf("test.Int(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if id2 != 2 {
			t.Errorf("%s not equals 2", key)
			t.FailNow()
		}
	}
	key = "id3"
	if id3, err := test1.Bool(key); err != nil {
		t.Errorf("test.Bool(\"%s\") failed (%s)", key, err.Error())
		t.FailNow()
	} else {
		if !id3 {
			t.Errorf("%s not equals false", key)
			t.FailNow()
		}
	}
	test1.Add("id4", "goconf baby", " hahah\n heihei,woshishei")
	save := "./examples/conf_reload.txt"
	if err := conf.Save(save); err != nil {
		t.Errorf("conf.Save(\"%s\") failed (%s)", save, err.Error())
		t.FailNow()
	}

	test1.Remove("id4")
	save = "./examples/conf_reload1.txt"
	if err := conf.Save(save); err != nil {
		t.Errorf("conf.Save(\"%s\") failed (%s)", save, err.Error())
		t.FailNow()
	}

	conf.Remove("test1")
	save = "./examples/conf_reload2.txt"
	if err := conf.Save(save); err != nil {
		t.Errorf("conf.Save(\"%s\") failed (%s)", save, err.Error())
		t.FailNow()
	}

	if _, err := conf.Reload(); err != nil {
		t.Errorf("conf.Reload() failed (%s)", err.Error())
		t.FailNow()
	}
	// test unmarshall
	tf := &TestConfig{}
	if err := conf.Unmarshal(tf); err != nil {
		t.Errorf("c.Unmarshal() failed (%s)", err.Error())
		t.FailNow()
	}
	if tf.ID != 1 {
		t.Errorf("TestConfig ID not equals 1")
		t.FailNow()
	}
	if tf.Col != "goconf" {
		t.Errorf("TestConfig Col not equals \"goconf\"")
		t.FailNow()
	}
	if len(tf.Arr) != 4 {
		t.Errorf("TestConfig Arr length not equals 4")
		t.FailNow()
	}
	if tf.Arr[0] != "1" {
		t.Errorf("TestConfig Arr[0] length not equals \"1\"")
		t.FailNow()
	}
	if tf.Arr[1] != "2" {
		t.Errorf("TestConfig Arr[1] length not equals \"2\"")
		t.FailNow()
	}
	if tf.Arr[2] != "3" {
		t.Errorf("TestConfig Arr[2] length not equals \"3\"")
		t.FailNow()
	}
	if tf.Arr[3] != "come on baby" {
		t.Errorf("TestConfig Arr[3] length not equals \"come on baby\"")
		t.FailNow()
	}
	if len(tf.Arr1) != 3 {
		t.Errorf("TestConfig Arr length not equals 4")
		t.FailNow()
	}
	if tf.Arr1[0] != 1 {
		t.Errorf("TestConfig Arr1[0] length not equals 1")
		t.FailNow()
	}
	if tf.Arr1[1] != 3 {
		t.Errorf("TestConfig Arr[1] length not equals 3")
		t.FailNow()
	}
	if tf.Arr1[2] != 4 {
		t.Errorf("TestConfig Arr1[2] length not equals 4")
		t.FailNow()
	}
	if tf.Test != 2*time.Hour {
		t.Errorf("TestConfig t_1 not equals 2 * time.Hour")
		t.FailNow()
	}
	if tf.Buf != 1*GB {
		t.Errorf("TestConfig t_1 not equals 1 * GB")
		t.FailNow()
	}
	if len(tf.M) != 2 {
		t.Errorf("TestConfig M length not equals 2")
		t.FailNow()
	}
	if v, ok := tf.M[1]; !ok {
		t.Errorf("TestConfig M no key 1")
		t.FailNow()
	} else {
		if v != "str" {
			t.Errorf("TestConfig M[1] not equals \"str\"")
			t.FailNow()
		}
	}
	if v, ok := tf.M[2]; !ok {
		t.Errorf("TestConfig M no key 2")
		t.FailNow()
	} else {
		if v != "str1" {
			t.Errorf("TestConfig M[2] not equals \"str1\"")
			t.FailNow()
		}
	}
}
