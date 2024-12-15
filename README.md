# **Chicha TCP Proxy: A Simple and Fast Port Forwarding Tool**

Chicha TCP Proxy is a lightweight **Layer 2 (L2) proxy** designed to efficiently forward traffic between local and remote ports. It’s faster and more optimized than tools like `xinetd`, making it a perfect choice for high-performance setups. 

---

### **Why Use Chicha TCP Proxy?**
- **Simple Port Forwarding**: Easily map local and remote ports with `LOCALPORT:REMOTEIP:REMOTEPORT`.
- **High Performance**: Fully utilizes all CPU cores for maximum efficiency.
- **Minimal Configuration**: One command to start—no complex setup required.
- **Automatic Log Management**: Daily log rotation with `.gz` compression.
- **Cross-Platform**: Compatible with all major operating systems and architectures.

---

### **Quick Start Example**
Forward local port `80` to `192.168.0.1:80` with a single command:
```bash
chicha-tcp-proxy -routes "80:192.168.0.1:80" -log /var/log/chicha-tcp-proxy.log
```

---

### **Key Features**
- **Super Fast**: Optimized for handling high traffic loads without slowing down.
- **Logs Included**: Real-time logging to both console and file for easy monitoring.
- **Autostart on Linux**: Comes with a `systemd` service file for startup on boot.
- **Simple Installation**: One binary download, ready to run.

---

### **How Is It Different?**
Unlike `xinetd` and similar tools, Chicha TCP Proxy focuses on:
- **Speed**: Designed for high-performance port forwarding.
- **Simplicity**: No complex configuration files—just one simple command.
- **Efficiency**: Lightweight and resource-friendly.

---

### **Systemd Autostart Setup (Linux)**
1. Create a `chicha-tcp-proxy.service` file using `mcedit`:
   ```bash
   sudo mcedit /etc/systemd/system/chicha-tcp-proxy.service
   ```
2. Add the following content:
   ```ini
   [Unit]
   Description=Chicha 80 and 443 TCP Proxy
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/chicha-tcp-proxy -routes "80:192.168.0.1:80,443:192.168.0.1:443," -log /var/log/chicha-tcp-proxy.log
   Restart=on-failure

   [Install]
   WantedBy=multi-user.target
   ```
3. Save and exit `mcedit`.

4. Enable and start the service:
   ```bash
   sudo systemctl enable chicha-tcp-proxy
   sudo systemctl start chicha-tcp-proxy
   ```

---

Chicha TCP Proxy offers **speed and simplicity** for system administrators who need reliable port forwarding without the hassle. Try it today and see the difference!
