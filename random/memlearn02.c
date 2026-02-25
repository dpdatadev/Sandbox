#include <stdio.h>

typedef struct counter {
    int val;
} counter;

// todo

int main() {

    char* response;
    
    counter *c = (counter*)malloc(sizeof(counter));

    printf("Please enter how many times you want me to say Hello\n");
    scanf("Enter Answer: %d", &response);


    printf("::GOODBYE::\n");
    return 0;
}