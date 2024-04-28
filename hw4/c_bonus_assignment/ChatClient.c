/*
 * ChatClient.c
 * 20190532 Sang yun Lee
 */

#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h> // inet_ntoa
#include <time.h> //localtime
#include <signal.h> //signal
#include <ctype.h> //islower, toupper
#include "map.h"

#define BUFFER_SIZE 1024
#define NICKNAME_SIZE 40
#define SERVER_IP "127.0.0.1" // server ip
#define SERVER_PORT 30532 // server port

// global client socket, using close fd when sigint occur
int client_socket;

// sigInt handler
void sigint_handler(int signum);
char **splitN(char *str, const char *delim, int n);

int main(int ac, char **av){
    // system args handling
    if (ac != 2) {
        printf("This program need to one argument\n");
        exit(EXIT_FAILURE);
    }
    // get nickname from system args
    char nicknameBuffer[NICKNAME_SIZE];
    memset(&nicknameBuffer, 0, sizeof (nicknameBuffer));
    strcpy(nicknameBuffer, av[1]);

    //signal
    signal(SIGINT, sigint_handler);

    // create commandMap
    Map *commandMap = createMap(5); // 초기 용량 5인 Map 생성
    // insert init command to Map
    insert(commandMap, "ls", 1);
    insert(commandMap, "secret", 2);
    insert(commandMap, "except", 3);
    insert(commandMap, "ping", 4);
    insert(commandMap, "quit", 5);

    // assign fd for client socket and IPv4, TCP
    // socket creation error handling
    if ((client_socket = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
        perror("Socket creation error");
        exit(EXIT_FAILURE);
    }

    // init server address to struct
    struct sockaddr_in server_addr;
    memset(&server_addr, 0, sizeof (server_addr));
    server_addr.sin_family = AF_INET; //  ipv4
    server_addr.sin_addr.s_addr = inet_addr(SERVER_IP); // server ip assign
    server_addr.sin_port = htons(SERVER_PORT); // big endian to little endian

    // connect to server
    if (connect(client_socket, (struct sockaddr*)&server_addr, sizeof(server_addr)) == -1){
        perror("connect");
        exit(EXIT_FAILURE);
    }

    // send nickname to server
    if (write(client_socket, nicknameBuffer, NICKNAME_SIZE) == -1){
        perror("Failed to connect to server\n");
        exit(EXIT_FAILURE);
    }

    char nicknameResBuffer[BUFFER_SIZE];
    // receive nickname response
    if (read(client_socket, nicknameResBuffer, BUFFER_SIZE) == -1){
        perror("Failed to connect to server\n");
        exit(EXIT_FAILURE);
    }
    // extract nickname status
    char** accessRes = splitN(nicknameResBuffer, "\n", 2);
    // print status msg
    printf("%s\n", accessRes[1]);
    // is status is not ok, exit client program
    if (strcmp(accessRes[0], "404") == 0){
        free(accessRes);
        exit(EXIT_FAILURE);
    }

    // select data structure init
    fd_set read_fds;
    FD_ZERO(&read_fds);
    // set stdin to fd_set
    FD_SET(STDIN_FILENO, &read_fds);
    // set client socket to fd_set
    FD_SET(client_socket, &read_fds);

    while (1){
        struct timeval timeout;
        timeout.tv_sec = 5;
        timeout.tv_usec = 0;

        // init tmp fd_set
        fd_set tmp_fds = read_fds;
        // select function modifies original fd_set array, find fd was changed
        if (select(client_socket + 1, &tmp_fds, 0, 0,&timeout) < 0){
            perror("select error");
            exit(EXIT_FAILURE);
        }

        // server data case
        if (FD_ISSET(client_socket, &tmp_fds)){
            char buffer[BUFFER_SIZE] = {0, };
            // get server response
            ssize_t bytes_received = read(client_socket, buffer, BUFFER_SIZE);
            // server disconnected case
            if (bytes_received <= 0){
                printf("\nServer disconnected\n");
                break;
            } else { // otherwise, print msg
                printf("%s", buffer);
            }
        }

        // user standard input case
        if (FD_ISSET(STDIN_FILENO, &tmp_fds)){
            char buffer[BUFFER_SIZE] = {0, };
            // input client msg, delimiter = "\n"
            fgets(buffer, sizeof (buffer), stdin);
            printf("\n");
            // check command, command case
            if (buffer[0] == '\\'){
                size_t newline_idx = strcspn(buffer, "\n");
                buffer[newline_idx] = '\0';
                // try to split 2 substring
                char** command_split = splitN(buffer + 1, " ", 2);
                int len = 0;
                // calculate array length
                while (command_split[len]){
                    len++;
                }
                // try to find that command is valid in map structure
                char command = find(commandMap, command_split[0]);
                if (command == -1 || // can't search command
                    ((command == 1 || command == 4 || command == 5) && len != 1) || // ls, ping, quit error
                    ((command == 2 || command == 3) && len != 2)){ // secret, except error
                    printf("invalid command\n");
                    continue;
                }
                // command quit
                if (command == 5){
                    free(command_split);
                    sigint_handler(SIGINT);
                }
                // command ls, command ping
                if (command == 1 || command == 4) {
                    char command_buffer[1];
                    command_buffer[0] = command;
                    write(client_socket, command_buffer, 1);
                } else {
                    char command_buffer[BUFFER_SIZE] = {0, };
                    command_buffer[0] = command;
                    command_buffer[1] = '\0';
                    strcat(command_buffer, command_split[1]);
                    command_buffer[1] = ' ';
                    write(client_socket, command_buffer, sizeof(command_buffer));
                }
            } else {
                // send to server
                write(client_socket, buffer, BUFFER_SIZE);
            }
        }
    }

    // memory free in map
    freeMap(commandMap);
    return 0;
}

void sigint_handler(int signum){
    printf("\ngg~\n");
    // if server socket opened then close that socket
    if (client_socket > 0) close(client_socket);
    exit(signum);
}

// clang version splitN function
char **splitN(char *str, const char *delim, int n) {
    // pointer array allocated
    char **tokens = (char **)malloc(sizeof(char *) * (n+1));
    if (tokens == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(1);
    }
    char *token;
    int count = 0;
    // get first token
    token = strtok(str, delim);
    // save token to array
    tokens[count++] = token;
    // loop string end, n times
    while (token != NULL && count < n) {
        // get next token
        token = strtok(NULL, delim);
        // save token to array
        tokens[count++] = token;
    }
    // add null pointer, that is end of array
    tokens[count] = NULL;
    // return token array
    return tokens;
}
