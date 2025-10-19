// Package audio handles some manipulations with .wav files made up according to PCM standard
package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Fixed values applicable to PCM audio files
const (
	chunkID              = "RIFF"
	format               = "WAVE"
	subChunkID           = "fmt "
	subChunk1Size uint32 = 16
	audioFormat   uint16 = 1
	subChunk2ID          = "data"
	mono, stereo  uint16 = 1, 2
)

type meta struct {
	// Contains "RIFF" in ASCII
	ChunkID [4]byte
	// Audio file size excluding 8 bytes
	ChunkSize uint32
	// Contains "WAVE" in ASCII
	Format [4]byte
	// Contains "fmt " in ASCII
	SubChunk1ID [4]byte
	// Subchunk1 size, (16 bytes for PCM)
	SubChunk1Size uint32
	// format of audio, (1 for PCM)
	AudioFormat uint16
	// 1 - Mono, 2 - stereo
	NumChannels uint16
	SampleRate  uint32
	ByteRate    uint32
	// the amount of bytes for a single sample including all NumChannels
	BlockAlign    uint16
	BitsPerSample uint16
	// Contains "data" in ASCII
	SubChunk2ID [4]byte
	// Raw audio data in bytes that goes after this header
	SubChunk2Size uint32
}

// PCMAudioMetadata is a specification of WAVE PCM audio files. Contains 44 bytes
type PCMAudioMetadata struct {
	meta
	valid bool
}

func (p *PCMAudioMetadata) Validate() (bool, error) {
	switch {
	case string(p.meta.ChunkID[:]) != chunkID:
		return false, errors.New("RIFF header doesn't match")
	case string(p.meta.Format[:]) != format:
		return false, errors.New("audio format doesn't match")
	case string(p.meta.SubChunk1ID[:]) != subChunkID:
		return false, errors.New("subchunk fmt doesn't match")
	case p.meta.SubChunk1Size != subChunk1Size || p.meta.AudioFormat != audioFormat:
		return false, errors.New("provided data is not PCM audio")
	case string(p.meta.SubChunk2ID[:]) != subChunk2ID:
		return false, errors.New("data header doesn't match")
	default:
		p.valid = true
		return true, nil
	}
}

func (p *PCMAudioMetadata) Merge(other *PCMAudioMetadata, output io.Writer, left, right io.ReadSeeker) error {
	if !p.valid || !other.valid {
		return errors.New("validate each audio file first")
	}

	if p.meta.SampleRate != other.meta.SampleRate {
		return errors.New("rate of both audio files must match")
	}

	if p.meta.BitsPerSample != other.meta.BitsPerSample {
		return errors.New("bits per sample of both audio files must match")
	}

	newMeta := meta{
		ChunkID:       [4]byte([]byte(chunkID)),
		ChunkSize:     36 + p.meta.SubChunk2Size + other.meta.SubChunk2Size,
		Format:        [4]byte([]byte(format)),
		SubChunk1ID:   [4]byte([]byte(subChunkID)),
		SubChunk1Size: subChunk1Size,
		AudioFormat:   audioFormat,
		NumChannels:   stereo,
		SampleRate:    p.meta.SampleRate,
		ByteRate:      p.meta.SampleRate * uint32(stereo) * uint32(p.meta.BitsPerSample/8),
		BlockAlign:    stereo * (p.meta.BitsPerSample / 8),
		BitsPerSample: p.meta.BitsPerSample,
		SubChunk2ID:   [4]byte([]byte(subChunk2ID)),
		SubChunk2Size: p.meta.SubChunk2Size + other.meta.SubChunk2Size,
	}
	err := binary.Write(output, binary.LittleEndian, newMeta)
	if err != nil {
		return fmt.Errorf("could not write headers: %v", err)
	}

	_, err = left.Seek(44, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek left: %v", err)
	}

	_, err = right.Seek(44, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek right: %v", err)
	}

	sampleSize := int(p.meta.BitsPerSample / 8)
	bufLeft := make([]byte, sampleSize)
	bufRight := make([]byte, sampleSize)

	for {
		nL, errL := left.Read(bufLeft)
		if errL != nil && errL != io.EOF {
			return fmt.Errorf("error reading left: %v", errL)
		}

		nR, errR := right.Read(bufRight)
		if errR != nil && errR != io.EOF {
			return fmt.Errorf("error reading right: %v", errR)
		}

		// left and right data are consumed
		if (errL == io.EOF || nL == 0) && (errR == io.EOF || nR == 0) {
			break
		}

		// writing samples and adding zeros if files are not of the same size in order to align them
		if nL == sampleSize {
			_, errW := output.Write(bufLeft)
			if errW != nil {
				return fmt.Errorf("error writing left sample: %v", errW)
			}
		} else {
			zeros := make([]byte, sampleSize)
			_, errW := output.Write(zeros)
			if errW != nil {
				return fmt.Errorf("error writing zero sample for left: %v", errW)
			}
		}

		if nR == sampleSize {
			_, errW := output.Write(bufRight)
			if errW != nil {
				return fmt.Errorf("error writing right sample: %v", errW)
			}
		} else {
			zeros := make([]byte, sampleSize)
			_, errW := output.Write(zeros)
			if errW != nil {
				return fmt.Errorf("error writing zero sample for right: %v", errW)
			}
		}

		if errL == io.EOF || errR == io.EOF {
			break
		}
	}

	return nil
}

// NewPCMAudioMetadata receives an io.Reader param and completes its fields
func NewPCMAudioMetadata(data io.Reader) (*PCMAudioMetadata, error) {
	var p PCMAudioMetadata
	if err := binary.Read(data, binary.LittleEndian, &p.meta); err != nil {
		return nil, fmt.Errorf("could not read and parse the recieved data: %v", err.Error())
	}
	return &p, nil
}
