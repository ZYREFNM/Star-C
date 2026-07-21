#include <stddef.h>
#include <stdlib.h>

typedef enum {} arena_state;

typedef struct {
    char* m_top;
    char* m_curr;
    char* m_last;
    size_t m_size;
    int m_ref_count;
    char* m_track[];
} t_Arena;

t_Arena* new_arena(size_t p_size) {
    t_Arena* a = malloc(sizeof(t_Arena));
    if (!a) return NULL;
    a->m_top = malloc(p_size);
    if (!a->m_top) return NULL;
    a->m_curr = a->m_top;
    a->m_last = a->m_top;
    a->m_size = p_size;
    a->m_ref_count = 0;
    return a;
}

void* arena_malloc(t_Arena* p_arena, size_t p_ref) {
    size_t used_storage = p_arena->m_curr - p_arena->m_top;
    void* alloc = p_arena->m_curr;
    p_arena->m_last = p_arena->m_curr;
    p_arena->m_curr += p_ref;
    return alloc;
}

void* arena_keeps_track_of(t_Arena* p_arena) {
    return p_arena->m_track[p_arena->m_ref_count];
}

void kill_arena(t_Arena* p_arena) {
    if (p_arena) {
        free(p_arena->m_top);
        free(p_arena);
    }
}