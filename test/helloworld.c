#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include "../src/compiler/runtime.h"
#define __ARENA_ALLOCATOR_SIZE__ 1000
t_Arena* __STAR_C_ARENA_ALLOCATOR__;



int star_main() {
	printf(string_new("Hello World\n").lit);
return 0;
}

int main() {
__STAR_C_ARENA_ALLOCATOR__ = new_arena(__ARENA_ALLOCATOR_SIZE__ * 1024);
int exit = star_main();
kill_arena(__STAR_C_ARENA_ALLOCATOR__);
return exit;
}
