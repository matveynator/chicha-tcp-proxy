# **Chicha TCP Proxy**

## **Overview**
Chicha TCP Proxy is a lightweight and efficient TCP proxy tool. It forwards traffic between specified local and remote ports and supports log rotation with compression. It is available for a wide range of operating systems and architectures.

---

## **Features**
- **Simple TCP Forwarding**: Define routes as `LOCALPORT:REMOTEIP:REMOTEPORT`.
- **High Performance**: Utilizes all available CPU cores.
- **Log Rotation**: Logs are rotated daily and compressed into `.gz` format.
- **Simultaneous Logging**: Outputs logs to both console and file.
- **Systemd Compatibility**: Includes a service file for easy autostart on Linux.

---

## **Download**

Chicha TCP Proxy binaries are available for multiple platforms and architectures.  
You can browse and download the latest versions here:  
[http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/)

### **Available Platforms**
- [AIX](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/aix/)
- [Android](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/android/)
- [DragonFly BSD](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/dragonfly/)
- [FreeBSD](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/freebsd/)
- [Illumos](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/illumos/)
- [iOS](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/ios/)
- [JavaScript/WebAssembly](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/js/)
- [Linux](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/linux/)
- [macOS](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/mac/)
- [NetBSD](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/netbsd/)
- [OpenBSD](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/openbsd/)
- [Plan 9](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/plan9/)
- [Solaris](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/solaris/)
- [WASI/WebAssembly](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/wasip1/)
- [Windows](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/windows/)
- [z/OS](http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/zos/)

---

### **Quick Install (Linux x86_64)**

For Linux (64-bit), use this command to download and install:
```bash
sudo curl -o /usr/local/bin/chicha-tcp-proxy http://files.zabiyaka.net/chicha-tcp-proxy/latest/no-gui/linux/amd64/chicha-tcp-proxy && sudo chmod +x /usr/local/bin/chicha-tcp-proxy
```

---

## **Usage**

### **Basic Command**
Run the proxy using:
```bash
chicha-tcp-proxy -routes "LOCALPORT:REMOTEIP:REMOTEPORT,..." -log /path/to/logfile.log
```

#### **Example**:
```bash
chicha-tcp-proxy -routes "80:192.168.0.1:80,443:192.168.0.1:443" -log /var/log/chicha-tcp-proxy.log
```

#### **Flags**:
| Flag          | Description                                                                              | Default                  |
|---------------|------------------------------------------------------------------------------------------|--------------------------|
| `-routes`     | Comma-separated list of routes in the format `LOCALPORT:REMOTEIP:REMOTEPORT`.            | Required                 |
| `-log`        | Path to the log file where logs will be written.                                         | `chicha-tcp-proxy.log`   |
| `-rotation`   | Frequency for log rotation, e.g., `24h` (24 hours), `1h` (1 hour).                       | `24h`                    |

---

## **Systemd Setup**

To ensure the proxy runs automatically on system startup, follow these steps:

1. **Create Service File**:
   ```bash
   sudo nano /etc/systemd/system/chicha-tcp-proxy.service
   ```

2. **Add the Following Content**:
   ```ini
   [Unit]
   Description=Chicha TCP Proxy Service
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/chicha-tcp-proxy -routes "80:192.168.0.1:80,443:192.168.0.1:443" -log /var/log/chicha-tcp-proxy.log
   Restart=on-failure
   RestartSec=5s
   User=root
   Group=root

   [Install]
   WantedBy=multi-user.target
   ```

3. **Enable and Start the Service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable chicha-tcp-proxy
   sudo systemctl start chicha-tcp-proxy
   ```

4. **Check Service Status**:
   ```bash
   sudo systemctl status chicha-tcp-proxy
   ```

---

## **Logs**

- Logs are saved to the file specified by the `-log` flag (e.g., `/var/log/chicha-tcp-proxy.log`).
- Logs rotate daily and are compressed into `.gz` format.
- Both console and file logging are enabled simultaneously.

---

## **Example Output**

**On Startup**:
```plaintext
========== CHICHA TCP PROXY ==========
Routes:
  LocalPort=80 -> RemoteIP=192.168.0.1 RemotePort=80
  LocalPort=443 -> RemoteIP=192.168.0.1 RemotePort=443
Log file: /var/log/chicha-tcp-proxy.log
Log rotation frequency: 24h0m0s
======================================
```

**Logs**:
```plaintext
2024/12/14 13:30:01 Starting proxy for route: local=80 remote=192.168.0.1:80
2024/12/14 13:30:02 Log file rotated successfully, now compressing old log...
2024/12/14 13:30:03 Compression successful: /var/log/chicha-tcp-proxy.log.2024-12-14.gz
```

---

## **Uninstallation**

To remove `chicha-tcp-proxy` from your system:

1. Stop and disable the systemd service:
   ```bash
   sudo systemctl stop chicha-tcp-proxy
   sudo systemctl disable chicha-tcp-proxy
   ```

2. Delete the binary:
   ```bash
   sudo rm /usr/local/bin/chicha-tcp-proxy
   ```

3. Remove the systemd service file:
   ```bash
   sudo rm /etc/systemd/system/chicha-tcp-proxy.service
   ```

4. Reload systemd:
   ```bash
   sudo systemctl daemon-reload
   ```

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.

--- 

Let me know if you need further changes!
