# ObscuraScan (OBScan)

ObscuraScan (**OBScan**) is a fast, lightweight network reconnaissance tool designed for port scanning and service identification. With added flexibility for handling hostnames, banner grabbing, and customizable protocols, OBScan is a powerful tool for security professionals, developers, or anyone exploring the depths of a network.

---

## ✨ Features
- **Protocol support**: Choose between TCP or other supported protocols for scanning.  
- **Hostname and IP scanning**: Seamlessly scan using IP addresses or domain names with automatic DNS resolution.  
- **Port range customization**: Specify single ports, multiple ports, or ranges (e.g., `80`, `20-100`).  
- **Banner grabbing**: Identify running services behind open ports.
- **Vulnerability Detection**: Identify known vulnerabilities base on the banner we grab ( using the NIST api ) 
- **Timeout settings**: Adjust connection timeouts to suit your needs.  
- **Multi-host scanning**: Scan multiple resolved IPs for a single hostname.  

---

## 🚀 Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/obscurascan.git
   ```
2. Navigate to the directory:
   ```bash
   cd obscurascan
   ```
3. Build the binary:
   ```bash
   go build -o obscan
   ```

---

## 🔮 Usage
Run OBScan with the following options:

```bash
./obscan -type <protocol> -ip <target-host> -p <port-range> [-t <timeout>] [-g]
```

### Example
#### Scan ports 20-80 on `192.168.1.1` with TCP:
```bash
./obscan -type tcp -ip 192.168.1.1 -p 20-80
```

#### Scan a hostname, resolve IPs, and grab banners:
```bash
./obscan -type tcp -ip example.com -p 443 -g
```

### Options
- `-type`: Protocol to use for scanning (default: `tcp`).  
- `-ip`: Target host (IP address or domain name).  
- `-p`: Port or range of ports to scan (e.g., `80,443` or `20-100`).  
- `-t`: Timeout for connections in seconds (default: `5`).  
- `-g`: Enable banner grabbing for open ports.

---

## 🌌 Roadmap
- [x] Basic TCP port scanning.  
- [x] Hostname resolution to IP addresses.  
- [x] Banner grabbing for open ports.
- [ ] handle UDP scanning
- [ ] enhance vulnerability detection   
- [ ] Export scan results to JSON or CSV.  
- [ ] Protocol-based scan optimization.  
- [ ] Enhanced output formatting.  
- [ ] IPv6 support.

---

## 📜 Contribution
We welcome contributions! Submit issues, suggest new features, or open pull requests to make ObscuraScan even better.

---

## ⚖️ License
ObscuraScan is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

## 🌟 Disclaimer
ObscuraScan is intended for ethical and educational purposes only. Unauthorized use on networks without explicit permission is strictly prohibited.
