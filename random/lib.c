#include <stdint.h>
#include <stdlib.h>
#include <assert.h>
#include <stdbool.h>
#include <stdio.h>

static char *file_get_first_line(char *filename)
{
    char *ret = NULL;
    FILE *f = fopen(filename, "r");
    if (f)
    {
        char *line;
        size_t n = 0;
        if (getline(&line, &n, f) > 0)
            ret = line;
        fclose(f);
    }
    return ret;
}

int swap(int *a, int *b)
{
    int temp = *a;
    *a = *b;
    *b = temp;
    return 0;
}

void *safe_malloc(size_t size)
{
    void *ptr = malloc(size);
    assert(ptr != NULL && "Memory allocation failed");
    return ptr;
}

void safe_free(void *ptr)
{
    if (ptr != NULL)
    {
        free(ptr);
    }
}

bool is_aligned(void *ptr, size_t alignment)
{
    return ((uintptr_t)ptr % alignment) == 0;
}

void *aligned_malloc(size_t size, size_t alignment)
{
    void *ptr = NULL;
    int result = posix_memalign(&ptr, alignment, size);
    assert(result == 0 && "Aligned memory allocation failed");
    return ptr;
}

int alignment_test()
{
    size_t size = 1024;
    size_t alignment = 64;

    // Allocate aligned memory
    void *aligned_ptr = aligned_malloc(size, alignment);
    int check = is_aligned(aligned_ptr, alignment);
    if (!check)
    {
        // This should never happen due to the assertion in aligned_malloc
        safe_free(aligned_ptr);
        return -1;
    }

    printf("Memory allocated at %p is aligned to %zu bytes\n", aligned_ptr, alignment);

    // Use the memory (omitted for brevity)

    // Free the allocated memory
    safe_free(aligned_ptr);

    return 0;
}

struct Data
{
    int a;
    float b;
    char c;
};

typedef struct AlignedData
{
    struct Data data;
    // Padding to ensure 16-byte alignment
    char padding[12];
} AlignedData;

// Padding Example
int main()
{

    size_t test_size = sizeof(struct Data);
    struct Data *data = (struct Data *)safe_malloc(test_size);
    printf("Size of Data struct: %zu bytes\n", test_size);
    safe_free(data);

    AlignedData *aligned_data = (AlignedData *)aligned_malloc(sizeof(AlignedData), 16);
    assert(is_aligned(aligned_data, 16) && "AlignedData is not 16-byte aligned");
    printf("AlignedData allocated at %p is 16-byte aligned\n", aligned_data);
    printf("Size of AlignedData struct: %zu bytes\n", sizeof(AlignedData));

    printf("Size: %zu bytes\n", sizeof(aligned_data->data));

    safe_free(aligned_data);
    return 0;
}
