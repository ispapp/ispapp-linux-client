# openwrt-ispapp IPK

![ISPApp Logo](/img/logo.png)

Learn more at https://ispapp.co

> ~usage:

```
    /etc/init.d/ispappd start
```
```
    /etc/init.d/ispappd stop
```

# about


## translation
```
xgettext --language=Lua --keyword=translate --output=po/ispapp.pot ./luasrc/controller/ispapp.lua ./luasrc/model/cbi/ispapp/overview.lua ./luasrc/model/cbi/ispapp/settings.lua ./luasrc/model/cbi/ispapp/logread.lua
```
This is an ISPApp client designed to monitor hosts running Linux.

ISPApp allows you to monitor, configure and command hosts quickly and easily with high resolution charts and realtime data.

There are realtime, daily, weekly, monthly and yearly charts for:

* All Network Interfaces - Traffic and Packet Rate
![Traffic](/img/if-traffic.png)
* All Wireless Interfaces - RSSI and Traffic per Connected Station, # of Stations per Interface
![RSSI](/img/rssi.png)
* All System Disks - Total, Used and Available Disk Space
![Disk](/img/disk.png)
* System Load
![Load](/img/load.png)
* System Memory
![Memory](/img/memory.png)
* Ping to a Host - Average, Maximum and Minimum RTT + Total Loss
![Ping](/img/ping.png)
* Environment - Temperature, Humidity, Precipitation, Barometric Pressure and others
* Industrial Sensors - Gas Levels, Pressure, Fill Levels, Rotation and Positioning Data
* Vehicle Data - Torque, Fuel Rate and others
* Electronic Data - Voltage, Power and Current
* Request Data - HTTP Request Rate, DNS Request Rate, Rate of Incoming Emails etc (easily added to software with our REST API)

ISPApp also provides outage notifications and maintenance/degradation analysis..

# dependencies

* json-c
* libnl3
* mbedtls
* timeout

This has been tested running on small devices with 4MB of RAM and 16MB of disk space.

# build

```
# install deps as root
yum -y install json-c-devel libnl3-devel timeout

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
git clone https://github.com/ispapp/ispapp-linux-client
cd ispapp-linux-client
make
```

# run

The parameter field names are shown by launching the program without any options.

It must be run as root to send ping packets, ping requires a raw network socket which is supposed to be represented on the network as a privileged action.

```
cd ispapp-linux-client
sudo LD_LIBRARY_PATH=/usr/local/lib ./collect-client subdomain.ispapp.co 8550 eth0 "long_host_key" "amazon" "ec2" "amazon linux 2" "nano" "NA" "1602685864" "" ./ /tmp/collect-client-config.json > /dev/null 2>&1 &

# you can place the launch command above into a startup script like /etc/rc.local
```

# license

The project ispapp-linux-client is licensed per the GNU General Public License, version 2

A copy is in the project directory, as a file named LICENSE

# roadmap

We plan to create packages for various distributions targeting:

* Servers
* Virtualized Instances
* IoT Devices
* Desktops/Laptops
* Routers
