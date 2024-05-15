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
#define PORT 20532 // Server port number
#define MAX_CLIENT 8 // Max client number

typedef struct s_clientInfo{
    char* nickname; // client nickname
    char* ip; //client ip
    int port; //client port
}client_info;

// sigInt handler
void sigint_handler(int signum);
// convert string to lowercase
char* to_lower(char *str);
// clang version splitN function
char **splitN(char *str, const char *delim, int n);

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
    client_info client_arr[MAX_CLIENT + 10];

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
                    	write(clnt_sock, nickname_res_buffer, strlen(nickname_res_buffer));
                        break;
                    } else if (find(client_map, nickname_buffer) != -1){
                        status_code = 404;
                        sprintf(nickname_res_buffer, "%d\n[nickname already used by another user. cannot connect.]\n", status_code);
                    	write(clnt_sock, nickname_res_buffer, strlen(nickname_res_buffer));
                        break;
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
						 struct sockaddr_in local_addr;
    					socklen_t addr_len = sizeof(local_addr);
    					if (getsockname(clnt_sock, (struct sockaddr*)&local_addr, &addr_len) == -1) {
        					perror("getsockname error");
        					exit(1);
    					}
    					char* server_ip = inet_ntoa(local_addr.sin_addr);

                        sprintf(nickname_res_buffer, "%d\n[welcome %s to CAU net-class chat room at %s:%d.]\n[There are %d users in the room.]\n",
                                                     status_code, nickname_buffer, server_ip, PORT, total_client_num);
                        printf("[%s has joined from %s:%d.]\n"
                               "[There are %d users in room.]\n\n", nickname_buffer, client_arr[clnt_sock].ip, client_port, total_client_num);
                    }
                    write(clnt_sock, nickname_res_buffer, strlen(nickname_res_buffer));
                } else{ // already connected client
                    char buffer[BUFFER_SIZE] = {0, };
                    memset(&buffer, 0, sizeof (buffer));
                    ssize_t str_len = read(fd, buffer, BUFFER_SIZE);
		    if(str_len == 1 && buffer[0] == '\n')
			    continue;
                    if(str_len == 0){ // disconnect request
                        char sendMsg[BUFFER_SIZE] = {0, };
                        FD_CLR(fd, &reads); //change that fd to 0
                        close(fd); //disconnect client
                        total_client_num--; //decrease total client num
                        // remove client information
                        delete(client_map, client_arr[fd].nickname);
                        sprintf(sendMsg, "[%s left the room.]\n"
                                         "[There are %d users in the chat room.]\n",
                                         client_arr[fd].nickname, total_client_num);
                        printf("%s\n", sendMsg);
                        // send to msg without itself
                        for(int i = 0; i < client_map->size; i++){
                            int otherFd = client_map->data[i].value;
                            if (otherFd == fd)
                                continue;
                            write(otherFd, sendMsg, sizeof(sendMsg));
                        }
                    } else {
                        char sendMsg[BUFFER_SIZE] = {0, };
                        //command case
                        if (buffer[0] > 0 && buffer[0] < 6){
                            int command_type = (unsigned char)buffer[0];
                            if (command_type == 1){
                                for(int i = 0; i < client_map->size; i++){
                                    int other_fd = (unsigned char)client_map->data[i].value;
                                    char client_info[BUFFER_SIZE] = {0, };
                                    sprintf(client_info, "<%s, %s, %d>\n",
                                            client_arr[other_fd].nickname, client_arr[other_fd].ip, client_arr[other_fd].port);
                                    strcat(sendMsg, client_info);
                                }
                            } else if (command_type == 2 || command_type == 3) {
                                // try to split 3 substring
                                char** command_split = splitN(buffer, " ", 3);
                                int len = 0;
                                // calculate array length
                                while (command_split[len])
                                    len++;
                                // find dst fd in map
                                char dst_fd = find(client_map, command_split[1]);
                                // secret, except error
                                if (len != 3 || dst_fd == -1){
                                    char* type_str = command_type == 2 ? "\\secret" : "\\except";
                                    printf("invalid command: %s %s %s\n", type_str, command_split[1], command_split[2]);
                                    if (dst_fd == -1){
                                        sprintf(sendMsg, "%s doesn't exist in the chat room.\n", command_split[1]);
                                        write(fd, sendMsg, sizeof(sendMsg));
                                    }
                                    continue;
                                } else {
                                    // prepare special message
                                    sprintf(sendMsg, "from: %s> %s\n", client_arr[fd].nickname, command_split[2]);
                                    free(command_split);
                                    // secret case
                                    if (command_type == 2){
                                        write(dst_fd, sendMsg, sizeof(sendMsg));
                                    } else { // except case
                                        for(int i = 0; i < client_map->size; i++){
                                            int otherFd = (unsigned char)client_map->data[i].value;
                                            if (otherFd == dst_fd)
                                                continue;
                                            write(otherFd, sendMsg, sizeof(sendMsg));
                                        }
                                    }
                                    char filter_buffer[BUFFER_SIZE] = {0, };
                                    if(strstr(to_lower(sendMsg), "i hate professor") != NULL){
                                        total_client_num--;
                                        // client shutdown case
                                        sprintf(filter_buffer, "[%s is disconnected.]\n"
                                                               "[There are %d users in the chat room.]\n", client_arr[fd].nickname, total_client_num);
                                        printf("%s\n", filter_buffer);
                                        // send msg to all client including itself
                                        for(int i = 0; i < client_map->size; i++){
                                            int otherFd = (unsigned char)client_map->data[i].value;
                                            write(otherFd, filter_buffer, sizeof(filter_buffer));
                                        }
                                        FD_CLR(fd, &reads); //change that fd to 0
                                        close(fd); // disconnect client
                                        // remove client info
                                        delete(client_map, client_arr[fd].nickname);
                                    }
                                    continue;
                                }
                            } else if (command_type == 4){
                                sendMsg[0] = (char)command_type;
                            } else if (command_type == 5){
                                char sendMsg[BUFFER_SIZE] = { 0,};
                                FD_CLR(fd, &reads); // change that fd to 0
                                close(fd);          // disconnect client
                                total_client_num--; // decrease total client num
                                // remove client information
                                delete (client_map, client_arr[fd].nickname);
                                sprintf(sendMsg, "[%s left the room.]\n"
                                                 "[There are %d users in the chat room.]\n",
                                        client_arr[fd].nickname, total_client_num);
                                printf("%s\n", sendMsg);
                                // send to msg without itself
                                for (int i = 0; i < client_map->size; i++){
                                    int otherFd = client_map->data[i].value;
                                    if (otherFd == fd)
                                        continue;
                                    write(otherFd, sendMsg, sizeof(sendMsg));
                                }
                                continue;
                            }
                            write(fd, sendMsg, sizeof(sendMsg));
                            continue;
                        }
                        // not command case
                        sprintf(sendMsg, "%s> %s\n", client_arr[fd].nickname, buffer);
                        for(int i = 0; i < client_map->size; i++){
                            int otherFd = (unsigned char)client_map->data[i].value;
                            if (otherFd == fd)
                                continue;
                            write(otherFd, sendMsg, sizeof(sendMsg));
                        }
                        // text filtering part, contains that text in buffer
                        if(strstr(to_lower(buffer), "i hate professor") != NULL){
                            total_client_num--;
                            // client shutdown case
                            sprintf(sendMsg, "[%s is disconnected.]\n"
                                             "[There are %d users in the chat room.]\n", client_arr[fd].nickname, total_client_num);
                            printf("%s\n", sendMsg);
                            // send msg to all client including itself
                            for(int i = 0; i < client_map->size; i++){
                                int otherFd = (unsigned char)client_map->data[i].value;
                                write(otherFd, sendMsg, sizeof(sendMsg));
                            }
                            FD_CLR(fd, &reads); //change that fd to 0
                            close(fd); // disconnect client
                            // remove client info
                            delete(client_map, client_arr[fd].nickname);
                        }
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
    while (token != NULL && count < n - 1) {
        // get next token
        token = strtok(NULL, delim);
        // save token to array
        tokens[count++] = token;
    }
    // if there's more text left, add it as the last token
    if (token != NULL && count == n - 1) {
        tokens[count++] = strtok(NULL, ""); // add string
    }
    // add null pointer, that is end of array
    tokens[count] = NULL;
    // return token array
    return tokens;
}
