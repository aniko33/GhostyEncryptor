package main

import "bytes"

func RLECompress(buf []byte) []byte {
    var outBuf bytes.Buffer
    prevByte := buf[0]
    count := byte(1)

    for _, b := range buf[1:] {
        if b == prevByte && count < 255 {
            count++
        } else {
            outBuf.WriteByte(count)
            outBuf.WriteByte(prevByte)
            prevByte = b
            count = 1
        }
    }

    outBuf.WriteByte(count)
    outBuf.WriteByte(prevByte)

    return outBuf.Bytes()
}

func RLEDecompress(buf []byte) []byte {
    var outBuf []byte

    for i := 0; i < len(buf) - 1; i+=2 {
        count := int(buf[i])
        b     := buf[i + 1]
        for range count {
            outBuf = append(outBuf, b)
        }
    }

    return outBuf
}
