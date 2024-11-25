# ObscuraScan (OBScan)

ObscuraScan (or **OBScan** for short) is a fast, lightweight port scanner designed for network reconnaissance and service discovery. OBScan not only scans open ports on a target host but also retrieves service banners, unveiling the secrets of the network.

Whether you're a security professional, developer, or just exploring the unknown, OBScan is your go-to tool for uncovering what lies behind the veil of the network.

---

## ✨ Features
- **Fast and concurrent scanning**: Harness the power of Go's goroutines to scan multiple ports simultaneously.  
- **Banner grabbing**: Peek behind open ports and identify running services.  
- **Customizable scanning**: Define port ranges and timeout values to suit your reconnaissance needs.  
- **Portable binary**: Cross-platform compatibility with no external dependencies.  

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
./obscan -host <target-host> -ports <port-range> [-timeout <timeout>]
```

### Example
Scan ports 1-100 on `192.168.1.1`:
```bash
./obscan -host 192.168.1.1 -ports 1-100
```

### Options
- `-host`: Target host (IP address or domain name).  
- `-ports`: Port range to scan (e.g., `1-1000`).  
- `-timeout`: Timeout for connections in seconds (default: 5 seconds).  

---

## 🌌 Roadmap
- [x] Basic TCP port scanning.  
- [x] Concurrent scanning with goroutines.  
- [x] Banner grabbing for open ports.  
- [ ] Export scan results to JSON or CSV.  
- [ ] Advanced service detection based on banners.  
- [ ] Add optional stealth mode with random delays.  

---

## 📜 Contribution
We welcome contributions from fellow explorers! Feel free to submit issues, suggest new features, or open pull requests to improve ObscuraScan.

---

## ⚖️ License
ObscuraScan is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

## 🌟 Disclaimer
ObscuraScan is intended for educational and ethical purposes only. Unauthorized use on networks without explicit permission is prohibited.
