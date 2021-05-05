package webp

import (
	"errors"
	"io"

	"encoding/binary"

	"github.com/dsoprea/go-logging/v2"
)

var (
	// ErrNoRiffFileHeader indicates that the file-header failed.
	ErrNoRiffFileHeader = errors.New("not on a RIFF file-header")

	// ErrNoWebpSignature indicates that the WEBP signature following the file-
	// header failed.
	ErrNoWebpSignature = errors.New("not a WEBP file")
)

var (
	// DefaultEndianness is the endianness of a standard RIFF stream.
	DefaultEndianness = binary.LittleEndian
)

var (
	riffFileHeaderBytes = [4]byte{'R', 'I', 'F', 'F'}
	webpBytes           = [4]byte{'W', 'E', 'B', 'P'}
	exifFourCc          = [4]byte{'E', 'X', 'I', 'F'}
)

// WebpParser parses WEBP RIFF streams.
type WebpParser struct {
	r io.Reader
}

// NewWebpParser returns a WebpParser.
func NewWebpParser(r io.Reader) *WebpParser {
	return &WebpParser{
		r: r,
	}
}

// readHeader returns the file-header if we're sitting on the first byte.
func (wp *WebpParser) readHeader() (fileSize int64, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = errRaw.(error)
		}
	}()

	var fileHeaderBytes [4]byte
	_, err = io.ReadFull(wp.r, fileHeaderBytes[:])
	log.PanicIf(err)

	if fileHeaderBytes != riffFileHeaderBytes {
		return 0, ErrNoRiffFileHeader
	}

	var fileSizeRaw uint32
	err = binary.Read(wp.r, DefaultEndianness, &fileSizeRaw)
	log.PanicIf(err)

	var webpSignatureBytes [4]byte
	_, err = io.ReadFull(wp.r, webpSignatureBytes[:])
	log.PanicIf(err)

	if webpSignatureBytes != webpBytes {
		return 0, ErrNoWebpSignature
	}

	return int64(fileSizeRaw), nil
}

// readChunkHeader returns the chunk-information if we're sitting on the first
// byte of a chunk.
func (wp *WebpParser) readChunkHeader() (fourCc [4]byte, chunkSize int64, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = errRaw.(error)
		}
	}()

	_, err = io.ReadFull(wp.r, fourCc[:])
	log.PanicIf(err)

	var chunkSizeRaw uint32
	err = binary.Read(wp.r, DefaultEndianness, &chunkSizeRaw)
	log.PanicIf(err)

	return fourCc, int64(chunkSizeRaw), nil
}

// DataGetterFunc is a lazy-getter for the payload data.
type DataGetterFunc func() (data []byte, err error)

// ChunkVisitorFunc is a callback that receives each chunk.
type ChunkVisitorFunc func(fourCc [4]byte, dataGetter DataGetterFunc) (err error)

// enumerateChunks enumerates each sequential chunk. Takes an optional callback
// (if no callback is given, execution essentially becomes simple validation).
func (wp *WebpParser) enumerateChunks(chunkVisitorCb ChunkVisitorFunc) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = errRaw.(error)
		}
	}()

	fileSize, err := wp.readHeader()
	log.PanicIf(err)

	remainingBytes := fileSize - 4

	for remainingBytes > 0 {
		fourCc, chunkSize, err := wp.readChunkHeader()
		log.PanicIf(err)

		hasRead := false
		dataGetter := func() (data []byte, err error) {
			defer func() {
				if errRaw := recover(); errRaw != nil {
					err = errRaw.(error)
				}
			}()

			data = make([]byte, chunkSize)
			_, err = io.ReadFull(wp.r, data)
			log.PanicIf(err)

			hasRead = true
			return data, nil
		}

		if chunkVisitorCb != nil {
			err := chunkVisitorCb(fourCc, dataGetter)
			log.PanicIf(err)
		}

		// Need to do this before we can process the padding, so, if the caller
		// didn't provide a callback, or if they did but didn't actually
		// retrieve the data, we'll read it here.
		if hasRead == false {
			_, err := dataGetter()
			log.PanicIf(err)
		}

		remainingBytes = remainingBytes - 8 - chunkSize

		// If there is an odd number of bytes, there will be one padding byte.
		if chunkSize%2 == 1 {

			padding := make([]byte, 1)
			_, err := io.ReadFull(wp.r, padding)
			log.PanicIf(err)

			remainingBytes--
		}
	}

	return nil
}
