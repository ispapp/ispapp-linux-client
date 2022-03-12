#include "utility.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

char *escape_string_for_json(char *str) {
  // allocate the length of str
  char *nstr = calloc(strlen(str) + 1, sizeof(char));

  // loop through each character
  long unsigned int c = 0;
  long unsigned int d = 0;
  while (c < strlen(str)) {
    // printf("character: %c\n", str[c]);

    // json needs everything from '\x00' to '\x1f' escaped
    if (str[c] == '"' || str[c] == '\\' || ('\x00' <= str[c] && str[c] <= '\x1f')) {
      // printf("\tescaping %c\n", str[c]);

      // add the escape character to nstr
      nstr[d] = '\\';

      // increment d to account for the extra space
      d++;

      // allocate that space in the nstr pointer
      nstr = realloc(nstr, d);

      // add the character
      nstr[d] = str[c];

    } else {
      // add the character to nstr
      nstr[d] = str[c];
    }

    c++;
    d++;
  }

  // add the \0 at the end
  nstr[d] = '\0';

  return nstr;
}

int l_strcpy(char *dest, char *src, int start, int end) {
  // returns number of characters copied

  int c = 0;
  int d = 0;
  while (c < strlen(src)) {
    if (c > end && end > 0) {
      break;
    }

    if (c > start) {
      dest[d] = src[c];
      d++;
    }

    c++;
  }

  dest[d] = '\0';

  return d;
}