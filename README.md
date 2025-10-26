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
f, err := os.Open("mono.wav")
if err != nil {
log.Fatal(err)
}
defer f.Close()

wav, err := audio.NewPCMAudio(f)
if err != nil {
log.Fatal(err)
}

if ok, err := wav.Validate(); !ok {
log.Fatal("Invalid WAV file:", err)
}

fmt.Printf("Sample rate: %d Hz\n", wav.SampleRate)
fmt.Printf("Bits per sample: %d\n", wav.BitsPerSample)
fmt.Printf("Channels: %d\n", wav.NumChannels)
fmt.Printf("Byte rate: %d B/s\n", wav.ByteRate)
```

* Merge two mono recordings into stereo

Perfect for combining separate caller/callee tracks from a phone call:
```go
f1, err := os.Open("client.wav")
if err != nil {
log.Fatal(err)
}
defer f1.Close()

client, err := audio.NewPCMAudio(f1)
if err != nil {
log.Fatal(err)
}
if ok, err := client.Validate(); !ok {
log.Fatal(err)
}

f2, err := os.Open("operator.wav")
if err != nil {
log.Fatal(err)
}
defer f2.Close()

operator, err := audio.NewPCMAudio(f2)
if err != nil {
log.Fatal(err)
}
if ok, err := operator.Validate(); !ok {
log.Fatal(err)
}

output, err := os.Create("output.wav")
if err != nil {
log.Fatal(err)
}
defer output.Close()

err = client.Merge(operator, output)
if err != nil {
log.Fatal(err)
}
```

