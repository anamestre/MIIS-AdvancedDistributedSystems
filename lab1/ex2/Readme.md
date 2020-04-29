# Advanced Topics in Distributed Systems

## Lab 1: Introduction to Distributed Systems Programming
### Files:
- machine.go: Code for one concurrent machine (server & client), where the server can read from more than one client and the client can write into more than one server.
- run_chat.bat: Script that executes machine.go with 5 different configuration files.
- configFile_600i.txt: where i = [1, 5], these are the configuration files for 5 different executions.

### How to compile:
This action is not necessary since the executables are already provided.
- go build machine.go

### How to execute:
#### Individual execution:
  - Open a terminal (per server), execute machine: ./machine <configuration file>
#### Batch execution:
  - Launch (double click) run_chat.bat : this will open 5 terminals, each one launching an instance of machine.go with its corresponding configuration file (configFile_600X.txt).

### COMMENT: whatever configuration file that follows this format can be replaced:
List of IPS and its Port: \<IP address> : \<Port number>
#### For example:
127.0.0.1:6002 \
10.80.29.90:6001 \
127.0.0.1:6001
