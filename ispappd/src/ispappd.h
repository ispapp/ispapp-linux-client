#ifndef ISPAPPD_H
#define ISPAPPD_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <netdb.h>
#include <unistd.h>
#include <pthread.h>
#include <curl/curl.h>
#include <jansson.h>
#include <regex.h>

/*
    requirements for ubuntu
    sudo apt install libjansson-dev libcurl4-openssl-dev
*/

// init global 
// CURLcode curl_global_init(long flags);

// Configuration structure
typedef struct
{
    int enabled;
    char login[20];
    char topDomain[100];
    int topListenerPort;
    int topSmtpPort;
    char topKey[100];
    char refreshToken[150];
    char accessToken[150];
    int connected;
    char ipbandswtestserver[100];
    char btuser[50];
    char btpwd[50];
} ispapp_config_t;

int isValidUUID(const char *uuid);
void loadConfig(ispapp_config_t *config);
void saveConfig(const ispapp_config_t *config);
void startService();
void stopService();
void handleCommand(int argc, char *argv[]);
void *updates_thread(void *arg);
void *configs_thread(void *arg);
int checkUuid(const char *domain, int port);
void *configs_thread(void *arg);
int isDomainLive(const char *domain);
void logError(const char *message);
int isPortOpen(const char *domain, int port);

struct curl_data {
    size_t size;
    char* data;
};
size_t WriteCallback(void *ptr, size_t size, size_t nmemb, struct curl_data *data);
void *healthcheck_thread(void *arg);

#endif // ISPAPPD_H
