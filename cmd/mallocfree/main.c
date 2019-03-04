#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#define PAGE_SIZE 4096

static void print_usage(){
    printf("Usage:\n");
    printf("mallocfree -m <memory in MiB>\n");
    printf("  -- allocate 2 large blocks of memory, touch all pages, then free one block\n");            
}

static uint8_t *allocate(long mib) {
    uint8_t *ptr;
    long bytes = mib * 1024 * 1024;
    if ((ptr = (uint8_t*)malloc(bytes)) == NULL) {
        perror("malloc");
        exit(EXIT_FAILURE);
    }
    fprintf(stderr, "Allocated %ld MiB\n", mib);
    return ptr;
}

static void touch(uint8_t *ptr, long mib) {
    long bytes = mib * 1024 * 1024;
    long count = 0;
    for (long off = 0; off < bytes; off += PAGE_SIZE) {
        *(ptr + off) = 1; 
        count ++;
    }
    fprintf(stderr, "Touched all %ld pages\n", count);
}

int main(int argc, char *argv[]) {
    int option = 0, mib;

    while ((option = getopt(argc, argv,"m:")) != -1) {
        switch (option) {
             case 'm':
                 mib = atol(optarg); 
                 break;
             default: print_usage(); 
                 exit(EXIT_FAILURE);
        }
    }

    printf("mib: %d\n", mib);
    uint8_t *one = allocate(mib);
    uint8_t *two = allocate(mib);
    touch(one, mib);
    touch(two, mib);
    sleep(30);
    free(one);
    fprintf(stderr, "Freed one block of memory\n");
    while (1){
        sleep(1);
        fprintf(stderr, ".");
    }
    return 0;
}
