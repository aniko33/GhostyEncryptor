CC=x86_64-w64-mingw32-gcc

all: env decryptor.exe ghostyencryptor
env:
	mkdir -p ./decryptor/build
	
ghostyencryptor: ./src/*.go
	go build -o ghostyencryptor ./src

decryptor.exe: ./decryptor/build/main.o ./decryptor/build/yenc.o ./decryptor/build/xor.o ./decryptor/build/compression.o ./decryptor/build/nibble.o
	$(CC) ./decryptor/build/*.o -o decryptor.exe

./decryptor/build/main.o: ./decryptor/main.c
	$(CC) ./decryptor/main.c -c -o ./decryptor/build/main.o

./decryptor/build/yenc.o: ./decryptor/yenc.c
	$(CC) ./decryptor/yenc.c -c -o ./decryptor/build/yenc.o

./decryptor/build/compression.o: ./decryptor/compression.c
	$(CC) ./decryptor/compression.c -c -o ./decryptor/build/compression.o

./decryptor/build/xor.o: ./decryptor/xor.c
	$(CC) ./decryptor/xor.c -c -o ./decryptor/build/xor.o

./decryptor/build/nibble.o: ./decryptor/nibble.c
	$(CC) ./decryptor/nibble.c -c -o ./decryptor/build/nibble.o
