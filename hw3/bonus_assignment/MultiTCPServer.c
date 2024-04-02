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

typedef struct s_clientInfo{
    char* ip;
    int port;
    int isvalid;
    int num;
}client_info;

int main(void){
    int total_client_num = 0;
    client_info client_arr[MAX_CLIENT + 3];
    for(int i = 0; i < MAX_CLIENT + 3; i++)
        client_arr[i].isvalid = 0;

    int serv_sock = socket(AF_INET, SOCK_STREAM, 0);
    int opt = 1;
    setsockopt(serv_sock, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt));

    struct sockaddr_in serv_addr; //서버 주소 초기화
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

    fd_set reads, temps; // setect 함수에 사용되는 fd_set 구조체 선언
    FD_ZERO(&reads); // fd_set을 0으로 초기화
    FD_SET(serv_sock, &reads); // server listen 소캣을 파일 디스크립터를 fd_set에 등록
    int fd_max = serv_sock; // 서버 소캣을 최대 fd로 설정, 이 값에 +1을 한 값을 반복문으로 처리해서 이후에 있는 클라이언트의 요청을 받는다.
    printf("serv_sock is : %d\n", serv_sock);

    while(1){
        struct timeval timeout;
        timeout.tv_sec = 5;
        timeout.tv_usec = 0;

        temps = reads; // select는 원본값을 변경하기 때문에 값을 복사해준다.
        if(select(fd_max + 1, &temps, 0, 0, &timeout) < 0){ // 이 함수에서 변화를 감지한다. 감지하면 그 fd를 1로 만들어줌
            perror("select error");
            exit(1);
        }
        for(int fd = 0; fd < fd_max + 1; fd++){ // 어떤 fd에서 변화가 일어났는지 확인
            if(FD_ISSET(fd, &temps)){ // client가 요청을 한 경우 fd가 1인 경우를 확인, 즉 변화가 있는 경우, 2가지 경우가 존재
                if(fd == serv_sock){ // server의 연결 요청인 경우, 즉 서버라면 다시 들어가서 요청을 처리함
                    struct sockaddr_in clnt_addr;
                    socklen_t clnt_len = sizeof(clnt_addr);
                    //해당 클라이언트와 연결시킨다.
                    int clnt_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_len);
                    char *client_ip = inet_ntoa(clnt_addr.sin_addr);
                    int client_port = ntohs(clnt_addr.sin_port);
                    printf("Connection request from %s:%d\n", client_ip, client_port);
                    client_arr[clnt_sock].ip = client_ip;
                    client_arr[clnt_sock].port = client_port;
                    client_arr[clnt_sock].isvalid = 1;
                    client_arr[clnt_sock].num = ++total_client_num;

                    FD_SET(clnt_sock, &reads); //연결했으므로 해당 원본 set에 1을 넣는다.
                    if(fd_max < clnt_sock){ //9 - 클라이언트가 들어올때마다 최대 fd_max를 갱신
                        fd_max = clnt_sock;
                    }
                } else{ //이미 연결된 클라이언트의 요청
                    char message[BUFSIZE];
                    int str_len = read(fd, message, BUFSIZE); //10
                    if(str_len == 0){ //11 연결 종료 요청일경우 여기서는 처리 필요
                        FD_CLR(fd, &reads); //해당 파일디스크립터 fd를 0으로 변경
                        close(fd);
                        total_client_num--;
                        printf("client end : socket discriptor %d\n", fd);
                    } else {
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

