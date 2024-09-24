#ifndef ISPAPPD_H
#define ISPAPPD_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <curl/curl.h>
#include <jansson.h>

/*
    requirements for ubuntu
    sudo apt install libjansson-dev libcurl4-openssl-dev
*/
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
void *configs_thread(void *arg);
int isDomainLive(const char *domain);
void logError(const char *message);
int isPortOpen(const char *domain, int port);
static size_t WriteCallback(void *contents, size_t size, size_t nmemb, void *userp);
int checkUuid(const char *domain, int port);
void *healthcheck_thread(void *arg);

#endif // ISPAPPD_H
