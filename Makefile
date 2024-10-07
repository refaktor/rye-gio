# https://github.com/refaktor/rye-gioui

BIN_ROOT=$(PWD)/.bin
export PATH:=$(BIN_ROOT):$(PATH)

NAME=rye-gioui

all: gen bin

print:
	@echo ""

gen:
	cd gioui && go run .
bin:
	mkdir -p $(BIN_ROOT)
	go build -o $(BIN_ROOT)/$(NAME)

RUN_PATH=$(PWD)/examples
RUN_NAME=hello_gio.rye
RUN_NAME=click_counter.rye

run-h:
	$(NAME) -h
run:
	$(NAME) $(RUN_PATH)/$(RUN_NAME)