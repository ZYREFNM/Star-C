#ifndef _STAR_C_ARENA_ALLOC
#define _STAR_C_ARENA_ALLOC

#include <stddef.h>
#include <stdlib.h>

typedef struct {
    char* m_top;
    char* m_curr;
    size_t m_size;
    int m_ref_count;
} t_Arena;

t_Arena* new_arena(const size_t p_size);
void* arena_malloc(t_Arena* p_arena, size_t p_size);
void* kill_arena(t_Arena* p_arena);
#endif