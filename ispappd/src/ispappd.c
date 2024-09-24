#include "ispappd.h"
#include <curl/curl.h>
#include <netdb.h>
#include <unistd.h>

ispapp_config_t shared_config;
pthread_t updates_thread_id, configs_thread_id, healthcheck_thread_id;

void logError(const char *message)
{
    FILE *logFile = fopen("/etc/config/ispapp_logs", "w");  // Open for writing, truncating any existing content
    if (logFile != NULL)
    {
        fprintf(logFile, "%s\n", message);
        fclose(logFile);
    }
}

int isDomainLive(const char *domain)
{
    struct hostent *he = gethostbyname(domain);
    return (he != NULL);  // Returns non-zero if the domain is live
}

int isPortOpen(const char *domain, int port)
{
    struct sockaddr_in server;
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) return 0;

    server.sin_family = AF_INET;
    server.sin_port = htons(port);
    server.sin_addr = *(struct in_addr *)gethostbyname(domain)->h_addr;

    // Set a timeout for the connection attempt
    struct timeval timeout;
    timeout.tv_sec = 2;  // 2 seconds
    timeout.tv_usec = 0;
    setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, (const char*)&timeout, sizeof(timeout));

    int result = connect(sock, (struct sockaddr *)&server, sizeof(server));
    close(sock);
    return (result == 0);  // Returns non-zero if the port is open
}

// Callback to write the response into a string
static size_t WriteCallback(void *contents, size_t size, size_t nmemb, void *userp)
{
    size_t totalSize = size * nmemb;
    strncat(userp, contents, totalSize);
    return totalSize;
}


int checkUuid(const char *domain, int port)
{
    CURL *curl;
    CURLcode res;
    char url[256];
    char response[256] = {0}; 
    snprintf(url, sizeof(url), "http://%s:%d/auth/uuid", domain, port);

    curl = curl_easy_init();
    if (curl)
    {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, response);  // Store response in this buffer

        char *response;
        res = curl_easy_perform(curl);
        if (res == CURLE_OK)
        {
            printf("UUID: %s\n", response);
            // Assuming response contains a valid UUID
            // strcpy(shared_config.login, uuid);  // Extracted UUID from response
            return 1;  // Successful UUID retrieval
        }
        else
        {
            logError(curl_easy_strerror(res));
        }

        curl_easy_cleanup(curl);
    }
    return 0;  // Failed to check UUID
}

void loadConfig(ispapp_config_t *config)
{
    // Load configurations from /etc/config/ispapp
    FILE *fp = fopen("/etc/config/ispapp", "r");
    if (fp == NULL)
    {
        perror("Failed to open config file");
        return;
    }

    char line[256];
    while (fgets(line, sizeof(line), fp))
    {
        if (strstr(line, "option enabled"))
            sscanf(line, " option enabled '%d'", &config->enabled);
        if (strstr(line, "option login"))
            sscanf(line, " option login '%[^']", config->login);
        if (strstr(line, "option topDomain"))
            sscanf(line, " option topDomain '%[^']", config->topDomain);
        if (strstr(line, "option topListenerPort"))
            sscanf(line, " option topListenerPort '%d'", &config->topListenerPort);
        if (strstr(line, "option topSmtpPort"))
            sscanf(line, " option topSmtpPort '%d'", &config->topSmtpPort);
        if (strstr(line, "option topKey"))
            sscanf(line, " option topKey '%[^']", config->topKey);
        if (strstr(line, "option ipbandswtestserver"))
            sscanf(line, " option ipbandswtestserver '%[^']", config->ipbandswtestserver);
        if (strstr(line, "option btuser"))
            sscanf(line, " option btuser '%[^']", config->btuser);
        if (strstr(line, "option btpwd"))
            sscanf(line, " option btpwd '%[^']", config->btpwd);
    }
    fclose(fp);
}

void updateLoginInConfig(const char *newLogin)
{
    FILE *fp = fopen("/etc/config/ispapp", "r");
    if (fp == NULL)
    {
        perror("Failed to open config file");
        return;
    }

    // Read the entire file into memory
    char *buffer = NULL;
    size_t size = 0;
    getline(&buffer, &size, fp); // This reads the entire line into the buffer
    fclose(fp);

    // Modify the login line
    char *pos = strstr(buffer, "option login");
    if (pos)
    {
        snprintf(pos, size - (pos - buffer), " option login '%s'\n", newLogin);
    }

    // Write back to the file
    fp = fopen("/etc/config/ispapp", "w");
    if (fp != NULL)
    {
        fputs(buffer, fp);
        fclose(fp);
    }
    free(buffer); // Don't forget to free the buffer allocated by getline
}

void startService()
{
    // Start the service based on the configuration
    char* operation = "Starting";
    if (shared_config.enabled != 1)
    {
        stopService();
        return;
    }
    printf("%s ispappd with configuration:\n", operation);
    printf("Enabled: %d\n", shared_config.enabled);
    printf("Login: %s\n", shared_config.login);
    printf("Top Domain: %s\n", shared_config.topDomain);
    printf("Top Listener Port: %d\n", shared_config.topListenerPort);
    printf("Top SMTP Port: %d\n", shared_config.topSmtpPort);
    printf("Top Key: %s\n", shared_config.topKey);
    printf("IP Bandwidth Test Server: %s\n", shared_config.ipbandswtestserver);
    printf("BT User: %s\n", shared_config.btuser);
    printf("BT Password: %s\n", shared_config.btpwd);

    // Start threads
    
    pthread_create(&updates_thread_id, NULL, updates_thread, NULL);
    pthread_create(&configs_thread_id, NULL, configs_thread, NULL);
    pthread_create(&healthcheck_thread_id, NULL, healthcheck_thread, NULL);
}

void stopService()
{
    // Stop the service (terminate threads)
    if (updates_thread_id)
    {
        pthread_cancel(updates_thread_id);
        pthread_join(updates_thread_id, NULL);
    }
    if (configs_thread_id)
    {
        pthread_cancel(configs_thread_id);
        pthread_join(configs_thread_id, NULL);
    }
    if (healthcheck_thread_id)
    {
        pthread_cancel(healthcheck_thread_id);
        pthread_join(healthcheck_thread_id, NULL);
    }
    printf("Service stopped\n");
}

void handleCommand(int argc, char *argv[])
{
    if (argc < 2)
    {
        printf("Usage: %s <command>\n", argv[0]);
        return;
    }
    if (strcmp(argv[1], "start") == 0)
    {
        startService();
    }
    else if (strcmp(argv[1], "stop") == 0)
    {
        stopService();
    }
    else
    {
        printf("Unknown command: %s\n", argv[1]);
    }
}

void *updates_thread(void *arg)
{
    // Updates thread logic
    while (1)
    {
        // Perform periodic updates
        sleep(60); // Placeholder: Adjust as needed
    }
    return NULL;
}

void *configs_thread(void *arg)
{
    // Configs thread logic
    while (1)
    {
        // Perform periodic config checks
        sleep(60); // Placeholder: Adjust as needed
    }
    return NULL;
}

void *healthcheck_thread(void *arg)
{
    // Healthcheck thread logic
    while (1)
    {
        // Only run UUID replacement if current UUID is '00:00:00:00:00:00'
        if (strcmp(shared_config.login, "00:00:00:00:00:00") == 0)
        {
            if (isDomainLive(shared_config.topDomain))
            {
                if (isPortOpen(shared_config.topDomain, shared_config.topListenerPort))
                {
                    if (checkUuid(shared_config.topDomain, shared_config.topListenerPort))
                    {
                        // Update the login in the configuration file if successful
                        updateLoginInConfig(shared_config.login);  // Assuming login is updated with the new UUID
                    }
                }
            }
        }

        sleep(60);  // Periodic health checks every minute
    }
    return NULL;
}

int main(int argc, char *argv[])
{
    loadConfig(&shared_config);
    handleCommand(argc, argv);
    return 0;
}
