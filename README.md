# Vulnerobot
Robot de collecte d'alertes de sécurité

## Get & Build
```
git clone https://github.com/linuxisnotunix/Vulnerobot.git && cd Vulnerobot && make
 - or -
go get -u -v github.com/linuxisnotunix/Vulnerobot
```

Binary for most platform can also be found [here](https://github.com/linuxisnotunix/Vulnerobot/releases).

## Start
```
./vulnerobot collect
./vulnerobot list
```

For more advance params : ```./vulnerobot help```
```
NAME:
   Vulnerobot - Index CVE related to a list of progs

USAGE:
   vulnerobot [global options] command [command options] [arguments...]

VERSION:
   testing

COMMANDS:
     collect, c  Collect CVE from modules and add them to database
     list, l     List known CVE in database from a application list
     info, i     Display global info like the of list plugins availables
     web, w      Start a web server to display result.
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d                   Turns on verbose logging [$DEBUG]
   --config value, -c value      Application list to monitor (default: "data/configuration")
   --database value, --db value  Application database to use (default: "data/sqlite.db")
   --help, -h                    show help
   --version, -v                 print the version
```

For more details please see the related [wiki](https://github.com/linuxisnotunix/Vulnerobot/wiki).
