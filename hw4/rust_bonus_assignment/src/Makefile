# Makefile
# 20190532 sang yun lee

RC = rustc

RS_NAME_CLIENT = RS_ChatClient
RS_SRCS_CLIENT = ChatClient.rs

RS_NAME_SERVER = RS_ChatServer
RS_SRCS_SERVER = ChatServer.rs

all: $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

$(RS_NAME_CLIENT): $(RS_SRCS_CLIENT)
	$(RC) $^ -o $@

$(RS_NAME_SERVER): $(RS_SRCS_SERVER)
	$(RC) $^ -o $@

%.o: %.c
	$(CC) -c $< -o $@

clean:
	rm -rf $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

fclean: clean
	rm -rf $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

re: fclean all

.PHONY: all re clean fclean