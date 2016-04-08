@ECHO OFF
setlocal enableDelayedExpansion
chdir C:\GoProj\src\github.com\rhino1998\cluster\test\0
start C:\GoProj\src\github.com\rhino1998\cluster\test\cluster.exe -port 3000
chdir C:\GoProj\src\github.com\rhino1998\cluster\test
for /L %%A in (1,1,50) do (
	mkdir %%A
	copy specs.json %%A\specs.json
        copy conf.json %%A\conf.json
	set /A port=%%A+3000
	echo port=!port!
	echo %%A 
	chdir C:\GoProj\src\github.com\rhino1998\cluster\test\%%A
	start C:\GoProj\src\github.com\rhino1998\cluster\test\cluster.exe -port !port!
	chdir C:\GoProj\src\github.com\rhino1998\cluster\test
)

