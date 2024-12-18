<img src="https://github.com/matveynator/chicha-tcp-proxy/blob/master/chicha-tcp-proxy.png?raw=true" width="100%" align="right"></img>


# Chicha TCP Proxy is the fastest and simplest solution for port forwarding.

Chicha TCP Proxy is a lightweight **Layer 2 (L2) proxy** written in **Go**, designed for efficient port forwarding. It’s faster and simpler than traditional tools like `xinetd`, making it ideal for handling high traffic with minimal setup.

---

### **Why Choose Chicha TCP Proxy?**
- **Blazing Fast**: Written in Go, it fully utilizes all CPU cores for optimal performance.
- **Simple to Use**: One command to forward traffic—no complex configurations.
- **Log Rotation**: Automatically compresses logs into `.gz` format daily.
- **Cross-Platform**: Compatible with all major operating systems and architectures.
- **Efficient**: Low resource usage even under heavy load.

---

### **Download and Install**
1. **Download the Binary**:  
   Choose your platform and architecture from the links below:
   - [Linux (x86_64)](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/linux/amd64/chicha-tcp-proxy)
   - [All Platforms](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/)

2. **Install the Binary**:  
   For Linux (x86_64), download and install with:
   ```bash
   sudo curl -o /usr/local/bin/chicha-tcp-proxy http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/linux/amd64/chicha-tcp-proxy && sudo chmod +x /usr/local/bin/chicha-tcp-proxy
   ```

---

### **Quick Start Examples**

#### **Single Port Example**
Forward traffic from local port `80` to remote `192.168.0.1:80`:
```bash
sudo chicha-tcp-proxy -routes "80:192.168.0.1:80" -log /var/log/chicha-tcp-proxy.log
```

#### **Multiple Ports Example**
Forward multiple ports simultaneously:
```bash
sudo chicha-tcp-proxy -routes "80:192.168.0.1:80,443:192.168.0.1:443,22:192.168.0.2:22" -log /var/log/chicha-tcp-proxy.log
```

In this example:
- Port `80` on the local machine forwards to `192.168.0.1:80`.
- Port `443` forwards to `192.168.0.1:443`.
- Port `22` forwards to `192.168.0.2:22`.

---

### **Systemd Autostart Setup**

1. **Create a Service File**:
   ```bash
   sudo mcedit /etc/systemd/system/chicha-tcp-proxy.service
   ```

2. **Add the Following Content**:
   ```ini
   [Unit]
   Description=Chicha TCP Proxy
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/chicha-tcp-proxy -routes "80:192.168.0.1:80,443:192.168.0.1:443" -log /var/log/chicha-tcp-proxy.log
   Restart=on-failure

   [Install]
   WantedBy=multi-user.target
   ```

3. Save and exit `mcedit`.

4. **Enable and Start the Service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable chicha-tcp-proxy
   sudo systemctl start chicha-tcp-proxy
   ```

---

### **Why Chicha TCP Proxy?**
- **Go-Powered Performance**: Written in Go, ensuring speed and reliability.
- **Multi-Port Support**: Easily forward traffic for one or multiple ports.
- **No Complexity**: Simple commands, no bloated configs.
- **Ready for Production**: Log rotation, compression, and systemd integration make it production-ready.

---
### **Benchmarks**

#### **chicha-tcp-proxy:**

```
siege http://localhost:8081 -t 15s -c 100
** SIEGE 4.0.4
** Preparing 100 concurrent users for battle.
The server is now under siege...
Lifting the server siege...
Transactions:		       15728 hits
Availability:		      100.00 %
Elapsed time:		       14.65 secs
Data transferred:	      124.07 MB
Response time:		        0.03 secs
Transaction rate:	     1073.58 trans/sec
Throughput:		        8.47 MB/sec
Concurrency:		       33.37
Successful transactions:       11824
Failed transactions:	           0
Longest transaction:	        0.25
Shortest transaction:	        0.00
```

#### **xinetd:**
```
siege http://localhost:8082 -t 15s -c 100
** SIEGE 4.0.4
** Preparing 100 concurrent users for battle.
The server is now under siege...
Lifting the server siege...
Transactions:		       14863 hits
Availability:		      100.00 %
Elapsed time:		       14.57 secs
Data transferred:	      117.20 MB
Response time:		        0.04 secs
Transaction rate:	     1020.11 trans/sec
Throughput:		        8.04 MB/sec
Concurrency:		       36.99
Successful transactions:       11178
Failed transactions:	           0
Longest transaction:	        0.55
Shortest transaction:	        0.00
 
```

#### **direct requests:**
```
siege http://files.zabiyaka.net:80 -t 15s -c 100
** SIEGE 4.0.4
** Preparing 100 concurrent users for battle.
The server is now under siege...
Lifting the server siege...
Transactions:		       14778 hits
Availability:		      100.00 %
Elapsed time:		       14.14 secs
Data transferred:	      116.52 MB
Response time:		        0.03 secs
Transaction rate:	     1045.12 trans/sec
Throughput:		        8.24 MB/sec
Concurrency:		       35.07
Successful transactions:       11112
Failed transactions:	           0
Longest transaction:	        0.19
Shortest transaction:	        0.00
```

---


Chicha TCP Proxy is the **fastest and simplest solution** for port forwarding. Whether forwarding one port or dozens, it's the ideal tool for sysadmins looking for performance and ease of use!
