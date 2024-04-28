/*
 * map.c
 * 20190532 Sang yun Lee
 */

#include "map.h"

// Map init function
Map* createMap(int capacity) {
    Map *map = (Map*)malloc(sizeof(Map));
    map->data = (KeyValue*)malloc(capacity * sizeof(KeyValue));
    map->capacity = capacity;
    map->size = 0;
    return map;
}

// map insert function
void insert(Map *map, char *key, int value) {
    // if size is fulled, increase size
    if (map->size == map->capacity) {
        map->capacity *= 2;
        map->data = (KeyValue*)realloc(map->data, map->capacity * sizeof(KeyValue));
    }

    // add new key-value pair
    map->data[map->size].key = strdup(key);
    map->data[map->size].value = value;
    map->size++;
}

// delete function
void delete(Map *map, char *key) {
    int i, found = 0;
    // find key
    for (i = 0; i < map->size; i++) {
        if (strcmp(map->data[i].key, key) == 0) {
            found = 1;
            break;
        }
    }
    // if find key, delete that element
    if (found) {
        free(map->data[i].key); // free memory
        for (int j = i; j < map->size - 1; j++) {
            map->data[j] = map->data[j + 1];
        }
        map->size--;
    }
}

// get element
char find(Map *map, char *key) {
    for (int i = 0; i < map->size; i++) {
        if (strcmp(map->data[i].key, key) == 0) {
            return map->data[i].value;
        }
    }
    // do not find key, return -1
    return -1;
}

// return map size
int size(Map *map) {
    return map->size;
}

// free map function
void freeMap(Map *map) {
    for (int i = 0; i < map->size; i++) {
        // memory free
        free(map->data[i].key);
    }
    free(map->data);
    free(map);
}

// Map print function
void printMap(Map *map) {
    printf("Map Contents:\n");
    for (int i = 0; i < map->size; i++) {
        printf("[%s: %d]\n", map->data[i].key, map->data[i].value);
    }
}