package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

func checkArgValue(index_to_value int) {
    if index_to_value >= len(os.Args) {
        panic("Argument value not found")
    }
}

func entropyCalc(buf []byte) float64 {
    if len(buf) == 0 {
        return 0.0
    }

    var counts [256]int

    for _, b := range buf {
        counts[b]++
    }

    dataLen := float64(len(buf))
    entropy := 0.0

    for _, count := range counts {
        if count == 0 {
            continue
        }

        p       := float64(count) / dataLen
        entropy -= p * math.Log2(p)
    }

    return entropy
}

func main() {
    shellcodePath := ""
    keyHex        := ""
    outFile       := ""
    enableVerify  := false

    if len(os.Args) <= 1 {
        fmt.Printf("Usage: %s <shellcode.bin> [-key 0xFF] [-out output.bin] [-verify]\n", os.Args[0])
        os.Exit(-1)
    }

    shellcodePath = os.Args[1]

    if len(os.Args) > 2 {
        for i, arg := range os.Args {
            switch arg {
                case "-key":
                    checkArgValue(i + 1)
                    keyHex = os.Args[i + 1]
                case "-out":
                    checkArgValue(i + 1)
                    outFile = os.Args[i + 1]
                case "-verify":
                    enableVerify = true
            }
        }
    }

    var keyBuf byte

    if keyHex == "" {
        keyBuf = byte(rand.Intn(0xFF))
    } else {
        key, err := strconv.ParseUint(keyHex, 0, 8)

        if err != nil {
            panic(err)
        }

        keyBuf = byte(key)
    }

    if outFile == "" {
        outFile = "shellcode.out"
    }

    stat, err := os.Stat(shellcodePath)

    if err != nil {
        panic(err)
    }

    if stat.IsDir() {
        panic("Shellcode can't be a directory")
    }

    shellcodeBytes, err := os.ReadFile(shellcodePath)

    if err != nil {
        panic(err)
    }

    fmt.Printf("[i] KEY = 0x%X\n", keyBuf)
    fmt.Printf("[i] Size of shellcode %d\n", len(shellcodeBytes))
    fmt.Printf("[i] Original Entropy %f\n", entropyCalc(shellcodeBytes))

    fmt.Println("[+] Encryption...")
    shellcodeBytes = XorEncryptDecrypt(shellcodeBytes, keyBuf)

    fmt.Println("[+] Compression...")
    shellcodeBytes = RLECompress(shellcodeBytes)

    fmt.Println("[+] Nibble encoding...")
    shellcodeBytes = NibbleEncode(shellcodeBytes)

    fmt.Println("[+] YEnc Encoding...")
    shellcodeBytes = YEncEncode(shellcodeBytes)

    err = os.WriteFile(outFile, shellcodeBytes, 0644)

    if err != nil {
        panic(err)
    }

    if enableVerify {
        fmt.Println("[+] Verifying if the shellcode is valid...")

        shellcodeVerifyBytes, _ := YEncDecode(shellcodeBytes)
        shellcodeVerifyBytes = NibbleDecode(shellcodeVerifyBytes)
        shellcodeVerifyBytes = RLEDecompress(shellcodeVerifyBytes)
        shellcodeVerifyBytes = XorEncryptDecrypt(shellcodeVerifyBytes, keyBuf)
    }

    fmt.Printf("[+] Entropy %f\n", entropyCalc(shellcodeBytes))
    fmt.Printf("[i] Size of shellcode %d\n", len(shellcodeBytes))
    fmt.Printf("[i] Shellcode output %s\n", outFile)
}
