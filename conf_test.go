package goconf

import (
	"testing"
)

func TestConf(t *testing.T) {
	file := "./examples/conf_test.txt"
	if conf, err := New(file); err != nil {
		t.Errorf("New(\"%s\") failed (%s)", file, err.Error())
	} else {
		key := "id"
		if id, err := conf.String(key); err != nil {
			t.Errorf("conf.String(\"%s\") failed (%s)", key, err.Error())
		} else {
			if id != "1" {
				t.Errorf("config key \"%s\" value not equals \"1\"", key)
			}
		}

		key = "id2"
		if id, err := conf.Int(key); err != nil {
			t.Errorf("conf.Int(\"%s\") failed (%s)", key, err.Error())
		} else {
			if id != 2 {
				t.Errorf("config key \"%s\" value not equals 2", key)
			}
		}

		key = "id3"
		if id, err := conf.String(key); err != nil {
			t.Errorf("conf.String(\"%s\") failed (%s)", key, err.Error())
		} else {
			if id != "yes" {
				t.Errorf("config key \"%s\" value not equals \"yes\"", key)
			}
		}

		key = "f"
		if f, err := conf.Float(key); err != nil {
			t.Errorf("conf.String(\"%s\") failed (%s)", key, err.Error())
		} else {
			if f != 1.23 {
				t.Errorf("config key \"%s\" value not equals 1.23", key)
			}
		}

		key = "test"
		if test, err := conf.Bool(key); err != nil {
			t.Errorf("conf.Bool(\"%s\") failed (%s)", key, err.Error())
		} else {
			if test != false {
				t.Errorf("config key \"%s\" value not equals false", key)
			}
		}

		key = "buf"
		if buf, err := conf.Byte(key); err != nil {
			t.Errorf("conf.Byte(\"%s\") failed (%s)", key, err.Error())
		} else {
			if buf != 1*GB {
				t.Errorf("config key \"%s\" value not equals %d", key, 1*GB)
			}
		}

		key = "sleep"
		if buf, err := conf.Duration(key); err != nil {
			t.Errorf("conf.Duration(\"%s\") failed (%s)", key, err.Error())
		} else {
			if buf != 10*Second {
				t.Errorf("config key \"%s\" value not equals %d", key, 10*Second)
			}
		}

		if len(conf.data) != 7 {
			t.Errorf("parse config file \"%s\" failed, map length error", file)
		}
	}
}
