# WARNING!!! THIS IS PROTOTYPE

# Capuchin

## Build comparison:

-   `capuchin_windows_tray_amd64.exe` - desktop edition with tray and auto-opening the browser when launching
-   `capuchin_windows_amd64.exe` - server edition as a service
-   `capuchin_linux_amd64` - server edition as a service

## Install

```bash
sudo -i
mkdir capuchin
cd capuchin
wget -O capuchin "https://github.com/capuchinapp/capuchin/releases/download/v1.0.0/capuchin_linux_amd64"
chmod 755 capuchin
wget -P /lib/systemd/user "https://github.com/capuchinapp/capuchin/releases/download/v1.0.0/capuchin.service"
chmod 755 /lib/systemd/user/capuchin.service
systemctl enable /lib/systemd/user/capuchin.service
```

## Run

```bash
service capuchin start
service capuchin status
```

## Update

```bash
sudo -i
service capuchin stop
service capuchin status
cd capuchin
wget -O capuchin "https://github.com/capuchinapp/capuchin/releases/download/v1.0.0/capuchin_linux_amd64"
chmod 755 capuchin
systemctl daemon-reload
service capuchin start
service capuchin status
```

## Development

Default ports:

-   `5173` - frontend
-   `8090` - backend
