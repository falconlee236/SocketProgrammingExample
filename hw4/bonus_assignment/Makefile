# Makefile
# 20190532 sang yun lee

CC = gcc
RC = rustc

C_NAME_CLIENT = C_ChatClient
C_SRCS_CLIENT = ChatClient.c map.c
C_OBJS_CLIENT = $(C_SRCS_CLIENT:.c=.o)

C_NAME_SERVER = C_ChatServer
C_SRCS_SERVER = ChatServer.c map.c
C_OBJS_SERVER = $(C_SRCS_SERVER:.c=.o)

RS_NAME_CLIENT = RS_ChatClient
RS_SRCS_CLIENT = ChatClient.rs

RS_NAME_SERVER = RS_ChatServer
RS_SRCS_SERVER = ChatServer.rs

all: $(C_NAME_CLIENT) $(C_NAME_SERVER) $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

$(C_NAME_CLIENT): $(C_OBJS_CLIENT)
	$(CC) $^ -o $@ -lpthread

$(C_NAME_SERVER): $(C_OBJS_SERVER)
	$(CC) $^ -o $@ -lpthread

$(RS_NAME_CLIENT): $(RS_SRCS_CLIENT)
	$(RC) $^ -o $@

$(RS_NAME_SERVER): $(RS_SRCS_SERVER)
	$(RC) $^ -o $@

%.o: %.c
	$(CC) -c $< -o $@

clean:
	rm -rf $(C_OBJS_CLIENT) $(C_OBJS_SERVER) $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

fclean: clean
	rm -rf $(C_NAME_CLIENT) $(C_NAME_SERVER) $(RS_NAME_CLIENT) $(RS_NAME_SERVER)

re: fclean all

.PHONY: all re clean fclean