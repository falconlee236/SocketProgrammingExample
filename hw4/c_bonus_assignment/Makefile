# Makefile
# 20190532 sang yun lee

CC = gcc

NAME_CLIENT = ChatClient
SRCS_CLIENT = ChatClient.c
OBJS_CLIENT = $(SRCS_CLIENT:.c=.o)

NAME_SERVER = ChatServer
SRCS_SERVER = ChatServer.c
OBJS_SERVER = $(SRCS_SERVER:.c=.o)

all: $(NAME_CLIENT) $(NAME_SERVER)

$(NAME_CLIENT): $(OBJS_CLIENT)
	$(CC) $^ -o $@ -lpthread

$(NAME_SERVER): $(OBJS_SERVER)
	$(CC) $^ -o $@ -lpthread

%.o: %.c
	$(CC) -c $< -o $@

clean:
	rm -rf $(OBJS_CLIENT) $(OBJS_SERVER)

fclean: clean
	rm -rf $(NAME_CLIENT) $(NAME_SERVER)

re: fclean all

.PHONY: all re clean fclean