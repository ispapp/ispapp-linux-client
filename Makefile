collect-client: collect-client.c
	gcc -g -L /usr/local/lib -I /usr/local/include/mbedtls -I/usr/include/libnl3 -I -I/usr/include/json-c -lmbedtls -lmbedx509 -lmbedcrypto -lpthread -lnl-genl-3 -lnl-3 -lnl-route-3 -ljson-c -o collect-client collect-client.c
