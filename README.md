# WAV PCM Audio Library for Go

[![Go](https://img.shields.io/badge/Go-1.18%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A lightweight, dependency-free Go library for reading, validating, and merging **WAV PCM audio files**. Designed for reliability and ease of use in audio processing pipelines—ideal for phone call recordings, batch operations, or simple audio editing tasks.

## ✨ Features

- Parse WAV file headers and access metadata (sample rate, bit depth, channels, etc.)
- Validate WAV file integrity and PCM format compliance
- Merge two mono WAV files into a single stereo file
- Concatenate multiple WAV files with matching formats
- Zero external dependencies — pure Go

## 📦 Installation

```bash
go get github.com/jadrezz/audio
```

## 🚀 Usage examples
* Read and inspect WAV metadata
```go
f, err := os.Open("recording.wav")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

meta, err := audio.NewPCMAudioMetadata(f)
if err != nil {
    log.Fatal(err)
}

if ok, err := meta.Validate(); !ok {
    log.Fatal("Invalid WAV file:", err)
}

fmt.Printf("Sample rate: %d Hz\n", meta.SampleRate)
fmt.Printf("Bits per sample: %d\n", meta.BitsPerSample)
fmt.Printf("Channels: %d\n", meta.NumChannels)
fmt.Printf("Byte rate: %d B/s\n", meta.ByteRate)
```

* Merge two mono recordings into stereo

Perfect for combining separate caller/callee tracks from a phone call:
```go
f1, err := os.Open("part1.wav")
if err != nil {
	log.Fatal(err)
}
defer f1.Close()

meta, err := audio.NewPCMAudioMetadata(f1)
if err != nil {
	log.Fatal(err)
}
if ok, err := meta.Validate(); !ok {
	log.Fatal(err)
}
	
f2, err := os.Open("part2.wav")
if err != nil {
	log.Fatal(err)
}
defer f2.Close()

meta2, err := audio.NewPCMAudioMetadata(f2)
if err != nil {
	log.Fatal(err)
}
if ok, err := meta2.Validate(); !ok {
	log.Fatal(err)
}
output, err := os.Create("output.wav")
if err != nil {
	log.Fatal(err)
}
defer output.Close()

err = meta.Merge(meta2, output, f1, f2))
```

