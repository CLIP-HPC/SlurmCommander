
# SlurmCommander

> Development repository, highly volatile code. Unpublished.
>
> Work is in progress to finish planned features (first iteration) and stabilize the code (read: remove obvious/known bugs).
>
> Main branch is the safest bet (not guarantee) to compile and try out.

## Intro

SlurmCommander is a text-based user interface (TUI) to your slurm cluster.

It ties together multiple slurm commands to provide you with a simple and effective interaction point with your cluster.

You can view, search, analize and interact with:

* Job queue - view, search, inspect, ssh to nodes, issue commands to jobs etc.
* Job history
* Historical job details
* Edit and submit jobs from predefined templates (can be your own or group/site ones)
* Examine state of cluster nodes and partitions

![demo](./images/jobqueue.gif)

## Installation

SlurmCommander does not require any special privileges to be installed, see instructions below.

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

SlurmCommander is developed for 256 color terminals (black background) and requires at least 185x43 (rolumns x rows) to work.

* If you experience _funky_ colors on startup, try setting your `TERM` environment variable to something like `xterm-256color`.
* If you get a message like this:
`FATAL: Window too small to run without breaking view. Have 80x24. Need at least 185x43.`, check your terminal resolution with `stty -a` and try resizing the window or reducting the font.


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


