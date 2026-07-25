#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stddef.h>
#include "arena.h"

typedef struct {
    char* lit;
    size_t length;
} string;

string string_new(const char* cstr) {
    string starstr;
    if (!cstr) cstr = "";
    size_t len = strlen(cstr);
    char* result = malloc(len + 1);
    if (result) {
        strcpy(result, cstr);
        starstr.lit = result;
        starstr.length = strlen(starstr.lit);
    } else {
        starstr.lit = NULL;
        starstr.length = 0;
    }
    return starstr;
}

typedef char byte;

string star_concat(string to, string from) {
    if (!to.lit) to.lit = "";
    if (!from.lit) from.lit = "";
    
    string res;
    char* result = malloc(strlen(to.lit) + strlen(from.lit) + 1);
    if (result) {
        strcpy(result, to.lit);
        strcat(result, from.lit);
    }
    res = string_new(result);
    free(result);
    return res;
}