If you want to upgrade your tool to the latest version please do the following:


Check the releases page of the tool to see which version you want to upgrade to. You can find the releases here: https://github.com/vitwit/matic-jagar/releases

Checkout to that release and build the binary.

```
cd matic-jagar
git fetch
git checkout <version-tag>
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
