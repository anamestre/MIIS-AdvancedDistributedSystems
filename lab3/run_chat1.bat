@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,3) DO (
    SET "NUM=configFile_1%%x.txt
	start cmd.exe /k "go run machine.go !NUM!"
)
