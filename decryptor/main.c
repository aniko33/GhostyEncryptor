#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "xor.h"
#include "compression.h"
#include "nibble.h"
#include "yenc.h"

int main(int argc, char** argv) {
    if (argc < 3) {
        printf("Usage: %s <shellcode.bin> <key>\n", argv[0]);
        return 0;
    }

    char* endptr;
    long value = strtol(argv[2], &endptr, 0);
    uint8_t key = (uint8_t) value;
    FILE* fd = fopen(argv[1], "rb");

    fseek(fd, 0L, SEEK_END);
    size_t fileSize = ftell(fd);
    fseek(fd, 0L, SEEK_SET);

    uint8_t* fileBuf = malloc(fileSize);
    fread(fileBuf, sizeof(uint8_t), fileSize, fd);
    fclose(fd);

    size_t outLen;
    uint8_t* shellcodeStep1 = YEncDecode(fileBuf, fileSize, &outLen);
    uint8_t* shellcodeStep2 = NibbleDecode(shellcodeStep1, outLen, &outLen);
    free(shellcodeStep1);
    uint8_t* shellcodeStep3 = RLEDecompress(shellcodeStep2, outLen, &outLen);
    free(shellcodeStep2);
    XorEncryptDecrypt(shellcodeStep3, key, outLen);

    FILE* outfd = fopen("test.out", "wb");
    fwrite(shellcodeStep3, 1, outLen, outfd);
    fclose(outfd);

    return 0;
}
