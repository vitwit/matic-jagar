If you want to upgrade your tool to the latest version please do the following:


Pull the latest code and build the binary.

```
cd matic-jagar
git fetch
git pull origin main
go build -o matic-jagar
```

Move the binary to your $GOBIN and restart the system service.

```
mv matic-jagar ~/go/bin
sudo systemctl restart matic_jagar.service
```

You can check the logs using:
```
journalctl -u matic_jagar -f
```
