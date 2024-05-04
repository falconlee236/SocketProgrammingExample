# Makefile
# 20190532 sang yun lee

CC = gcc

C_NAME_CLIENT = C_ChatClient
C_SRCS_CLIENT = ChatClient.c map.c
C_OBJS_CLIENT = $(C_SRCS_CLIENT:.c=.o)

C_NAME_SERVER = C_ChatServer
C_SRCS_SERVER = ChatServer.c map.c
C_OBJS_SERVER = $(C_SRCS_SERVER:.c=.o)

all: $(C_NAME_CLIENT) $(C_NAME_SERVER)

$(C_NAME_CLIENT): $(C_OBJS_CLIENT)
	$(CC) $^ -o $@ -lpthread

$(C_NAME_SERVER): $(C_OBJS_SERVER)
	$(CC) $^ -o $@ -lpthread

%.o: %.c
	$(CC) -c $< -o $@

clean:
	rm -rf $(C_OBJS_CLIENT) $(C_OBJS_SERVER)

fclean: clean
	rm -rf $(C_NAME_CLIENT) $(C_NAME_SERVER)

re: fclean all

.PHONY: all re clean fclean
