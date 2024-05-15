# Makefile
# 20190532 sang yun lee

SRC_DIR := .
OUT_DIR := .
SRCS := $(wildcard $(SRC_DIR)/*.java)

CLS := $(SRCS:$(SRC_DIR)/%.java=$(OUT_DIR)/%.class)

# compiler and compiler flags
JC := javac
# suffixes
.SUFFIXES: .java

# default target(s)
all: $(CLS)

$(CLS): $(OUT_DIR)/%.class: $(SRC_DIR)/%.java
	$(JC) $<

##
# clean up any output files
##
clean:
	rm -rf $(OUT_DIR)/*.class

re: clean all

# targets that do not produce output files
.PHONY: all clean