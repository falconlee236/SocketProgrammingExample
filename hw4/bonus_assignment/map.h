/*
 * map.h
 * 20190532 Sang yun Lee
 */

#ifndef __MAP_H__
#define __MAP_H__

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// store key, value struct
typedef struct KeyValue {
    char *key;
    char value;
} KeyValue;

// Map data struct
typedef struct Map {
    KeyValue *data; // key-value pair store list
    int capacity;   // capacity of list
    int size;       // number of elements
} Map;

// Map init function
Map* createMap(int capacity);
// map insert function
void insert(Map *map, char *key, int value);
// delete function
void delete(Map *map, char *key);
// get element
char find(Map *map, char *key);
// return map size
int size(Map *map);
// free map function
void freeMap(Map *map);
// Map print function
void printMap(Map *map);

#endif
