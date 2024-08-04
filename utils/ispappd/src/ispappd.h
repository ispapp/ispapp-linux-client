#ifndef ISPAPPD_H
#define ISPAPPD_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <curl/curl.h>
#include <jansson.h>

// Configuration structure
typedef struct
{
    int enabled;
    char login[20];
    char topDomain[100];
    int topListenerPort;
    int topSmtpPort;
    char topKey[100];
    char ipbandswtestserver[100];
    char btuser[50];
    char btpwd[50];
} ispapp_config_t;

void loadConfig(ispapp_config_t *config);
void startService();
void stopService();
void handleCommand(int argc, char *argv[]);
void *updates_thread(void *arg);
void *configs_thread(void *arg);

#endif // ISPAPPD_H
