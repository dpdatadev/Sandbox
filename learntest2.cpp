using namespace std;
#include <iostream>
#include <memory>
#include <vector>
#include <cstring>
#include <mcheck.h>

#define LABEL_SIZE 256

void* xmalloc(size_t size) {
    void *ptr = malloc(size);
    if (ptr == 0)
        throw exception();
    return ptr;
}

void* save_copy_of_string(const char* str, size_t len) {
    char *new_string = (char*)xmalloc(len + 1);
    new_string[len] = '\0';
    return memcpy(new_string, str, len);
}


struct Data {
    int value;
    char* label[LABEL_SIZE];
};


int main() {

    #ifdef DEBUGGING
    mtrace();
    #endif

    Data d;
    d.value = 42;
    
    cout << "The answer is: " << d.value << endl;

    d.value += 1;

    cout << "The new answer is: " << d.value << endl;
    return 0;
}