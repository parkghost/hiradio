package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testfile() (string, error) {
	dir, err := ioutil.TempDir("", "hiradio")
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, "config.json")
	return path, nil
}

func TestSetAndGetString(t *testing.T) {
	testKey := "player"
	testValue := "/usr/bin/vlc"
	file, err := testfile()
	if err != nil {
		t.Fatalf("unexpected error on testfile: %s", err)
	}

	c := New()
	c.Set(testKey, testValue)
	err = SaveTo(file, c)
	if err != nil {
		t.Fatalf("unexpected error on SaveTo: %s", err)
	}

	c, err = From(file)
	if err != nil {
		t.Fatalf("unexpected error on From: %s", err)
	}

	got := c.GetString(testKey, "")
	if got != testValue {
		t.Fatalf("got %s, want %s", got, testValue)
	}
}

func TestSetAndGetInt(t *testing.T) {
	testKey := "port"
	testValue := 1077
	file, err := testfile()
	if err != nil {
		t.Fatalf("unexpected error on testfile: %s", err)
	}

	c := New()
	c.Set(testKey, testValue)
	err = SaveTo(file, c)
	if err != nil {
		t.Fatalf("unexpected error on SaveTo: %s", err)
	}

	c, err = From(file)
	if err != nil {
		t.Fatalf("unexpected error on From: %s", err)
	}

	got := c.GetInt(testKey, -1)
	if got != testValue {
		t.Fatalf("got %d, want %d", got, testValue)
	}
}

func TestGetStringDefaultValue(t *testing.T) {
	testKey := "player"
	testValue := "/usr/bin/vlc"

	c := New()
	want := "vvvvv"
	got := c.GetString(testKey, want)
	if got != want {
		t.Fatalf("got %s, want %s", got, testValue)
	}

	c.Set(testKey, 99999)
	got = c.GetString(testKey, want)
	if got != want {
		t.Fatalf("got %s, want %s", got, testValue)
	}
}

func TestGetIntDefaultValue(t *testing.T) {
	testKey := "port"
	testValue := 1077

	c := New()
	want := 1234
	got := c.GetInt(testKey, want)
	if got != want {
		t.Fatalf("got %s, want %s", got, testValue)
	}

	c.Set(testKey, "AAAAA")
	got = c.GetInt(testKey, want)
	if got != want {
		t.Fatalf("got %s, want %s", got, testValue)
	}
}

func TestConfigChange(t *testing.T) {
	testKey := "player"
	testValue := "/usr/bin/vlc"

	c := New()
	c.Set(testKey, testValue)
	if !c.changed {
		t.Fatalf("set new kv, config.changed should be true")
	}

	// reset the flag
	c.changed = false
	c.Set(testKey, testValue)
	if c.changed {
		t.Fatalf("set same kv, config.changed should be false")
	}

	c.Set(testKey, "vvvvv")
	if !c.changed {
		t.Fatalf("set new value to existed key, config.changed should be true")
	}
}

func TestFromError(t *testing.T) {
	file, err := testfile()
	if err != nil {
		t.Fatalf("unexpected error on testfile: %s", err)
	}

	_, err = From(file)
	if err != ErrEmptyFile {
		t.Fatalf("expected ErrEmptyFile on From: %s", err)
	}

	err = ioutil.WriteFile(file, []byte("bad format"), 0644)
	if err != nil {
		t.Fatalf("unexpected error on WriteFile: %s", err)
	}
	_, err = From(file)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError when parsing bad format of data: %#+v", err)
	}
}

func TestSaveToError(t *testing.T) {
	testKey := "player"
	testValue := "/usr/bin/vlc"
	file, err := testfile()
	if err != nil {
		t.Fatalf("unexpected error on testfile: %s", err)
	}

	c := New()
	c.Set(testKey, testValue)

	// reset flag
	c.changed = false
	err = SaveTo(file, c)
	if err != nil {
		t.Fatalf("unexpected error on SaveTo: %s", err)
	}

	_, err = os.Open(file)
	if err == nil {
		t.Fatal("Open should have failed")
	}
}
