1. `"external": false` 
- Interface is internal to the system
- Not an externally pluggable device

2. `"present": true`
- Physical device exists and is detected by system
- Hardware interface is available

3. `"type": "Network device"`
- Standard network interface type
- Generic network device classification

4. `"up": true`
- Interface is enabled/activated
- Administrative state is up

5. `"carrier": false`
- No physical link detected
- Cable might be unplugged or no physical connection
- Despite interface being "up", there's no active link

6. `"mtu": number` Common MTU values:
- Ethernet: 1500 bytes
- Jumbo frames: 9000 bytes
- PPPoE: 1492 bytes
- IPv6: 1280 bytes minimum

7. `carrier_up_count`: 16
- Number of times the link has come UP
- Cable connected events

8. `carrier_down_count`: 15
- Number of times the link has gone DOWN
- Cable disconnected events

9. `carrier_changes`: 31
- Total number of link state changes
- Should equal `carrier_up_count + carrier_down_count`
- Each connect/disconnect cycle counts as 2 changes
