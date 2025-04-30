iperf -c speedtest.is.cc -p 5201-5209 -i 1 -t 1 -P 10 -f m | awk '/SUM/ {if (down == "") down=$7; else up=$7} END {printf "{ \"down\": %.2f, \"up\": %.2f }\n", down, up}'

iperf3 -c 192.168.200.169 -i 1 -t 1 -P 10 -f m | awk '/SUM.*sender/ {up=$6} /SUM.*receiver/ {down=$6} END {print "{ \"down\": " down ", \"up\": " up " }"}'
iperf3 -c speedtest.phx1.us.leaseweb.net -p 5201-5210 -i 1 -t 3 -P 10 -f m