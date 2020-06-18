# Advanced Topics in Distributed Systems

## Lab 4: Leader election in an anonymous network
### Members of the group
- Ana Mestre
- Cristian Ramirez

### Files:
- anon.go: Code for one concurrent machine (server & client), where the server can read from more than one client and the client can send messages to more than one server.
- anon.exe: executable file for anon.go
- Networks: a ppt file describing the networks for every test.
- Configuration folder: contains different configuration files:
    - configFile_linei.txt: where i = [1, 4], these are the configuration files for 4 different executions corresponding to the run_chat line.bat script.
    - configFile_ringi.txt: where i = [1, 4], these are the configuration files for 4 different executions corresponding to the run_chat ring.bat script.
    - configFile_i.txt: where i = [1, 5], these are the configuration files for 5 different executions corresponding to the run_chat2.bat script.
    - configFile_600i.txt: where i = [1, 5], these are the configuration files for 5 different executions corresponding to the run_chat.bat script.
- run_chat.bat: Script that executes anon.go with 5 different configuration files.
- run_chat2.bat: Script that executes anon.go with 5 different configuration files.
- run_chat line.bat: Script that executes anon.go with 4 different configuration files.
- run_chat ring.bat: Script that executes anon.go with 4 different configuration files.


### How to compile:
This action is not necessary since the executable is already provided.
- go build anon.go

### How to execute:
#### Individual execution:
  - Open a terminal (per server), execute the program: ./anon.exe <configuration file> <integer: network size>
#### Batch execution:
  - Launch (double click) run_chat.bat (or any other script) : this will open 5 terminals, each one launching an instance of anon.go with its corresponding configuration file (configFile_600X.txt).

### COMMENT: whatever configuration file that follows this format can be replaced:
List of IPS and its Port: \<IP address> : \<Port number>
#### For example:
127.0.0.1:6002 \
10.80.29.90:6001 \
127.0.0.1:6001
