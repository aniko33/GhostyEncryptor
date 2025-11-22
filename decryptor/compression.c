#include <stdint.h>
#include <stdlib.h>
#include <string.h>

uint8_t* RLEDecompress(const uint8_t* buf, size_t bufLen, size_t* outLen) {
    if (bufLen & 1) return NULL;

    size_t capacity = 0;
    for (size_t i = 0; i + 1 < bufLen; i += 2) {
        capacity += buf[i];
    }

    uint8_t* outBuf = malloc(capacity);
    if (!outBuf) return NULL;

    size_t pos = 0;
    for (size_t i = 0; i + 1 < bufLen; i += 2) {
        uint8_t count = buf[i];
        uint8_t val   = buf[i+1];
        memset(outBuf + pos, val, count);
        pos += count;
    }

    *outLen = pos;
    return outBuf;
}
