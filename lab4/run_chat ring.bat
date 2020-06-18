@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,4) DO (
    SET "NUM=configFile_ring%%x.txt 4
	start cmd.exe /k "go run anon.go Configuration/!NUM!"
)
