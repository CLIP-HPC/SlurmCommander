
# SlurmCommander

## Description

SlurmCommander is a simple, lightweight, no-dependencies text-based user interface (TUI) to your cluster.
It ties together multiple slurm commands to provide you with a simple and efficient interaction point with slurm.

Installation does not require any special privileges or environment. Simply download the [binary](https://github.com/pja237/SlurmCommander-dev/releases/latest), fill out a small [config file](./cmd/scom/scom.conf) and it's ready to run.

You can view, search, analyze and interact with:

* Job queue - view, search, inspect, ssh to nodes, issue commands to jobs etc.
* Job history - browse, search and inspect past jobs
* Edit and submit jobs from predefined templates (can be your own or group/site ones)
  * In the config file, set the `templatedirs` list of directories where to look for _.sbatch_ templates and their _.desc_ description files
* Examine state of cluster nodes and partitions

Example Job Queue tab demo:
![demo](./images/jobqueue.gif)

## Installation

SlurmCommander does not require any special privileges to be installed, see instructions below.

> Hard requirement: json-output capable slurm commands

### Regular users

1. Download the pre-built [binary](https://github.com/pja237/SlurmCommander-dev/releases/latest)
2. Download the [annotated config](./cmd/scom/scom.conf) file
3. Edit the config file, follow instructions inside
4. Create scom directory in your $HOME and place the edited config there: `mkdir $HOME/scom`
5. Run

### Site administrators

Instructions are same as for the regular users, with one minor perk. 
Place the [config file](./cmd/scom/scom.conf) in `/etc/scom/scom.conf` to be used as global configuration source for all scom instances on that machine.

__NOTE__: Users can still override global configuration options by changing config stanzas in their local `$HOME/scom/scom.conf`

## Usage tips

SlurmCommander is developed for 256 color terminals (black background) and requires at least 185x43 (columns x rows) to work.

* If you experience _funky_ colors on startup, try setting your `TERM` environment variable to something like `xterm-256color`.
* If you get a message like this:
`FATAL: Window too small to run without breaking view. Have 80x24. Need at least 185x43.`, check your terminal resolution with `stty -a` and try resizing the window or reduce the font.


```
[pja@ SlurmCommander-dev]$ [DEBUG=1] [TERM=xterm-256color]./scom -h
Welcome to Slurm Commander!

Usage of ./scom:
  -d uint
        Jobs history fetch last N days (default 7)
  -t uint
        Job history fetch timeout, seconds (default 30)
  -v    Display version

```

To run in _debug_ mode, set `DEBUG` env. variable. You will see an extra debug message line in the user interface and scom will record a `scdebug.log` file with (_lots of_) internal log messages.

> Tested on: 
> * slurm 21.08.8
> * slurm 22.05.5

## Feedback

__Is most welcome. Of any kind.__

* bug reports
* broken UI elements
* code panic reports
* ideas
* wishes
* code contributions
* usage stories
* kudos
* coffee and/or beer
* ...

## Acknowledgments

Powered by amazing [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework/ecosystem. Kudos to glamurous Charm developers and community.

