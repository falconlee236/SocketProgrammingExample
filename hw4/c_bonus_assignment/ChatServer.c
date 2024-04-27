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

#define BUFFER_SIZE 100 // Max Buffer size
#define PORT 30532 // Server port number
#define MAX_CLIENT 100 // Max client number
#define SEC(t) ((t).tv_sec + (t).tv_nsec / 1e+9) // second to millisecond

typedef struct s_clientInfo{
    char* ip; //client ip
    int port; //client port
    int num; // client id
}client_info;

// print current connect information
void print_connect_status(int client_num, int total_client_num, int is_connected);
// sigInt handler
void sigint_handler(int signum);

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
    printf("Server is ready to receive on port %d\n", PORT);

    fd_set reads, temps;
    FD_ZERO(&reads); // fd_set init
    FD_SET(serv_sock, &reads); // assign server fd to reads fd set
    int fd_max = serv_sock; // assign server socket to max, client request fd is always after this fd_max.
    int client_id = 0;

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
                    // connect that client
                    int clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_len);
                    char *client_ip = inet_ntoa(clnt_addr.sin_addr);
                    int client_port = ntohs(clnt_addr.sin_port);
                    printf("Connection request from %s:%d\n", client_ip, client_port);
                    // get next client number
                    client_arr[clnt_sock].num = ++client_id;
                    client_arr[clnt_sock].ip = client_ip;
                    client_arr[clnt_sock].port = client_port;
                    total_client_num++;
                    print_connect_status(client_arr[clnt_sock].num, total_client_num, 1);

                    FD_SET(clnt_sock, &reads); // set that client fd to 1
                    if(fd_max < clnt_sock){ // set max fd to that client fd
                        fd_max = clnt_sock;
                    }
                } else{ // already connected client
                    char type_str[BUFFER_SIZE];
                    int str_len = read(fd, type_str, BUFFER_SIZE);
                    if(str_len == 0){ // disconnect request
                        FD_CLR(fd, &reads); //change that fd to 0
                        close(fd);
                        total_client_num--;
                        print_connect_status(client_arr[fd].num, total_client_num, 0);
                    } else {
                        type_str[str_len] = '\0';
                        printf("Command %s", type_str);

                        char res[BUFFER_SIZE * 5] = {0, };
                        if (strncmp(type_str, "1\n", str_len) == 0){
                            char text[BUFFER_SIZE];
                            str_len = read(fd, text, BUFFER_SIZE);
                            for(int i = 0; i < str_len; i++){
                                if (islower(text[i])) res[i] = toupper(text[i]);
                                else res[i] = text[i];
                            }
                        } else if (strncmp(type_str, "2\n", str_len) == 0){
                            sprintf(res, "client IP = %s, port = %d\n", client_arr[fd].ip, client_arr[fd].port);
                        } else if (strncmp(type_str, "3\n", str_len) == 0){
                            sprintf(res, "request served %d=\n",1);
                        } else if (strncmp(type_str, "4\n", str_len) == 0){
                            struct timespec cur_time;
                            clock_gettime(CLOCK_MONOTONIC, &cur_time);
                            long duration = SEC(cur_time) - SEC(server_start);
                            int hours = duration / 3600;
                            int minutes = (hours % 3600) / 60;
                            int seconds = duration % 60;
                            sprintf(res, "run time = %02d:%02d:%02d\n", hours, minutes, seconds);
                        }
                        // send to client
                        write(fd, res, strlen(res));
                    }
                }
            }
        }
    }
}

void print_connect_status(int client_num, int total_client_num, int is_connected){
    time_t raw_time;
    struct tm *time_info;
    char time_str[9];

    // calculate current time
    time(&raw_time);
    time_info = localtime(&raw_time);

    strftime(time_str, sizeof(time_str), "%H:%M:%S", time_info);
    if (is_connected){
        printf("[Time: %s] Client %d connected. Number of clients connected = %d\n",
               time_str, client_num, total_client_num);
    } else {
        printf("[Time: %s] Client %d disconnected. Number of clients connected = %d\n",
               time_str, client_num, total_client_num);

    }
}

void sigint_handler(int signum){
    printf("\nBye bye~\n");
    // if server socket opened then close that socket
    if (serv_sock > 0) close(serv_sock);
    exit(signum);
}