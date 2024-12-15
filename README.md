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
chicha-tcp-proxy -routes "80:192.168.0.1:80" -log /var/log/chicha-tcp-proxy.log
```

#### **Multiple Ports Example**
Forward multiple ports simultaneously:
```bash
chicha-tcp-proxy -routes "80:192.168.0.1:80,443:192.168.0.1:443,22:192.168.0.2:22" -log /var/log/chicha-tcp-proxy.log
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

Chicha TCP Proxy is the **fastest and simplest solution** for port forwarding. Whether forwarding one port or dozens, it's the ideal tool for sysadmins looking for performance and ease of use!
