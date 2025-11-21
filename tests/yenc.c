#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <stdbool.h>

uint8_t* YEncDecode(const unsigned char* buf, size_t bufLen, size_t* outLen) {
    uint8_t* outBuf = malloc(bufLen);
    size_t cursor = 0;
    for (int i = 0; i < bufLen; i++) {
        uint8_t c = buf[i];

        if (c == '\r' || c == '\n') {
            continue;
        }

        if (c == '=') {
            i++;
            if (i >= bufLen) {
                free(outBuf);
                return NULL;
            }
            c = (c - 64) % 256;
        }
        c = (c - 42) % 256;
        outBuf[cursor++] = c;
    }

    *outLen = cursor;
    return outBuf;
}
