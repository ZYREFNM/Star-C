#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include "src/compiler/runtime.h"

#define PACKAGE "IO"
typedef struct {
} Input;
void Input_printf(char* fmt, char* text) {
	printf(fmt, text);
}
void Input_printfln(char* fmt, char* text) {
	printf(star_concat(fmt, "\n"), text);
}
typedef struct {
} Output;



int main() {
	Input_printfln("%s", "Hello, world!");
	return 0;
}
