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

# How do: Nibble Encoding

Nibble encoding is a method to encode the data in simple way with overhead of *1 byte to 2 bytes*

## Abstract

### Encode

1. Take a byte like: `11110000b` (`1` are first 4 bits and `0` are the last 4 bit)
2. split it in H (high) and L (low).
3. ```
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

> [!TODO]
> Add "How to:" Xor encrypt, YEnc encode, RLE compress
