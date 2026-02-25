#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdalign.h>

void* safe_malloc(size_t size) {
    void* ptr = malloc(size);
    assert(ptr != NULL && "Memory allocation failed");
    return ptr;
}

void safe_free(void* ptr) {
    if (ptr != NULL) {
        free(ptr);
    }
}


typedef struct Point {
    int x;
    int y;
} Point;


void print_point_coordinates(Point* p) {
    if (p != NULL) {
        printf("Point coordinates: (%d, %d)\n", p->x, p->y);
    } else {
        printf("Point is NULL\n");
    }
}

void swap_point(Point* a, Point* b) {
    Point temp = *a;
    *a = *b;
    *b = temp;
}

void print_point_size_alignment(Point* p) {
    printf("The Size of the Point: %zu\n", sizeof(Point));
    printf("The alignement of the Point: %d\n", alignof(&p));
}


int main() {
    printf("\nGET TO THE POINT!\n\n");
    printf("<<C Program for managing POINTS>>\n\n");
    size_t point_size = sizeof(Point);
    
    Point* a = (Point*)safe_malloc(point_size);
    Point* b = (Point*)safe_malloc(point_size);
    Point* c = (Point*)safe_malloc(point_size);

    Point x;

    a->x = 10; a->y = 20;
    b->x = 30; b->y = 40;
    c->x = 50; c->y = 60;

    x.x = 70; x.y = 80;

    print_point_coordinates(a);
    print_point_coordinates(b);
    print_point_coordinates(c);
    print_point_coordinates(&x);

    printf("The size of an int: %zu bytes\n", sizeof(int));
    printf("The size of a float: %zu bytes\n", sizeof(float));
    printf("The size of a double: %zu bytes\n", sizeof(double));
    printf("The size of a char: %zu bytes\n", sizeof(char));

    printf("The align of an int: %d\n", alignof(int));
    printf("The align of a float: %d\n", alignof(float));
    printf("The align of a double: %d\n", alignof(double));
    printf("The align of a char: %d\n", alignof(char));

    print_point_size_alignment(a);
    print_point_size_alignment(b);
    print_point_size_alignment(c);
    print_point_size_alignment(&x);

    //compare pass by value vs pass by reference
    printf("Before swap: a = (%d, %d), b = (%d, %d)\n", a->x, a->y, b->x, b->y);
    swap_point(a, b);
    printf("After swap: a = (%d, %d), b = (%d, %d)\n", a->x, a->y, b->x, b->y);

    safe_free(a);
    safe_free(b);
    safe_free(c);

    Point points[] = { {1,2}, {3,4}, {5,6} };

    for (size_t i = 0; i < sizeof(points)/sizeof(points[0]); i++) {
        print_point_coordinates(&points[i]);
    }

    return 0;
}