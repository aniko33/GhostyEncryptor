package main

func XorEncryptDecrypt(buf []byte, key byte) []byte {
    var outBuf []byte

    for _, b := range buf {
        outBuf = append(outBuf, byte(b) ^ key)
        key = (key + 1) & 0xFF
    }

    return outBuf
}
