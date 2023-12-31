# Docker tui managment
This is a command line app that connect to the docker local daemon in order to execute some actions.

This is and alternative if you want to stay in the console using something more nice that the default cli docker and do not want to use a web UI like portainer.

It uses the great lib bubbletea for create a nice UI in the terminal

https://github.com/charmbracelet/bubbletea

## Demo
![demo](https://github.com/ernesto27/dcli/assets/1366157/7ca15b20-da6f-4541-9691-46d0a627ffef)


## Installation
Go to release page and download the binary for your OS.

[https://github.com/ernesto27/container-nice-cli/releases](https://github.com/ernesto27/dcli/releases)

Also if you have Golang installed you can install it with:
```
go install github.com/ernesto27/dcli@latest
```


## Key bindings
| Key              | Description                                 |
|:-----------------|:--------------------------------------------|
| <kbd>ctrl+f</kbd>     | Search containers by name              |
| <kbd>ctrl+l</kbd>     | View logs containers                 |
| <kbd>ctrl+o</kbd>     | Options for container (stop, start, remove)|
| <kbd>ctrl+e</kbd>     | Exec in a contaner                    |
| <kbd>ctrl+b</kbd>     | List images
| <kbd>ctrl+f</kbd>     | On image list, search by image name    |
| <kbd>ctrl+o</kbd>     | Options image    |
| <kbd>ctrl+n</kbd>     | Network list    |
| <kbd>ctrl+f</kbd>     | Search network by name    |
| <kbd>ctrl+o</kbd>     | Option network    |
| <kbd>ctrl+v</kbd>     | Volume list    |
| <kbd>ctrl+f</kbd>     | Search volume by name    |
| <kbd>ctrl+o</kbd>     | Option volume    |







