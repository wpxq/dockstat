# dockstat
![GitHub release](https://img.shields.io/github/v/release/wpxq/dockstat?style=for-the-badge&color=blue)
![License](https://img.shields.io/github/license/wpxq/dockstat?style=for-the-badge&color=green)
![GitHub all releases](https://img.shields.io/github/downloads/wpxq/dockstat/total?style=for-the-badge&color=green)

![dockstat](https://github.com/wpxq/dockstat/blob/main/dockstat.png)
![dockstat running](https://github.com/wpxq/dockstat/blob/main/dockstat_running.png)

---

A minimalist, fast, and clean command-line tool for monitoring your running Docker containers. Forget about the clutter of standard `docker ps` outputs—`dockstat` delivers all the essential information you need at a single glance.

---

## Features
- **Container Overview**: ID, name, and image name in clean, readable format.
- **Status Monitoring**: See instantly if your container is running, paused, or stopped.
- **RAM Usage**: Real-time memory consumption in MB.
- **Uptime**: Know exactly how long your container has been alive.
- **Port Mapping**: Quick overview of active container ports

## Requirements
- Go 1.26 or higher (if you build it yourself)
- `github.com/fatih/color`
- `github.com/docker`
## Setup
### Quick Install (Recommended)
**Linux / MacOS**:
```bash
curl -sSL https://raw.githubusercontent.com/wpxq/dockstat/main/dockstat_setup.sh | bash
```
### Manual Install (From Source)
if you prefer to clone and build yourself:
```bash
git clone https://github.com/wpxq/dockstat
cd dockstat
chnod +x dockstat_setup.sh
./dockstat_setup.sh
```

## Troubleshooting (PATH issues)
if the command is not found, add this to your `.bashrc` (or `.zshrc`):
```bash
export PATH="$HOME/.local/bin:$PATH"
```
### Then save and run: `source .bashrc` (or `source .zshrc`)