/*
 * MultiTCPServer.c
 * 20190532 Sang yun Lee
 */


// https://m.blog.naver.com/whtie5500/221692806173
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include <sys/socket.h>
#include <netinet/in.h>


#define BUFSIZE 100
#define PORT 20532

int main(void){
    int serv_sock;
    struct sockaddr_in serv_addr;

    fd_set reads, temps;
    int fd_max;
    int binded;

    char message[BUFSIZE];
    struct timeval timeout;

    int fd, str_len;
    int clnt_sock, clnt_len;
    struct sockaddr_in clnt_addr;
    int opt = 1;

    serv_sock = socket(AF_INET, SOCK_STREAM, 0);

    //1
    setsockopt(serv_sock, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt));
    bzero(&serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = INADDR_ANY;
    serv_addr.sin_port = htons(PORT);

    binded = bind(serv_sock, (struct sockaddr*)&serv_addr, sizeof(serv_addr));
    if(binded != 0){
        perror("bind error");
        exit(-1);
    }

    binded = listen(serv_sock, 5);
    if(binded != 0){
        perror("listen error");
        exit(-1);
    }
    printf("server listening\n");

    FD_ZERO(&reads); //2
    FD_SET(serv_sock, &reads); //3
    fd_max = serv_sock; //4
    printf("serv_sock is : %d\n", serv_sock);

    while(1){
        temps = reads;
        timeout.tv_sec = 5;
        timeout.tv_usec = 0;
        //5
        if(select(fd_max + 1, &temps, 0, 0, &timeout) < 0){
            perror("select error");
            exit(1);
        }
        //6
        for(fd = 0; fd < fd_max+1; fd++){
            //7
            if(FD_ISSET(fd, &temps)){
                //8 연결요청인경우
                if(fd == serv_sock){
                    clnt_len = sizeof(clnt_addr);
                    clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_len);
                    //해당 클라이언트와 연결시킨다.

                    FD_SET(clnt_sock, &reads); //클라이언트 번지..?에 1을 넣는다.
                    //9
                    if(fd_max < clnt_sock){
                        fd_max = clnt_sock;
                    }
                    printf("client access : socket discriptor %d\n", clnt_sock); //해당 client socket번지에
                    // fd_set을 1로 변경하였다~
                } else{
                    str_len = read(fd, message, BUFSIZE); //10
                    //11 연결 종료 요청일경우
                    if(str_len == 0){
                        FD_CLR(fd, &reads); //해당 파일디스크립터 fd를 0으로 변경
                        close(fd);
                        printf("client end : socket discriptor %d\n", fd);
                    } else { //12
                        message[str_len] = '\0';
                        printf("client : %d message : %s\n", fd, message);
                        write(fd, message, str_len);
                    }
                }
            }
        }
    }
    return 0;
}