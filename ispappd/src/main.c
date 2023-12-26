#include "main.h"
// just an example to how the agent will work in core function

// Function to handle the 'status' command
void handleStatus() {
    printf("Executing 'ispappd status'\n");
}

// Function to handle the 'start' command
void handleStart() {
    printf("Executing 'ispappd start'\n");
}

// Function to handle the 'stop' command
void handleStop() {
    printf("Executing 'ispappd stop'\n");
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Usage: %s <command>\n", argv[0]);
        return 1;
    }
    if (strcmp(argv[1], "status") == 0) {
        handleStatus();
    } else if (strcmp(argv[1], "start") == 0) {
        handleStart();
    } else if (strcmp(argv[1], "stop") == 0) {
        handleStop();
    } else {
        printf("Unknown command: %s\n", argv[1]);
        return 1;
    }

    return 0;
}
