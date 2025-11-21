#include <stdint.h>
#include <stdlib.h>

uint8_t* NibbleDecode(const uint8_t* buf, size_t bufSize, size_t* outLen) {
    uint8_t* outBuf = malloc(bufSize / 2);

    for (int i = 0, j = 0; i < bufSize; i+=2) {
        uint8_t hi = buf[i] - 0x40;
        uint8_t lo = buf[i + 1] - 0x50;

        outBuf[j++] = ((hi & 0xF) << 4 | (lo & 0xF));
    }

    *outLen = bufSize / 2;
    return outBuf;
}
