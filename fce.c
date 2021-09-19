#include <stdio.h>

int writeln(int x) {
    printf("%d\n", x);
    return 0;
}
int write(char* x) {
    printf("%s", x);
    return 0;
}
int readln(int *x) {
    scanf("%d", x);
    return 0;
}
