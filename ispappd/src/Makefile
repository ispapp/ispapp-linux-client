# Compiler
CC = gcc

# Compiler Flags
CFLAGS = -Wall -g

# Include the directory where curl/curl.h is located
CFLAGS += -I/usr/include

# Linker Flags
LDFLAGS = -lcurl

# Source files
SRCS = ispappd.c

# Output executable
TARGET = ispappd

# Build target
$(TARGET): $(SRCS)
	$(CC) $(CFLAGS) $(SRCS) -o $(TARGET) $(LDFLAGS)

# Clean target
clean:
	rm -f $(TARGET)
