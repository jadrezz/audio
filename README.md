# WAV PCM Audio Library for Go

[![Go](https://img.shields.io/badge/Go-1.18%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A lightweight, dependency-free Go library for reading, validating, merging **WAV PCM audio files** and concatenate them. Designed for
reliability and ease of use in audio processing pipelines—ideal for phone call recordings, batch operations, or simple
audio editing tasks.

## ✨ Features

- Parse WAV file headers and access metadata (sample rate, bit depth, channels, etc.)
- Validate WAV file integrity and PCM format compliance
- Merge two mono WAV files into a single stereo file
- Concatenate multiple WAV files with matching formats
- Zero external dependencies — pure Go

## 📦 Installation

```bash
go get github.com/jadrezz/audio@latest
```

## 👨🏻‍💻 API Reference

* ### **NewPCMAudio️(io.ReadSeeker)**

Function that accepts io.ReadSeeker object, parses it and creates \*PCMAudio
object ready to use

* ### **PCMAudio**

PCMAudio is a structure that has methods to operate with WAVE data and
access to its fields such as audioFormat, byteRate, blockAlign etc...

* ### **\*PCMAudio.Validate()**

This method should be called before proceeding to Merge or Concat WAVE data.

Thus, we make sure both files are ready for operations

* ### **\*PCMAudio.Merge(\*PCMAudio, io.Writer)**

This method is used to merge 2 WAVE PCM audio files.
Since you have an instance of *PCMAudio object, you
need to create another one which is going to represent the file
for merging

* ### **\*PCMAudio.Concat(\*PCMAudio, io.Writer)**

This method is used to concatenate 2 WAVE audio files.
As well as PCMAudio.Merge, it requires other \*PCMAudio object
and the output for a new instance

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

* Concatenate two WAVE files

A good solution if you need concatenate 2 .wav files efficiently

```go
f1, err := os.Open("welcome.wav")
if err != nil {
    log.Fatal(err)
}
defer f1.Close()

welcomeSound, err := audio.NewPCMAudio(f1)
if err != nil {
    log.Fatal(err)
}
if ok, err := welcomeSound.Validate(); !ok {
    log.Fatal(err)
}

f2, err := os.Open("ring.wav")
if err != nil {
    log.Fatal(err)
}
defer f2.Close()

ring, err := audio.NewPCMAudio(f2)
if err != nil {
    log.Fatal(err)
}
if ok, err := ring.Validate(); !ok {
    log.Fatal(err)
}

output, err := os.Create("output.wav")
if err != nil {
    log.Fatal(err)
}
defer output.Close()

err = welcomeSound.Concat(ring, output)
if err != nil {
    log.Fatal(err)
}
```
