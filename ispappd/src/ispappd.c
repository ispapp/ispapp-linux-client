#include "ispappd.h"

ispapp_config_t shared_config;
pthread_t updates_thread_id, configs_thread_id;

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
            sscanf(line, " option login '%s'", config->login);
        if (strstr(line, "option topDomain"))
            sscanf(line, " option topDomain '%s'", config->topDomain);
        if (strstr(line, "option topListenerPort"))
            sscanf(line, " option topListenerPort '%d'", &config->topListenerPort);
        if (strstr(line, "option topSmtpPort"))
            sscanf(line, " option topSmtpPort '%d'", &config->topSmtpPort);
        if (strstr(line, "option topKey"))
            sscanf(line, " option topKey '%s'", config->topKey);
        if (strstr(line, "option ipbandswtestserver"))
            sscanf(line, " option ipbandswtestserver '%s'", config->ipbandswtestserver);
        if (strstr(line, "option btuser"))
            sscanf(line, " option btuser '%s'", config->btuser);
        if (strstr(line, "option btpwd"))
            sscanf(line, " option btpwd '%s'", config->btpwd);
    }
    fclose(fp);
}

void startService()
{
    // Start the service based on the configuration
    if (shared_config.enabled != 1)
    {
        stopService();
        return;
    }

    printf("Starting ispappd with configuration:\n");
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
}

void stopService()
{
    // Stop the service (terminate threads)
    pthread_cancel(updates_thread_id);
    pthread_cancel(configs_thread_id);
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

int main(int argc, char *argv[])
{
    loadConfig(&shared_config);
    handleCommand(argc, argv);
    return 0;
}
