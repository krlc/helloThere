package helloThere

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	const geoIPFile = "test"

	temp, err := ioutil.TempFile(os.TempDir(), "config")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(temp.Name()); err != nil {
			t.Fatalf("can't remove temp file %v: %v", temp.Name(), err)
		}
	}()

	_, err = NewConfig("")
	if err == nil {
		t.Fatal("NewConfig() no error with empty path")
	}

	if _, err := temp.WriteString(fmt.Sprintf("geoipfile: %s\n", geoIPFile)); err != nil {
		t.Fatal(err)
	}

	t.Run("TestCorrectConfig", func(tt *testing.T) {
		config, err := NewConfig(temp.Name())
		if err != nil {
			tt.Fatal(err)
		}

		if config.GeoIPFile != geoIPFile {
			tt.Fatalf("NewConfig() GeoIPFile got = %v, want %v", config.GeoIPFile, geoIPFile)
		}
	})

	t.Run("TestIncorrectConfig", func(tt *testing.T) {
		if _, err := temp.WriteString("geoipfile"); err != nil {
			tt.Fatalf("can't write to temp file %v: %v", temp.Name(), err)
		}

		_, err := NewConfig(temp.Name())
		if err == nil {
			tt.Fatal("NewConfig() got nil error with erroneous config")
		}
	})
}
