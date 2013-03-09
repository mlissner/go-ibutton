go-ibutton
==========

Go application to use Maxim iButtons

# installation

Go package management can install ibutton directly from github:
```
go install github.com/maxhille/go-ibutton/ibutton
```

# usage

start a new mission
```
ibutton -command start
```

stop the currently running mission
```
ibutton -command stop
```

print out the sample log
```
ibutton -command read
```

show the button status
```
ibutton -command status
```

clear the button mission memory
```
ibutton -command clear
```
