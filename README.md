#gotmux-mem-cpu

----------------------------------
CPU and RAM monitor for use with [tmux](https://tmux.github.io/)
-----------------------------------

Description
===========

Simple CPU and RAM monitor, written on golang. CPU utilization shows in percent over all cores.

```
5375/7859MB ▆ 27.2% ▃ 
 ^    ^     ^   ^   ^
 |    |     |   |   |
 1    2     3   4   5
 ```

1. Currently used memory.
2. Total memory.
3. RAM usage bar graph, <50% - green, >50% - yellow, >75% - red.
4. CPU utilization, default for 2sec.
5. CPU usage bar graph, same as RAM.
 
 
Installation
============

Building
--------
 
* \>= golang 1.0                                                                                                                                                                                               

Download
--------
The link to source code of [monitor](https://github.com/Spirans/gotmux-mem-cpu).

Build & install
-----
```
cd <source dir>
go install gotmux-mem-cpu
```
or
```
cd <source dir>
go build gotmux-mem-cpu.go
sudo cp gotmux-mem-cpu /usr/local/bin
```
or 

download gotmux-mem-cpu for your architecture from the repository and put it in the directory, which is in the $PATH

Configuring
-----------
```
➜  gotmux-mem-cpu -h
Usage of gotmux-mem-cpu:
  -b string
        Background color for cpu and mem status bar (default "black")
  -f string
        Foreground color for cpu and mem status bar (default "white")
  -i int
        Interval in seconds for calculating CPU utilization (default 2)

```

Configuring tmux
--------------
Edit $HOME/.tmux.conf and add for left or right side of status bar:
```
set -g status-interval 2
set -g status-right "#(gotmux-mem-cpu)"
```

_Note: status-interval must be the same like interval for gotmux-mem-cpu (2sec by default)_

Author
------

Veniamin Stepanov <[vics31@gmail.com](mailto:vics31@gmail.com)>
