triggerd
========
`triggerd` executes commands when a GPIO attached button is pressed.
It sits and listens for changes on GPIO pins and executes a command assigned to the pin if an edge has been detected.


Building
--------
Check out the repository and execute
```
go build
```

To cross compile for the Raspberry Pi execute
```
GOOS=linux GOARCH=arm go build
```


Configuration
-------------
The daemon expects one file for each GPIO pin to watch.
By default these files are searched in `etc/triggerd`.
Each file has the following layout:
```
Name = Example Trigger      # Name used for logging and display only

[GPIO]
Pin = 27                    # Pin number in BCM notation
Pull = true                 # Enable pull up/down resistors
Reverse = true              # Reverse the input pin

[Action]
Command = /bin/echo test    # The command to execute
```
