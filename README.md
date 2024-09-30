# Caddy GUI

![](https://img.shields.io/badge/Golang-1.22-blue)
![](https://img.shields.io/badge/Fyne-v2-blue)
![](https://img.shields.io/badge/PRs-welcome-green)

<p align="center">
  <img src="https://github.com/user-attachments/assets/09d70ea6-2c53-406f-889f-c6e05a919e99" width="45%"/>
</p>

Easily configure and start the Caddy server with a user-friendly interface. This tool allows you to quickly set up and manage your Caddy server, making it ideal for creating test servers and managing configurations effortlessly.

## Start

#### 1. Download `caddy` at [caddyserver.com](https://caddyserver.com/download)

```bash
mv caddy-[your-platform] caddy
```

#### 2. Build

```bash
fyne package -os [your-platform] -icon icon.png
```

#### 3. Run

##### 3.1. MacOS

Move the `caddy` binary to the `Caddy GUI.app` directory

```bash
cp caddy Caddy\ GUI.app/Contents/MacOS
```

Then, move `Caddy GUI.app` to the Applications directory

##### 3.2. Windows

Ensure that `caddy.exe` is located in the same directory as the CaddyGUI executable.

```
.
├── caddy.exe
└── caddy-gui.exe
```
