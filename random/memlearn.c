#include <stdio.h>
#include <stdlib.h>
#include <string.h>


int main()
{

    printf("Memory Learning Example - Input Integers\n");
    int n, i, sum = 0;
    int *ptr;

    printf("Enter number of Integers: \n");
    scanf("%d", &n);

    ptr = (int *)malloc(n * sizeof(int));
    if (ptr == NULL)
    {
        printf("Memory allocation failed\n");
        return 1;
    }

    printf("Enter %d integers:\n", n);
    for (i = 0; i < n; i++)
    {
        scanf("%d", &ptr[i]);
        sum += *(ptr + i);
    }

    printf("The sum of the entered integers is: %d\n", sum);

    free(ptr);

    return 0;
}