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

    // total buffer
    char buffer[BUFFER_SIZE];
    memset(&buffer, 0, sizeof (buffer));

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
    // receive nickname response
    if (read(client_socket, buffer, BUFFER_SIZE) == -1){
        perror("Failed to connect to server\n");
        exit(EXIT_FAILURE);
    }
    while (1){
        memset(buffer, 0, sizeof (buffer));
        fgets(buffer, sizeof (buffer), stdin);
        write(client_socket, buffer, BUFFER_SIZE);
        printf("you enter %s\n", buffer);
        memset(buffer, 0, sizeof (buffer));
        read(client_socket, buffer, BUFFER_SIZE);
        printf("you received %s\n", buffer);
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