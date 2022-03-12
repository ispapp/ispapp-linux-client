.PHONY: all, install, clean

TARGET=collect-client

CSRCS+= collect-client.c utility.c webshell.c

OBJSDIR=./build
OBJS:= $(patsubst %.c, $(OBJSDIR)/%.o, $(CSRCS))

CFLAGS := -I/usr/local/include/mbedtls -I/usr/include/libnl3 -I -I/usr/include/json-c
CFLAGS += -DDEBUG -Wall -g

LDFLAGS := -L/usr/local/lib
LDFLAGS += -lmbedtls -lmbedx509 -lmbedcrypto -lpthread -lnl-genl-3 -lnl-3 -lnl-route-3 -ljson-c

CC:=gcc

all: $(TARGET)
$(TARGET) : $(OBJS)
	@echo " [LINK] $@"
	@mkdir -p $(shell dirname $@)
	@$(CC) $(OBJS) -o $@ $(LDFLAGS)

$(OBJSDIR)/%.o: %.c
	@echo " [CC]  $@"
	@mkdir -p $(shell dirname $@)
	@$(CC) -c $< -o $@ $(CFLAGS)

install:
	cp -rf $(TARGET) /usr/local/bin

clean:
	rm -rf $(OBJSDIR)/*.o
	rm -rf $(TARGET)
