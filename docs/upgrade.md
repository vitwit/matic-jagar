If you want to upgrade your tool to the latest version please do the following:

Remove any previous repository present on the instance
```
rm -rf matic-jagar
```

Clone the tool repository

```
git clone https://github.com/vitwit/matic-jagar.git
```

Check the releases page of the tool to see which version you want to upgrade to. You can find the releases here: https://github.com/vitwit/matic-jagar/releases

Checkout to that release and build the binary.

```
cd matic-jagar
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
journalctl -u matic-jagar -f
```
