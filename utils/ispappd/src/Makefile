# build main executable when user executes "make"
CFLAGS:=-lpthread -curl
ispappd: ispappd.o
	$(CC) $(LDFLAGS) ispappd.o -o ispappd
main.o: ispappd.c
	$(CC) $(CFLAGS) -c ispappd.c
# utils.o: utils.c
# 	$(CC) $(CFLAGS) -c utils.c

# remove object files and executable when user executes "make clean"
clean:
	rm *.o ispappd