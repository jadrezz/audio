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

// headers are the fixed fields of a WAV file
type headers struct {
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

// PCMAudio is a specification of WAVE PCM audio files. Contains 44 bytes
type PCMAudio struct {
	headers
	data  io.ReadSeeker
	valid bool
}

func (p *PCMAudio) Validate() (bool, error) {
	switch {
	case string(p.ChunkID[:]) != chunkID:
		return false, errors.New("RIFF header doesn't match")
	case string(p.Format[:]) != format:
		return false, errors.New("audio format doesn't match")
	case string(p.SubChunk1ID[:]) != subChunkID:
		return false, errors.New("subchunk fmt doesn't match")
	case p.SubChunk1Size != subChunk1Size || p.AudioFormat != audioFormat:
		return false, errors.New("provided data is not PCM audio")
	case string(p.SubChunk2ID[:]) != subChunk2ID:
		return false, errors.New("data header doesn't match")
	default:
		p.valid = true
		return true, nil
	}
}

// Merge creates a single stereo file made from 2 mono files.
//
// Accepts other PCMAudio and io.Writer to write the result
// If something goes wrong, returns an error
func (p *PCMAudio) Merge(other *PCMAudio, output io.Writer) error {
	if !p.valid || !other.valid {
		return errors.New("validate each audio file first")
	}

	if p.SampleRate != other.SampleRate {
		return errors.New("rate of both audio files must match")
	}

	if p.BitsPerSample != other.BitsPerSample {
		return errors.New("bits per sample of both audio files must match")
	}

	newChunkSize := 36 + p.SubChunk2Size + other.SubChunk2Size
	newBlockAlign := stereo * (p.BitsPerSample / 8)
	newSubChunk2Size := p.SubChunk2Size + other.SubChunk2Size

	newHeaders := headers{
		ChunkID:       [4]byte([]byte(chunkID)),
		ChunkSize:     newChunkSize,
		Format:        [4]byte([]byte(format)),
		SubChunk1ID:   [4]byte([]byte(subChunkID)),
		SubChunk1Size: subChunk1Size,
		AudioFormat:   audioFormat,
		NumChannels:   stereo,
		SampleRate:    p.SampleRate,
		ByteRate:      p.SampleRate * uint32(newBlockAlign),
		BlockAlign:    newBlockAlign,
		BitsPerSample: p.BitsPerSample,
		SubChunk2ID:   [4]byte([]byte(subChunk2ID)),
		SubChunk2Size: newSubChunk2Size,
	}
	err := binary.Write(output, binary.LittleEndian, newHeaders)
	if err != nil {
		return fmt.Errorf("could not write headers: %v", err)
	}

	_, err = p.data.Seek(44, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek left: %v", err)
	}

	_, err = other.data.Seek(44, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek right: %v", err)
	}

	sampleSize := int(p.BitsPerSample / 8)
	bufLeft := make([]byte, sampleSize)
	bufRight := make([]byte, sampleSize)

	for {
		nL, errL := p.data.Read(bufLeft)
		if errL != nil && errL != io.EOF {
			return fmt.Errorf("error reading left: %v", errL)
		}

		nR, errR := other.data.Read(bufRight)
		if errR != nil && errR != io.EOF {
			return fmt.Errorf("error reading right: %v", errR)
		}

		// both data are consumed
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

// NewPCMAudio receives an io.ReadSeeker param and completes headers fields
func NewPCMAudio(data io.ReadSeeker) (*PCMAudio, error) {
	var p PCMAudio
	if err := binary.Read(data, binary.LittleEndian, &p.headers); err != nil {
		return nil, fmt.Errorf("could not read and parse the recieved data: %v", err.Error())
	}
	p.data = data
	return &p, nil
}
