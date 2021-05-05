package webp

import (
	"os"
	"testing"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"
)

var (
	exifFourCc = [4]byte{'E', 'X', 'I', 'F'}
)

func TestWebpParser_readHeader(t *testing.T) {
	filepath := GetTestImageFilepath()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	wp := NewWebpParser(f)

	fileSize, err := wp.readHeader()
	log.PanicIf(err)

	if fileSize != 403682 {
		t.Fatalf("File-size not correct: (%d)", fileSize)
	}
}

func TestWebpParser_readChunkHeader(t *testing.T) {
	filepath := GetTestImageFilepath()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	wp := NewWebpParser(f)

	_, err = wp.readHeader()
	log.PanicIf(err)

	fourCc, chunkSize, err := wp.readChunkHeader()
	log.PanicIf(err)

	if fourCc != [4]byte{'V', 'P', '8', 'X'} {
		t.Fatalf("Unexpected four-CC: [%04s]", fourCc[:])
	} else if chunkSize != 10 {
		t.Fatalf("Chunk-size not corret: (%d)", chunkSize)
	}
}

func TestWebpParser_enumerateChunks(t *testing.T) {
	filepath := GetTestImageFilepath()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	wp := NewWebpParser(f)

	exifFound := false
	visitorCb := func(fourCc [4]byte, dataGetter DataGetterFunc) (err error) {
		if fourCc != exifFourCc {
			return nil
		}

		exifData, err := dataGetter()
		log.PanicIf(err)

		im, err := exifcommon.NewIfdMappingWithStandard()
		log.PanicIf(err)

		ti := exif.NewTagIndex()

		_, _, err = exif.Collect(im, ti, exifData)
		log.PanicIf(err)

		exifFound = true

		return nil
	}

	err = wp.enumerateChunks(visitorCb)
	log.PanicIf(err)

	if exifFound != true {
		t.Fatalf("EXIF not found.")
	}
}
