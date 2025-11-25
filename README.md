# GhostyEncryptor

```
     .-----.
   .' -   - '.
  /  .-. .-.  \
  |  | | | |  |
   \ \o/ \o/ /
  _/    ^    \_
 | \  '---'  / |   ▗▄▄▖▗▖ ▗▖ ▗▄▖  ▗▄▄▖▗▄▄▄▖▗▖  ▗▖
 / /`--. .--`\ \  ▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌     █   ▝▚▞▘
/ /'---` `---'\ \ ▐▌▝▜▌▐▛▀▜▌▐▌ ▐▌ ▝▀▚▖  █    ▐▌
'.__.       .__.' ▝▚▄▞▘▐▌ ▐▌▝▚▄▞▘▗▄▄▞▘  █    ▐▌
    `|     |`
     |     \
     \      '--.
      '.        `\
        `'---.   |
           ,__) /
            `..'
```

# Build

```
$ git clone https://github.com/aniko33/GhostyEncryptor
$ cd src
$ make ghostyencryptor
```
Now the executable `ghostyencryptor` is the project root

# Usage

```
Usage: ./ghostyencryptor <shellcode.bin> [-key 0xFF] [-out output.bin] [-verify]
```

## Encrypt shellcode with random key

```
$ ./ghostyencryptor shellcode.bin
[i] KEY = 0xBC <--- GENERATED KEY!
[i] Size of shellcode 201798
[i] Original Entropy 6.046684
[+] Encryption...
[+] Compression...
[+] Nibble encoding...
[+] YEnc Encoding...
[+] Entropy 3.848014
[i] Size of shellcode 791272
[i] Shellcode output shellcode.out
```

## Encrypt shellcode with custom key

```
$ ./ghostyencryptor shellcode.bin -key 0xAA
[i] KEY = 0xAA
[i] Size of shellcode 201798
[i] Original Entropy 6.046684
[+] Encryption...
[+] Compression...
[+] Nibble encoding...
[+] YEnc Encoding...
[+] Entropy 3.843809
[i] Size of shellcode 792132
[i] Shellcode output shellcode.out
```

## Encrypt shellcode with custom key, verify option and save output as test.bin

```
$ ./ghostyencryptor shellcode.bin -key 0xFF -verify -out test.bin
[i] KEY = 0xFF
[i] Size of shellcode 201798
[i] Original Entropy 6.046684
[+] Encryption...
[+] Compression...
[+] Nibble encoding...
[+] YEnc Encoding...
[+] Verifying if the shellcode is valid...
[+] Entropy 3.841705
[i] Size of shellcode 798764
[i] Shellcode output test.bin
```

# How to decrypt it?

Into `decryptor/` there are a fully function example of decryptor written in C and can be ported in every programming language.

## How test it (default compiler `x86_64-w64-mingw32-gcc`)

```
$ mkdir -p decryptor/build
$ make decryptor.exe
```

> [!NOTE]
> To change compiler: set environment key `CC`
> ```
> $ CC=customCompiler make decryptor.exe
> ```

```
$ decryptor.exe <shellcode.bin> <key>
```

The decryptor will save the shellcode output to `test.out`

# Build your decryptor

The decryptor needs to perform the following steps:

- [YEnc](https://en.wikipedia.org/wiki/YEnc) Decoding
- [Nibble](#Nibble-Encoding) Decoding
- [RLE](https://en.wikipedia.org/wiki/Run-length_encoding) Decompression
- [Xor](https://en.wikipedia.org/wiki/XOR_cipher) Decryption

## Pseudocode (Python like)

```py
shellcodeStep1 = YEncDecode(shellcodeEncrypted)
shellcodeStep2 = NibbleDecode(shellcodeStep1)
shellcodeStep3 = RLEDecompress(shellcodeStep2)
decrypted      = XorEncryptDecrypt(shellcodeStep3, key)
```
# How do: yEnc Encoding

yEnc is a simple binary-to-text encoding used for Usenet/text-only transports. It maps raw bytes to mostly printable characters with low overhead and a small escape mechanism.

## Abstract

### Encode

1. Take a byte, e.g. `0xF0`.
2. Compute: `enc = (b + 42) % 256`.
3. If `enc` is one of the reserved bytes `{ 0x00, 0x0A, 0x0D, 0x3D }` then:
   - Output the escape byte `'='` with `(enc + 64) % 256`.
4. Otherwise output `enc` as a single byte

### Decode

1. Read bytes left-to-right.
2. If byte is `'='` (`0x3D`):
   - Read next byte
   - `enc = (t - 64 + 256) % 256`
   - If no next byte throw an error: truncated escape sequence.
3. Else: recover original so `b = (byte - 42 + 256) % 256`
4. Output the decoded byte

### Escape rules

- Escape byte: `=` (0x3D)
- Bytes that must be escaped (after adding 42): `0x00`, `0x0A` (LF), `0x0D` (CR), `0x3D` (`=`)
- Escaped transmitted value: `(enc + 64) % 256`
- Decoder reverses by subtracting 64 (mod 256)

# How do: Nibble Encoding

Nibble encoding is a method to encode the data in simple way with overhead of *1 byte to 2 bytes*

## Abstract

### Encode

1. Take a byte like: `11110000b` (`1` are first 4 bits and `0` are the last 4 bit)
2. split it in H (high) and L (low).
    -  ```
        H = 1111
        L = 0000
        ```
4. Increase H by `0x40` so `1001111b`
5. Increase L by `0x50` so `1010000b`
6. Return the value as 2 bytes separated: `[H, L]` so `[1001111b, 1010000b]`

### Decode

1. Take bytes like: `[1001111b, 1010000b]` - The first one is the ***high***, and the second is the ***low***
4. Decrease H by `0x40` so `1111b`
5. Decrease L by `0x50` so `0000b`
6. Return the recomposed byte: `((high & 0xF) << 4) | (low & 0xF)` so `11110000b`

# How do: RLE (Run-Length Encoding)

Simple RLE encodes runs of repeated bytes as pairs (count, byte). Overhead is 1 extra byte per run (count) — encoded stream is: [count, value, count, value, ...]. Counts are single bytes (0–255), so maximum run length is 255.

## Abstract

### Encode

1. Start with an input byte buffer.
2. Track the current run value (`prev`) and a run length `count` (byte).
3. Iterate input bytes:
   - If the next byte equals `prev` and `count < 255`, increment `count`.
   - Otherwise, write the run as two bytes: `[count, prev]`, then set `prev` to the new byte and `count = 1`.
4. After the loop, write the final run `[count, prev]`.
5. Return the concatenated output bytes.

> [!NOTE]
> - The first input byte initializes `prev` and `count = 1`.
> - A run of a single byte becomes `[1, value]`.
> - Large runs are split into multiple runs when longer than 255.

### Decode

1. Read the compressed buffer two bytes at a time: first byte is `count`, second is `value`.
2. For each pair, append `count` copies of `value` to the output.
3. Stop when you've consumed all pairs. If the compressed buffer length is odd, it's malformed.

## Stupid Example

Input: [0xAA, 0xAA, 0xAA, 0xBB]
Encoded: [3, 0xAA, 1, 0xBB]  
Decoded: [0xAA, 0xAA, 0xAA, 0xBB]

# How do: XOR stream (byte-wise) Encryption/Decryption

XOR stream here means each byte is XORed with a single-byte key that increments after every byte. The operation is symmetric: the same function encrypts and decrypts.

## Abstract

### Encode / Encrypt

1. Take input bytes and an initial single-byte key (0–255).
2. For each input byte `b`:
   - Compute out = `b XOR key` so `b ^ key`.
   - Append out to output.
   - Increment key: `key = (key + 1) & 0xFF` (wraps at 255→0).
3. Return the output bytes.

## Stupid Example

Input: [0x10, 0x20, 0x30], key = `0x01`
Keystream: [0x01, 0x02, 0x03]
Encrypted: [0x11, 0x22, 0x33]

# Benchmark

Go to [`BENCHMARK.md`](./BENCHMARK.md)

# Important notes

> [!IMPORTANT]
> I'm very lazy when it comes to writing good documentation, so "How to: *X*" may be documented incorrectly because I used AI (but i review it first).
