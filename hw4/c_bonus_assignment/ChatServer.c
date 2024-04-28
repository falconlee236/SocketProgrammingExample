/*
 * ChatServer.c
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

#define BUFFER_SIZE 1024 // Max Buffer size
#define NICkNAME_SIZE 40 // nickname size
#define PORT 30532 // Server port number
#define MAX_CLIENT 8 // Max client number
#define SEC(t) ((t).tv_sec + (t).tv_nsec / 1e+9) // second to millisecond

typedef struct s_clientInfo{
    char* nickname; // client nickname
    char* ip; //client ip
    int port; //client port
}client_info;

// print current connect information
void print_connect_status(int client_num, int total_client_num, int is_connected);
// sigInt handler
void sigint_handler(int signum);
char* to_lower(char *str);

// global server socket, using close fd when sigint occur
int serv_sock;

int main(void){
    //server start time
    struct timespec server_start;
    clock_gettime(CLOCK_MONOTONIC, &server_start);

    //signal
    signal(SIGINT, sigint_handler);

    //total client number variable
    int total_client_num = 0;

    // client info array init
    client_info client_arr[MAX_CLIENT + 3];

    // assign fd for server socket, and IPv4, TCP
    serv_sock = socket(AF_INET, SOCK_STREAM, 0);
    // allow server to use same port
    int opt = 1;
    setsockopt(serv_sock, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    struct sockaddr_in serv_addr; //server address init
    bzero(&serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET; //internet protocol
    serv_addr.sin_addr.s_addr = htonl(INADDR_ANY); // multiple ip handling
    serv_addr.sin_port = htons(PORT); // big endian to little endian

    // bind socket info to server fd
    int binded;
    binded = bind(serv_sock, (struct sockaddr*)&serv_addr, sizeof(serv_addr));
    if(binded != 0){
        perror("bind error");
        exit(1);
    }

    // waiting client request
    binded = listen(serv_sock, 5);
    if(binded != 0){
        perror("listen error");
        exit(1);
    }
    struct sockaddr_in local_addr;
    socklen_t addr_len = sizeof(local_addr);
    if (getsockname(serv_sock, (struct sockaddr*)&local_addr, &addr_len) == -1) {
        perror("getsockname error");
        exit(1);
    }
    char* server_ip = inet_ntoa(local_addr.sin_addr);
    printf("Server is ready to receive on port %d\n", PORT);

    fd_set reads, temps;
    FD_ZERO(&reads); // fd_set init
    FD_SET(serv_sock, &reads); // assign server fd to reads fd set
    int fd_max = serv_sock; // assign server socket to max, client request fd is always after this fd_max.

    // create clientMap
    Map* client_map = createMap(MAX_CLIENT);

    while(1){
        struct timeval timeout;
        timeout.tv_sec = 5;
        timeout.tv_usec = 0;

        temps = reads; // select change original fd_set value
        if(select(fd_max + 1, &temps, 0, 0, &timeout) < 0){ // check array changed
            perror("select error");
            exit(1);
        }
        for(int fd = 0; fd < fd_max + 1; fd++){ // check which fd changed
            if(FD_ISSET(fd, &temps)){ // if changed
                if(fd == serv_sock){ // client request to server connecting case
                    struct sockaddr_in clnt_addr;
                    socklen_t clnt_len = sizeof(clnt_addr);
                    // accept client's connection
                    int clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_len);

                    // check nickname
                    char nickname_buffer[NICkNAME_SIZE] = {0, }, nickname_res_buffer[BUFFER_SIZE] = {0, };
                    int status_code = 200;
                    read(clnt_sock, nickname_buffer, NICkNAME_SIZE); // name\0, len = 4
                    if (total_client_num == 8){
                        status_code = 404;
                        sprintf(nickname_res_buffer, "%d\n[chatting room full. cannot connect.]\n", status_code);
                    } else if (find(client_map, nickname_buffer) != -1){
                        status_code = 404;
                        sprintf(nickname_res_buffer, "%d\n[nickname already used by another user. cannot connect.]\n", status_code);
                    } else {
                        insert(client_map, nickname_buffer, clnt_sock);
                        char *client_ip = inet_ntoa(clnt_addr.sin_addr);
                        int client_port = ntohs(clnt_addr.sin_port);
                        // get next client number
                        client_arr[clnt_sock].ip = strdup(client_ip);
                        client_arr[clnt_sock].port = client_port;
                        client_arr[clnt_sock].nickname = strdup(nickname_buffer);
                        total_client_num++;
                        FD_SET(clnt_sock, &reads); // set that client fd to 1
                        if(fd_max < clnt_sock){ // set max fd to that client fd
                            fd_max = clnt_sock;
                        }
                        sprintf(nickname_res_buffer, "%d\n[welcome %s to CAU net-class chat room at %s:%d.]\n[There are %d users in the room.]\n",
                                                     status_code, nickname_buffer, server_ip, PORT, total_client_num);
                        printf("[%s has joined from %s:%d.]\n"
                               "[There are %d users in room.]\n\n", nickname_buffer, client_ip, client_port, total_client_num);
                    }
                    write(clnt_sock, nickname_res_buffer, strlen(nickname_res_buffer));
                } else{ // already connected client
                    char buffer[BUFFER_SIZE] = {0, };
                    memset(&buffer, 0, sizeof (buffer));
                    ssize_t str_len = read(fd, buffer, BUFFER_SIZE);
                    if(str_len == 0){ // disconnect request
                        char sendMsg[BUFFER_SIZE] = {0, };
                        FD_CLR(fd, &reads); //change that fd to 0
                        close(fd);
                        total_client_num--;
                        delete(client_map, client_arr[fd].nickname);
                        sprintf(sendMsg, "[%s left the room.]\n"
                                         "[There are %d users in the chat room.]\n",
                                         client_arr[fd].nickname, total_client_num);
                        printf("%s\n", sendMsg);
                        for(int i = 0; i < client_map->size; i++){
                            int otherFd = client_map->data[i].value;
                            if (otherFd == fd)
                                continue;
                            write(otherFd, sendMsg, sizeof(sendMsg));
                        }
                    } else {
                        char sendMsg[BUFFER_SIZE] = {0, };
                        if (buffer[0] > 0 && buffer[0] < 6){
                            int command_type = buffer[0];
                            if (command_type == 1){
                                for(int i = 0; i < client_map->size; i++){
                                    int other_fd = (unsigned char)client_map->data[i].value;
                                    char client_info[BUFFER_SIZE] = {0, };
                                    sprintf(client_info, "<%s, %s, %d>\n",
                                            client_arr[other_fd].nickname, client_arr[other_fd].ip, client_arr[other_fd].port);
                                    strcat(sendMsg, client_info);
                                }
                            }
                            write(fd, sendMsg, sizeof(sendMsg));
                            continue;
                        }
                        sprintf(sendMsg, "%s> %s\n", client_arr[fd].nickname, buffer);
                        for(int i = 0; i < client_map->size; i++){
                            int otherFd = (unsigned char)client_map->data[i].value;
                            if (otherFd == fd)
                                continue;
                            write(otherFd, sendMsg, sizeof(sendMsg));
                        }
                        if(strstr(to_lower(buffer), "i hate professor") != NULL){
                            total_client_num--;
                            sprintf(sendMsg, "[%s is disconnected.]\n"
                                             "[There are %d users in the chat room.]\n", client_arr[fd].nickname, total_client_num);
                            printf("%s\n", sendMsg);
                            for(int i = 0; i < client_map->size; i++){
                                int otherFd = (unsigned char)client_map->data[i].value;
                                write(otherFd, sendMsg, sizeof(sendMsg));
                            }
                            FD_CLR(fd, &reads); //change that fd to 0
                            close(fd);
                            delete(client_map, client_arr[fd].nickname);
                        }
//                        type_str[str_len] = '\0';
//                        printf("Command %s", type_str);
//
//                        char res[BUFFER_SIZE * 5] = {0, };
//                        if (strncmp(type_str, "1\n", str_len) == 0){
//                            char text[BUFFER_SIZE];
//                            str_len = read(fd, text, BUFFER_SIZE);
//                            for(int i = 0; i < str_len; i++){
//                                if (islower(text[i])) res[i] = toupper(text[i]);
//                                else res[i] = text[i];
//                            }
//                        } else if (strncmp(type_str, "2\n", str_len) == 0){
//                            sprintf(res, "client IP = %s, port = %d\n", client_arr[fd].ip, client_arr[fd].port);
//                        } else if (strncmp(type_str, "3\n", str_len) == 0){
//                            sprintf(res, "request served %d=\n",1);
//                        } else if (strncmp(type_str, "4\n", str_len) == 0){
//                            struct timespec cur_time;
//                            clock_gettime(CLOCK_MONOTONIC, &cur_time);
//                            long duration = SEC(cur_time) - SEC(server_start);
//                            int hours = duration / 3600;
//                            int minutes = (hours % 3600) / 60;
//                            int seconds = duration % 60;
//                            sprintf(res, "run time = %02d:%02d:%02d\n", hours, minutes, seconds);
//                        }
//                        // send to client
//                        write(fd, res, strlen(res));
                    }
                }
            }
        }
    }
}

void sigint_handler(int signum){
    printf("\ngg~\n");
    // if server socket opened then close that socket
    if (serv_sock > 0) close(serv_sock);
    exit(signum);
}

// convert string to lowercase
char* to_lower(char *str){
    int i = 0;
    char* res = (char*)malloc(sizeof (char) * (strlen(str) + 1));
    while (str[i]){
        res[i] = (char)tolower(str[i]);
        i++;
    }
    res[i] = 0;
    return res;
}