#ifndef WEB_SHELL_H
#define WEB_SHELL_H

void execCommand(const char* cmd_string, char ** out_string, unsigned long *out_size, char** err_string, unsigned long* err_size);

#endif // WEB_SHELL_H