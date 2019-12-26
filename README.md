# password-manager

A local Password Manager

### Synopsis

A local password manager

### Options

```
  -h, --help                    help for password-manager
  -m, --masterPassword string   Master password
```

### How to install

#### Linux

1. Download the binary from the respective release for Linux.
1. Move the downloaded binary to `/usr/local/bin` with the following command,         
     ```sudo mv password-manager-darwin /usr/local/bin/password-manager```
1. Give executable permission for the binary with the following command,  
    ```sudo chmod +x /usr/local/bin/password-manager```

#### MacOS

1. Download the binary from the respective release for MacOS.
1. Move the downloaded binary to `/usr/local/bin` with the following command,         
     ```sudo mv password-manager-darwin /usr/local/bin/password-manager```
1. Give executable permission for the binary with the following command,  
    ```sudo chmod +x /usr/local/bin/password-manager```
    
### How to use

1. Initialize the password manager with the following command,
    ```$xslt
    password-manager init
    ```
    Enter the Master password once prompted.
    
2. Add a password entry with the following command,
    ```$xslt
    password-manager add TEST -i
    ```
    Here password entry id is **TEST** which always should be a unique value. By passing **-i** parameter we are enabling
    the **Interactive** mode. Once this command is executed, the following entries will be listed,
    
    ```$xslt
        Username: <Enter the username for this password entry>
        Password: <Enter the password for this password entry>
        Enter the Password again: <Enter the password again for this password entry>
        Lables: <Enter lables that can be used for searching. This can be list of comma seperated values>
        Master password: <Enter Master password>
        
        Example:
        
        Username: username@test.com
        Password: ***********
        Enter the Password again: ***********
        Lables: test,first,firstTime
        Master password: <Enter Master password>
    ```

1. Search a password entry with **ID**
    ```$xslt
        password-manager get <ID> 
    ```     
    Enter the Master password once prompted. For an example,
    ```$xslt
        password-manager get TEST
    ```
    Once you enter the Master password the password entry will be listed that match the given id. Once the entry is selected, the **password** will be copied to the clipboard. 
1. Search a password entry with a **label**
    ```$xslt
        password-manager search-label <LABEL>
    ```
    Enter the Master password once prompted. For an example,
    ```$xslt
        password-manager search-label test
    ```
    Once you enter the Master password the a list of password entries will be listed that match the given label. Once the entry is selected, the **password** will be copied to the clipboard.
     
### SEE ALSO

* [password-manager init](doc/password-manager_init.md)	 - Initialize the Password Manager
* [password-manager add](doc/password-manager_add.md)	 - Add a new password
* [password-manager change](doc/password-manager_change.md)	 - Change a password entry
* [password-manager change-master-password](doc/password-manager_change-master-password.md)	 - Change Master password
* [password-manager generate-password](doc/password-manager_generate-password.md)	 - Generate a secure password
* [password-manager get](doc/password-manager_get.md)	 - Get a password
* [password-manager search-label](doc/password-manager_search-label.md)	 - Search Password with Label
* [password-manager search-id](doc/password-manager_search-id.md)	 - Search Password with ID
* [password-manager import](doc/password-manager_import.md)	 - Import passwords
* [password-manager export](doc/password-manager_export.md)	 - Export password repository to a file
* [password-manager remove](doc/password-manager_remove.md)	 - Remove a password



