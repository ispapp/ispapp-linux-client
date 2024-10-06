#include "ispappd.h"

ispapp_config_t shared_config;
pthread_t updates_thread_id, configs_thread_id, healthcheck_thread_id;

void logError(const char *message)
{
    FILE *logFile = fopen("/etc/config/ispapp_logs", "a"); // Open for appending
    if (logFile != NULL)
    {
        fprintf(logFile, "%s\n", message);
        fclose(logFile);
    }
}

int isDomainLive(const char *domain)
{
    struct hostent *he = gethostbyname(domain);
    return (he != NULL); // Returns non-zero if the domain is live
}

int isPortOpen(const char *domain, int port)
{
    struct sockaddr_in server;
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0)
        return 0;

    server.sin_family = AF_INET;
    server.sin_port = htons(port);
    server.sin_addr = *(struct in_addr *)gethostbyname(domain)->h_addr;

    // Set a timeout for the connection attempt
    struct timeval timeout;
    timeout.tv_sec = 2; // 2 seconds
    timeout.tv_usec = 0;
    setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, (const char *)&timeout, sizeof(timeout));

    int result = connect(sock, (struct sockaddr *)&server, sizeof(server));
    close(sock);
    return (result == 0); // Returns non-zero if the port is open
}

// Callback to write the response into a string
size_t WriteCallback(void *ptr, size_t size, size_t nmemb, struct curl_data *data)
{
    size_t index = data->size;
    size_t n = (size * nmemb);
    char *tmp;

    data->size += (size * nmemb);

#ifdef DEBUG
    fprintf(stderr, "data at %p size=%ld nmemb=%ld\n", ptr, size, nmemb);
#endif
    tmp = realloc(data->data, data->size + 1); /* +1 for '\0' */

    if (tmp)
    {
        data->data = tmp;
    }
    else
    {
        if (data->data)
        {
            free(data->data);
        }
        fprintf(stderr, "Failed to allocate memory.\n");
        return 0;
    }

    memcpy((data->data + index), ptr, n);
    data->data[data->size] = '\0';

    return size * nmemb;
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
        if (strstr(line, "option connected"))
            sscanf(line, " option connected '%d'", &config->connected);
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
        if (strstr(line, "option refreshToken"))
            sscanf(line, " option refreshToken '%[^']", config->refreshToken);
        if (strstr(line, "option accessToken"))
            sscanf(line, " option accessToken '%[^']", config->accessToken);
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
    if (initConfig())
    {
        char *operation = "Starting";
        if (shared_config.enabled != 1)
        {
            printf("service already running!\n");
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
    else
    {
        printf("Failed to initialize config\n");
    }
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

void *healthcheck_thread(void *arg)
{
    // Healthcheck thread logic
    while (1)
    {
        // Only run UUID replacement if current UUID is '00:00:00:00:00:00'
        if (!isDomainLive(shared_config.topDomain) || !isPortOpen(shared_config.topDomain, shared_config.topListenerPort))
        {
            stopService();
            return 0; // Failed
            // todo
        }

        sleep(60); // Periodic health checks every minute
    }
    return NULL;
}

int isValidUUID(const char *uuid)
{
    // Regular expression pattern for UUID
    const char *pattern = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$";
    regex_t regex;

    // Compile the regex
    if (regcomp(&regex, pattern, REG_EXTENDED) != 0)
    {
        return 0; // Failed to compile regex
    }

    // Execute regex
    int status = regexec(&regex, uuid, 0, NULL, 0);
    regfree(&regex); // Free the compiled regex

    return (status == 0); // Returns non-zero if the UUID matches the pattern
}

int checkUuid(const char *domain, int port)
{
    CURL *curl;
    CURLcode res;
    char url[256];
    struct curl_data response;
    response.size = 0;
    response.data = malloc(4096); /* reasonable size initial buffer */
    snprintf(url, sizeof(url), "https://%s:%d/auth/uuid", domain, port);

    curl_global_init(CURL_GLOBAL_DEFAULT);
    /* cache the CA cert bundle in memory for a week */
    // curl_easy_setopt(curl, CURLOPT_CA_CACHE_TIMEOUT, 604800L);
    curl = curl_easy_init();
    if (curl)
    {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response); // Store response in this buffer
        // #ifdef SKIP_PEER_VERIFICATION
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L);
        // #endif
        // #ifdef SKIP_HOSTNAME_VERIFICATION
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);
        // #endif
        res = curl_easy_perform(curl);
        if (res == CURLE_OK)
        {
            printf("is valid uuid (%s): %s -> %d (%d)\n", url, response.data, isValidUUID(response.data), res);
            if (isValidUUID(response.data))
            {
                printf("Valid UUID: %s\n", response.data);
                if (strcmp(shared_config.login, "00:00:00:00:00:00") != -1)
                {
                    strcpy(shared_config.login, response.data);
                    saveConfig(&shared_config);
                }
                return 1;
            }
            else
            {
                logError("Invalid UUID format received");
            }
        }
        else
        {
            logError(curl_easy_strerror(res));
            fprintf(stderr, "curl_easy_perform() failed: %s\n",
                    curl_easy_strerror(res));
        }

        curl_easy_cleanup(curl);
    }
    curl_global_cleanup();
    return 0; // Failed to check UUID
}

int initConfig()
{
    CURL *curl;
    CURLcode res;
    char url[256];
    struct curl_data response;
    response.data = malloc(4096);
    response.size = 0;

    curl_global_init(CURL_GLOBAL_DEFAULT);
    if (strcmp(shared_config.login, "00:00:00:00:00:00") == 0)
    {
        if (!checkUuid(shared_config.topDomain, shared_config.topListenerPort))
        {
            return 0;
        }
    }
    snprintf(url, sizeof(url), "https://%s:%d/initconfig?login=%s&key=%s",
             shared_config.topDomain, shared_config.topListenerPort,
             shared_config.login, shared_config.topKey);

    curl = curl_easy_init();
    if (curl)
    {
        printf("sending request with login: %s\n", shared_config.login);
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response); // Store response in this buffer
        // #ifdef SKIP_PEER_VERIFICATION
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L);
        // #endif
        // #ifdef SKIP_HOSTNAME_VERIFICATION
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);
        res = curl_easy_perform(curl);
        if (res == CURLE_OK && response.size > 0)
        {
            printf("req: %s => %s \n", url, response.data);
            json_error_t error;
            json_t *root = json_loads(response.data, 0, &error);
            if (root)
            {
                json_t *accessTokenJson = json_object_get(root, "accessToken");
                json_t *refreshTokenJson = json_object_get(root, "refreshToken");

                if (accessTokenJson && refreshTokenJson)
                {
                    const char *accessToken = json_string_value(accessTokenJson);
                    const char *refreshToken = json_string_value(refreshTokenJson);

                    if (accessToken && refreshToken)
                    {
                        strcpy(shared_config.accessToken, accessToken);   // Save accessToken
                        strcpy(shared_config.refreshToken, refreshToken); // Save refreshToken

                        // Update config and start threads
                        shared_config.connected = 1;
                        shared_config.enabled = 1;
                        saveConfig(&shared_config); // Save updated config to file
                        json_decref(root);          // Decrement reference count
                        return 1;                   // Success
                    }
                }
                else
                {
                    logError("Missing accessToken or refreshToken in response");
                }
                json_decref(root); // Clean up
            }
            else
            {
                logError("JSON parsing error");
            }
        }
        else
        {
            logError(curl_easy_strerror(res));
        }
        curl_easy_cleanup(curl);
    }
    curl_global_cleanup();
    return 0; // Failure
}

void saveConfig(const ispapp_config_t *config)
{
    FILE *fp = fopen("/etc/config/ispapp", "w");
    if (fp != NULL)
    {
        fprintf(fp, "config settings\n");
        fprintf(fp, "\toption enabled '%d'\n", config->enabled);
        fprintf(fp, "\toption connected '%d'\n", config->connected);
        fprintf(fp, "\toption login '%s'\n", config->login);
        // fprintf(fp, "\toption topDomain '%s'\n", config->topDomain);
        fprintf(fp, "\toption topListenerPort '%d'\n", config->topListenerPort);
        fprintf(fp, "\toption topSmtpPort '%d'\n", config->topSmtpPort);
        fprintf(fp, "\toption topKey '%s'\n", config->topKey);
        fprintf(fp, "\toption refreshToken '%s'\n", config->refreshToken);
        fprintf(fp, "\toption accessToken '%s'\n", config->accessToken);
        fprintf(fp, "\toption ipbandswtestserver '%s'\n", config->ipbandswtestserver);
        fprintf(fp, "\toption btuser '%s'\n", config->btuser);
        fprintf(fp, "\toption btpwd '%s'\n", config->btpwd);
        fclose(fp);
    }
    else
    {
        perror("Failed to open config file for writing");
    }
}

int main(int argc, char *argv[])
{
    loadConfig(&shared_config);
    handleCommand(argc, argv);
    return 0;
}
