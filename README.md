![ISPApp Logo]
(/logo.png)

# about

This is an ISPApp client which is designed to monitor hosts running Linux.

It will automatically monitor and allow ISPApp to generate realtime, daily, weekly, monthly and annual charts for:

* All Network Interfaces - Traffic and Packet Rate
* All Wireless Interfaces - RSSI and Traffic per Connected Station, # of Stations per Interface
* All System Disks - Total, Used and Available Disk Space
* System Load
* System Memory
* Ping to a Host - Average, Maximum and Minimum RTT + Total Loss

# build

```
# install deps as root
yum -y install json-c-devel libnl3-devel

# add to /etc/profile
export SHARED=1
export LD_LIBRARY_PATH=/usr/local/lib

# build and install mbedtls
cd
wget https://github.com/ARMmbed/mbedtls/archive/v2.24.0.tar.gz
tar -xzvf v2.24.0.tar.gz
cd mbedtls-2.24.0
make
sudo --preserve-env make install

# get ispapp-linux-client source
cd
git clone https://github.com/andrewhodel/ispapp-linux-client
cd ispapp-linux-client
make
```

# launch

The launch parameter field names are shown by launching the program without any options.

```
LD_LIBRARY_PATH=/usr/local/lib ./collect-client subdomain.ispapp.co 8550 eth0 "long_host_key" "amazon" "ec2" "amazon linux 2" "nano" "NA" "1602685864" "" "" /home/ec2-user/ispapp-keys/__ispapp_co.ca-bundle /tmp/collect-client-config.json > /dev/null 2>&1 &

# you can place the launch command above into a startup script like /etc/rc.local
```

# license

The project ispapp-linux-client is licensed per the GNU General Public License, version 2

A copy is in the project directory, as a file named LICENSE
