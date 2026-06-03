#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include "src/compiler/runtime.h"


typedef struct {
} Input;
typedef struct {
} Output;
void IO__Output_print(char* text) {
	printf(text);
}
void IO__Output_printf(char* fmt, char* text) {
	printf(fmt, text);
}
void IO__Output_printfln(char* fmt, char* text) {
	printf(star_concat(fmt, "\n"), text);
}



