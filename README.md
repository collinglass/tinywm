# TinyWM

Port of TinyWM in Go.

## Usage

Add the ```.desktop``` file to the xsessions directory. This will add a new selection to the list of window managers at the login screen. 

```
$ sudo cp tinywm.desktop /usr/share/xsessions/
```

This file runs /usr/bin/tinywm-session. Create that file now.

```
$ sudo cp tinywm-session /usr/bin/
$ sudo chmod a+x /usr/bin/tinywm-session
```

tinwm-session preloads a terminal and then runs the go binary. Build the binary and move it to /usr/bin/. On Ubuntu, `go build` may throw `fatal error: X11/Xlib.h: No such file or directory`. Do `sudo apt-get install libx11-dev` to add Xlib.h .

```
$ go build
$ sudo mv tinywm /usr/bin/
$ sudo chmod a+x /usr/bin/tinywm
```

Log out. You should now see tinywm as an option in the list. Choose it a log in. You should have a bare WM with a terminal emulator running.

## Development

Issues and PR's are welcome! :) Feel free to hit me on twitter @collinglass.
