/*
 * MultiTCPServer.c
 * 20190532 Sang yun Lee
 */

/*
 * 참고자료
 * 1. fd_set에 대한 설명 https://blog.naver.com/tipsware/220810795410
 * 2. 이 코드의 reference 한글 https://m.blog.naver.com/whtie5500/221692806173
 */
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h> // inet_ntoa


#define BUFSIZE 100
#define PORT 20532
#define MAX_CLIENT 100

int main(void){
    int serv_sock = socket(AF_INET, SOCK_STREAM, 0);
    // 서버 소캣 생성, AF_INET이 TCP를 가리킴
    // 소캣 옵션 생성
    int opt = 1; // 단순 옵션 정보를 저장하는 변수
    setsockopt(serv_sock, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt));

    // 서버 주소 초기화
    struct sockaddr_in serv_addr;
    bzero(&serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET; //주소 체계, internet protocol
    serv_addr.sin_addr.s_addr = htonl(INADDR_ANY); // 여러 ip주소에서 들어오는 데이터를 모두처리
    serv_addr.sin_port = htons(PORT); // big endian to little endian

    // 서버 소캣 파일 디스크립터에 소캣 정보 (serv_addr) 바인딩
    int binded;
    binded = bind(serv_sock, (struct sockaddr*)&serv_addr, sizeof(serv_addr));
    if(binded != 0){
        perror("bind error");
        exit(1);
    }

    // 클라이언트의 연결요청 대기, 뒤에 값은 클라이언트 대기 큐에 있는 요청 최대 개수
    binded = listen(serv_sock, 5);
    if(binded != 0){
        perror("listen error");
        exit(1);
    }
    printf("Server is ready to receive on port %d\n", PORT);

    // setect 함수에 사용되는 fd_set 구조체 선언
    fd_set reads, temps;
    // fd_set을 0으로 초기화
    FD_ZERO(&reads);
    // server listen 소캣을 파일 디스크립터를 fd_set에 등록
    FD_SET(serv_sock, &reads);
    // 서버 소캣을 최대 fd로 설정, 이 값에 +1을 한 값을 반복문으로 처리해서 이후에 있는 클라이언트의 요청을 받는다.
    int fd_max = serv_sock;
    printf("serv_sock is : %d\n", serv_sock);

    while(1){
        struct timeval timeout;
        timeout.tv_sec = 5;
        timeout.tv_usec = 0;

        // select는 원본값을 변경하기 때문에 값을 복사해준다.
        temps = reads;
        // 이 함수에서 변화를 감지한다. 감지하면 그 fd를 1로 만들어줌
        if(select(fd_max + 1, &temps, 0, 0, &timeout) < 0){
            perror("select error");
            exit(1);
        }
        // 어떤 fd에서 변화가 일어났는지 확인
        for(int fd = 0; fd < fd_max + 1; fd++){
            // client가 요청을 한 경우 fd가 1인 경우를 확인, 즉 변화가 있는 경우, 2가지 경우가 존재
            if(FD_ISSET(fd, &temps)){
                // server의 연결 요청인 경우, 즉 서버라면 다시 들어가서 요청을 처리함
                if(fd == serv_sock){
                    struct sockaddr_in clnt_addr;
                    socklen_t clnt_len = sizeof(clnt_addr);
                    //해당 클라이언트와 연결시킨다.
                    int clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_len);
                    char *client_ip = inet_ntoa(clnt_addr.sin_addr);
                    int client_port = ntohs(clnt_addr.sin_port);
                    printf("Connection request from %s:%d\n", client_ip, client_port);

                    FD_SET(clnt_sock, &reads); //연결했으므로 해당 원본 set에 1을 넣는다.
                    //9 - 클라이언트가 들어올때마다 최대 fd_max를 갱신
                    if(fd_max < clnt_sock){
                        fd_max = clnt_sock;
                    }
                    //해당 client socket번지에 fd_set을 1로 변경하였다~
                } else{ //이미 연결된 클라이언트의 요청
                    char message[BUFSIZE];
                    int str_len = read(fd, message, BUFSIZE); //10
                    //11 연결 종료 요청일경우 여기서는 처리 필요
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
//
//void clientHandler(){
//
//}

