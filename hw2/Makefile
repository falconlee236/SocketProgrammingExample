SRC_DIR := bouns_assignment
OUT_DIR := .
SRCS := $(wildcard $(SRC_DIR)/*.java)

CLS := $(SRCS:$(SRC_DIR)/%.java=$(OUT_DIR)/%.class)

# compiler and compiler flags
JC := javac
JCFLAGS := -d $(OUT_DIR)/ -cp $(SRC_DIR)/
# suffixes
.SUFFIXES: .java

# default target(s)
all: $(CLS)

$(CLS): $(OUT_DIR)/%.class: $(SRC_DIR)/%.java
	$(JC) $(JCFLAGS) $<

##
# clean up any output files
##
clean:
	rm -rf $(OUT_DIR)/*.class

re: clean all

# targets that do not produce output files
.PHONY: all clean