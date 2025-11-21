#include <stdint.h>

void XorEncryptDecrypt(uint8_t* buf, uint8_t key, size_t bufLen) {
    for (int i = 0; i < bufLen; i++) {
        buf[i] = buf[i] ^ key;
        key = (key + 1) & 0xFF;
    }
}
