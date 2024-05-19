# Makefile
# 20190532 sang yun lee

RS_NAME_CLIENT = rs_split_client
RS_SRCS_CLIENT = src/SplitFileClient.rs

RS_NAME_SERVER = rs_split_server
RS_SRCS_SERVER = src/SplitFileServer.rs

all: $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

$(RS_NAME_CLIENT):
	cargo build --release --bin $(RS_NAME_CLIENT)
	mv ./target/release/$(RS_NAME_CLIENT) ./src

$(RS_NAME_SERVER):
	cargo build --release --bin $(RS_NAME_SERVER)
	mv ./target/release/$(RS_NAME_SERVER) ./src

clean:
	cargo clean
	rm -rf ./src/$(RS_NAME_CLIENT) ./src/$(RS_NAME_SERVER)

fclean: clean
	rm -rf target

re: fclean all

.PHONY: all re clean fclean