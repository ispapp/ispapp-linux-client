# build main executable when user executes "make"
CFLAGS:=-lpthread
ispappd: main.o
	$(CC) $(LDFLAGS) main.o -o ispappd
main.o: main.c
	$(CC) $(CFLAGS) -c main.c
# utils.o: utils.c
# 	$(CC) $(CFLAGS) -c utils.c

# remove object files and executable when user executes "make clean"
clean:
	rm *.o ispappd