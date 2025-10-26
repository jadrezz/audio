package audio

import "errors"

// checkRequiredParams validates required WAVE properties
// to ensure both WAVEs are available for processing.
func checkRequiredParams(left, right *PCMAudio) error {
	if !left.valid || !right.valid {
		return errors.New("validate each audio file first")
	}

	if left.SampleRate != right.SampleRate {
		return errors.New("rate of both audio files must match")
	}

	if left.BitsPerSample != right.BitsPerSample {
		return errors.New("bits per sample of both audio files must match")
	}

	if left.NumChannels != right.NumChannels {
		return errors.New("number of channels of both audio files must match")
	}

	return nil
}

// chooseMode figures out which audio mode a new audio should have
func chooseMode(left, right *PCMAudio) uint16 {
	if left.NumChannels == mono && right.NumChannels == mono {
		return mono
	}
	return stereo
}
