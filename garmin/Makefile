-include .env

SDK := $(shell ls -d "$(CIQ_SDK_DIR)"/connectiq-sdk-* 2>/dev/null | sort -V | tail -1)
MONKEYC := $(SDK)/bin/monkeyc
DEVICE := fr165
OUTPUT := bin/trainwatch.prg

.PHONY: build

build:
	"$(MONKEYC)" -d $(DEVICE) -f monkey.jungle -o $(OUTPUT) -y developer_key --release
