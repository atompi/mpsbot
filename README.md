# mpsbot

[metrics-post-station](https://github.com/atompi/metrics-post-station) bot

Scan keys and generate Prometheus scrape configs.

A tool that iterates over Redis keys and generates Prometheus scrape configuration files.

## Build

```
go build -ldflags '-s -w' -o examples/mpsbot
```

## Usage

```
mkdir /app/mpsbot
cp examples/mpsbot /app/mpsbot
chmod +x /app/mpsbot/mpsbot

cat > /app/mpsbot/mpsbot.yaml <<EOF
core:
  log:
    level: info
    path: ./logs/mpsbot.log
    maxsize: 100
    maxage: 7
    compress: true

task:
  interval: 15
  outputpath: ./

redis:
  addr: 127.0.0.1:6379
  password: 123456
  db: 0
  dialtimeout: 5
  prefix: mps-
EOF
cp init/systemd/mpsbot.service /lib/systemd/system/
systemctl daemon-reload
systemctl enable --now mpsbot
```
