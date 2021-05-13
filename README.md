![ISPApp Logo](/img/logo.png)

Learn more at https://ispapp.co

Watch a YouTube Video about how ISPApp can help you - https://www.youtube.com/watch?v=BQN8FdMqApo

# about

This is an ISPApp client which is designed to monitor hosts running Linux.

ISPApp allows you to monitor thousands of hosts or IoT devices quickly and easily with high resolution charts and realtime data.

It will automatically monitor a host when `collect-client` is ran on that host and send ISPApp data to generate realtime, daily, weekly, monthly and annual charts for:

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

ISPApp also provides outage notifications and maintenance/degradation analysis for each of the monitored data types.

We have ISPApp Instances running with tens of thousands of charts and are ready for you to be a customer.

# dependencies

* json-c
* libnl3
* mbedtls

This has been tested running on devices with 4MB of RAM and 16MB of Disk Space.

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
git clone https://github.com/ispapp/ispapp-linux-client
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

# roadmap

We plan to create packages for various distributions targetting:

* Servers
* Virtualized Instances
* IoT Devices
* Desktops/Laptops
* Raspberry Pi/Etc
* Routers (we already have an OpenWRT package, email me for more information)

If you would like to be a package maintainer, please email me via andrew@ispapp.co
