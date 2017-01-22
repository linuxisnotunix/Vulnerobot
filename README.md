# Vulnerobot
Robot de collecte d'alertes de sécurité

## Build
```
make
```

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
   testing-develop#a219df3@2017-01-22-1734-UTC

COMMANDS:
     collect  Collect CVE from modules and add them to database
     list     List known CVE in database from a application list
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
