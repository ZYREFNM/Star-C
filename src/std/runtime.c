#include <stdlib.h>
#include <stdio.h>
#include <string.h>

char* star_concat(char* to, char* from) {
    if (!to) to = "";
    if (!from) from = "";
    
    char* result = malloc(strlen(to) + strlen(from) + 1);
    if (result) {
        strcpy(result, to);
        strcat(result, from);
    }
    free(result)
    return result;
}