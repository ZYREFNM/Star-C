#ifndef _STAR_C_RUNTIME
#define _STAR_C_RUNTIME
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stddef.h>
#include "arena.h"

typedef struct {
    char* lit;
    size_t length;
} string;

string string_new(char* cstr);
string star_concat(string to, string from);

#endif