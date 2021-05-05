package webp

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"
)

func TestWebpMediaParser_ParseBytes(t *testing.T) {
	wmp := NewWebpMediaParser()

	filepath := GetTestImageFilepath()

	// f, err := os.Open(filepath)
	// log.PanicIf(err)

	// defer f.Close()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	// f, err := os.Open(filepath)
	// log.PanicIf(err)

	// defer f.Close()

	mc, err := wmp.ParseBytes(data)
	log.PanicIf(err)

	rootIfd, _, err := mc.Exif()
	log.PanicIf(err)

	ii := rootIfd.IfdIdentity()
	if ii != exifcommon.IfdStandardIfdIdentity {
		t.Fatalf("Root IFD is not the correct type: [%s]", ii)
	}
}

func TestWebpMediaParser_ParseFile(t *testing.T) {
	wmp := NewWebpMediaParser()

	filepath := GetTestImageFilepath()

	mc, err := wmp.ParseFile(filepath)
	log.PanicIf(err)

	rootIfd, _, err := mc.Exif()
	log.PanicIf(err)

	ii := rootIfd.IfdIdentity()
	if ii != exifcommon.IfdStandardIfdIdentity {
		t.Fatalf("Root IFD is not the correct type: [%s]", ii)
	}
}

func TestWebpMediaParser_Parse(t *testing.T) {
	wmp := NewWebpMediaParser()

	filepath := GetTestImageFilepath()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	mc, err := wmp.Parse(f, 0)
	log.PanicIf(err)

	rootIfd, _, err := mc.Exif()
	log.PanicIf(err)

	ii := rootIfd.IfdIdentity()
	if ii != exifcommon.IfdStandardIfdIdentity {
		t.Fatalf("Root IFD is not the correct type: [%s]", ii)
	}
}
