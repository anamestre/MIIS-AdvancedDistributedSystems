@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,5) DO (
    SET "NUM=configFile_600%%x.txt 5
	start cmd.exe /k "go run anon.go Configuration/!NUM!"
)
