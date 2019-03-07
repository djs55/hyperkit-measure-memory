#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdint.h>

#define PAGE_SIZE 4096

unsigned long long getTotalSystemMemory()
{
    long pages = sysconf(_SC_PHYS_PAGES);
    long page_size = sysconf(_SC_PAGE_SIZE);
    return pages * page_size;
}

static void print_usage(){
    printf("Usage:\n");
    printf("touch -m <memory in MiB>\n");
    printf("  -- allocate memory, touch all pages every second to keep it in RAM.\n");            
}

static void touch(long mib) {
    uint8_t *ptr;
    long bytes = mib * 1024 * 1024;
    if ((ptr = (uint8_t*)malloc(bytes)) == NULL) {
        perror("malloc");
        exit(EXIT_FAILURE);
    }
    fprintf(stderr, "Allocated %ld MiB\n", mib);
    int i = 1;
    for (i = 0; i < 120; i ++) {
        long count = 0;
        for (long off = 0; off < bytes; off += PAGE_SIZE) {
            *(ptr + off) = i % 256; 
            count ++;
        }
        fprintf(stderr, "Touched all %ld pages\n", count);
        sleep(1);
    }
    for (i = 0; ; i++) {
        sleep(1);
    }
}

int main(int argc, char *argv[]) {
    int option = 0;
    // add a bit extra to increase the memory pressure
    int mib = getTotalSystemMemory() / 1048576 + 50;

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
    touch(mib);
    return 0;
}