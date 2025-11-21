all: env decryptor.exe ghostyencryptor
env:
	mkdir -p ./tests/build
	
ghostyencryptor: ./src/*.go
	go build -o ghostyencryptor ./src

decryptor.exe: ./tests/build/main.o ./tests/build/yenc.o ./tests/build/xor.o ./tests/build/compression.o ./tests/build/nibble.o
	x86_64-w64-mingw32-gcc ./tests/build/*.o -o decryptor.exe

./tests/build/main.o: ./tests/main.c
	x86_64-w64-mingw32-gcc ./tests/main.c -c -o ./tests/build/main.o

./tests/build/yenc.o: ./tests/yenc.c
	x86_64-w64-mingw32-gcc ./tests/yenc.c -c -o ./tests/build/yenc.o

./tests/build/compression.o: ./tests/compression.c
	x86_64-w64-mingw32-gcc ./tests/compression.c -c -o ./tests/build/compression.o

./tests/build/xor.o: ./tests/xor.c
	x86_64-w64-mingw32-gcc ./tests/xor.c -c -o ./tests/build/xor.o

./tests/build/nibble.o: ./tests/nibble.c
	x86_64-w64-mingw32-gcc ./tests/nibble.c -c -o ./tests/build/nibble.o
