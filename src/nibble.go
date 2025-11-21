package main

func NibbleEncode(buf []byte) []byte {
    var outBuf []byte

    for _, b := range buf {
        hi, lo := (b >> 4) & 0xF, b & 0xF
        outBuf = append(outBuf, byte(0x40 + hi))
        outBuf = append(outBuf, byte(0x50 + lo))
    }

    return outBuf
}

func NibbleDecode(buf []byte) []byte {
    var outBuf []byte

    for i := 0; i < len(buf) - 1; i += 2 {
        hi := buf[i] - 0x40
        lo := buf[i + 1] - 0x50

        outBuf = append(outBuf,  ((hi & 0xF) << 4) | (lo & 0xF))
    }

    return outBuf
}

