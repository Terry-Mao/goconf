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

		if len(conf.data) != 3 {
			t.Errorf("parse config file \"%s\" failed, map length error", file)
		}
	}
}
