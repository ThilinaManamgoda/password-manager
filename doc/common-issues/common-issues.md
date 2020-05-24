# Common Issues

### Unable to add APT repository
The command to configure Bintray public GPG key,

```
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
```

can be failed if the **gnupg** package is not installed. Please install the **gnupg** package with following command,

```bash
sudo apt install gnupg
```

### apt-get update failed with certificate verification

Command for updating the APT repository,
 ```
 sudo apt update
```
can be failed in the stage of certificate verification.
  
Please install **ca-certificates** package with the following command,

```bash
sudo apt install ca-certificates
```

