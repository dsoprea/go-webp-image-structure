package webp

import (
	"bytes"
	"io"
	"os"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"
	"github.com/dsoprea/go-utility/v2/image"
)

// WebpMediaParser is a `riimage.MediaParser` that knows how to parse JPEG
// images.
type WebpMediaParser struct {
}

// NewWebpMediaParser returns a new WebpMediaParser.
func NewWebpMediaParser() *WebpMediaParser {

	// TODO(dustin): Add test

	return new(WebpMediaParser)
}

type RawExifData []byte

// Exif returns the EXIF's root IFD.
func (ref RawExifData) Exif() (rootIfd *exif.Ifd, data []byte, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	im, err := exifcommon.NewIfdMappingWithStandard()
	log.PanicIf(err)

	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, ref)
	log.PanicIf(err)

	return index.RootIfd, ref, nil
}

// Parse parses a JPEG uses an `io.ReadSeeker`. Even if it fails, it will return
// the list of segments encountered prior to the failure.
func (wmp *WebpMediaParser) Parse(rs io.ReadSeeker, size int) (mc riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	wp := NewWebpParser(rs)

	visitorCb := func(fourCc [4]byte, dataGetter DataGetterFunc) (err error) {
		if fourCc != exifFourCc {
			return nil
		}

		exifData, err := dataGetter()
		log.PanicIf(err)

		mc = RawExifData(exifData)

		return nil
	}

	err = wp.enumerateChunks(visitorCb)
	log.PanicIf(err)

	if mc != nil {
		return mc, nil
	}

	return nil, exif.ErrNoExif
}

// ParseFile parses a JPEG file. Even if it fails, it will return the list of
// segments encountered prior to the failure.
func (wmp *WebpMediaParser) ParseFile(filepath string) (mc riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): Add test

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	stat, err := f.Stat()
	log.PanicIf(err)

	size := stat.Size()

	mc, err = wmp.Parse(f, int(size))
	log.PanicIf(err)

	return mc, nil
}

// ParseBytes parses a JPEG byte-slice. Even if it fails, it will return the
// list of segments encountered prior to the failure.
func (wmp *WebpMediaParser) ParseBytes(data []byte) (mc riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	br := bytes.NewReader(data)

	mc, err = wmp.Parse(br, len(data))
	log.PanicIf(err)

	return mc, nil
}

// LooksLikeFormat indicates whether the data looks like a JPEG image.
func (wmp *WebpMediaParser) LooksLikeFormat(data []byte) bool {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			log.Panic(errRaw.(error))
		}
	}()

	br := bytes.NewReader(data)
	wp := NewWebpParser(br)

	_, err := wp.readHeader()
	return err == nil
}

var (
	// Enforce interface conformance.
	_ riimage.MediaParser = new(WebpMediaParser)
)
