package main 

import "errors"

func YEncEncode(buf []byte) []byte {
    out := make([]byte, 0, len(buf))
    escape := map[byte]struct{}{
        0x00: {}, 0x0A: {}, 0x0D: {}, 0x3D: {},
    }

    for _, b := range buf {
        enc := (int(b) + 42) % 256
        eb := byte(enc)

        if _, needsEscape := escape[eb]; needsEscape {
            out = append(out, '=')
            eb = byte((enc + 64) % 256)
        }
        out = append(out, eb)
    }
    return out
}

func YEncDecode(buf []byte) ([]byte, error) {
    out := make([]byte, 0, len(buf))
    i := 0
    for i < len(buf) {
        b := buf[i]
        if b == '=' {
            i++
            if i >= len(buf) {
                return nil, errors.New("truncated escape sequence")
            }
            b = byte((int(buf[i]) - 64 + 256) % 256)
        }
        dec := byte((int(b) - 42 + 256) % 256)
        out = append(out, dec)
        i++
    }
    return out, nil
}
