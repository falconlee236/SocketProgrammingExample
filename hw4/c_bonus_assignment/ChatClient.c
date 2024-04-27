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

int main(int ac, char **av){
    // system args handling
    if (ac != 2) {
        printf("This program need to one argument\n");
        exit(1);
    }
    // get nickname from system args
    char *nickname = av[1];
    // create commandMap
    Map *commandMap = createMap(5); // 초기 용량 5인 Map 생성
    // insert init command to Map
    insert(commandMap, "ls", 1);
    insert(commandMap, "secret", 2);
    insert(commandMap, "except", 3);
    insert(commandMap, "ping", 4);
    insert(commandMap, "quit", 5);

    // memory free in map
    freeMap(commandMap);
    return 0;
}