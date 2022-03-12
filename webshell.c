#include "webshell.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <limits.h>

int popenTHREE(int *threepipe, const char *command)
{
  // threepipe[0] is the stdin fd
  // threepipe[1] is the stdout fd
  // threepipe[2] is the stderr fd

  int in[2];
  int out[2];
  int err[2];
  int pid;
  int rc;

  // pipe returns r[0] : the fd for the read end
  // and r[1] : the fd for the write end
  // this means you have to open the fd with fdopen()
  rc = pipe(in);
  if (rc < 0)
    goto error_in;

  rc = pipe(out);
  if (rc < 0)
    goto error_out;

  rc = pipe(err);
  if (rc < 0)
    goto error_err;

  pid = fork();
  if (pid > 0)
  { /* parent */
    // this is the parent process that produces the pipes
    close(in[0]);
    close(out[1]);
    close(err[1]);
    // stdin, write to this
    threepipe[0] = in[1];
    // stdout, read from this
    threepipe[1] = out[0];
    // stderr, read from this
    threepipe[2] = err[0];
    return pid;
  }
  else if (pid == 0)
  { /* child */
    // this is the child process that is replaced by the executed process
    // via execve
    close(in[1]);
    close(out[0]);
    close(err[0]);
    close(0);
    if (!dup(in[0]))
    {
      ;
    }
    close(1);
    if (!dup(out[1]))
    {
      ;
    }
    close(2);
    if (!dup(err[1]))
    {
      ;
    }
    // this replaces the child process with whatever file is executed
    // it returns -1 when there is a failure and on success it does
    // not return
    char *timeout_str = calloc(strlen(command) + 20, sizeof(char));
    sprintf(timeout_str, "timeout 4 %s", command);
    int r = execl("/bin/sh", "sh", "-c", timeout_str, NULL);
    free(timeout_str);
    printf("execl returned: %i\n", r);

    if (r == -1)
    {
      printf("execl error: %s\n", strerror(errno));
    }

    _exit(1);
  }
  else
    goto error_fork;

  return pid;

error_fork:
  close(err[0]);
  close(err[1]);
error_err:
  close(out[0]);
  close(out[1]);
error_out:
  close(in[0]);
  close(in[1]);
error_in:
  return -1;
}

int pcloseTHREE(int pid, int *threepipe)
{
  int status;
  close(threepipe[0]);
  close(threepipe[1]);
  close(threepipe[2]);
  waitpid(pid, &status, 0);
  return status;
}

void execCommand(const char *cmd_string, char **out_string, unsigned long *out_size, char **err_string, unsigned long *err_size)
{
  int first_pipe[3];
  // popenTHREE uses the timeout command
  int pid = popenTHREE(first_pipe, cmd_string);

  printf("Executing Command (pid: %i): %s.\n", pid, cmd_string);

  FILE *f_stdout;
  FILE *f_stderr;

  if (NULL == (f_stdout = fdopen(first_pipe[1], "r")))
  {
    perror("fdopen failed");
  }

  if (NULL == (f_stderr = fdopen(first_pipe[2], "r")))
  {
    perror("fdopen failed");
  }

  *out_string = calloc(PATH_MAX, sizeof(char));
  *err_string = calloc(PATH_MAX, sizeof(char));

  int ch;
  unsigned long c = 0;
  *out_size = PATH_MAX;

  while (1)
  {
    if (c > PATH_MAX)
    {
      printf("PATH_MAX response size reached in command output\n");
      break;
    }

    ch = getc(f_stdout);

    if (ch != EOF)
    {
      if (c > *out_size)
      {
        *out_size += PATH_MAX;
        *out_string = realloc(*out_string, *out_size);
      }

      (*out_string)[c] = ch;
    }
    else
    {
      // done
      break;
    }

    c++;
  }

  // write the stderr to out_stderr
  c = 0;
  *err_size = PATH_MAX;
  while (1)
  {
    if (c > PATH_MAX)
    {
      printf("PATH_MAX response size reached in command output\n");
      break;
    }

    ch = getc(f_stderr);

    if (ch != EOF)
    {
      if (c > *err_size)
      {
        *err_size += PATH_MAX;
        *err_string = realloc(*err_string, *err_size);
      }

      (*err_string)[c] = ch;
    }
    else
    {
      // done
      break;
    }

    c++;
  }

  fclose(f_stdout);
  fclose(f_stderr);
  pcloseTHREE(pid, first_pipe);

  return;
}