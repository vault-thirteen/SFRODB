:: This script builds the SFHS server.
@ECHO OFF

SET build_dir=_build_\SFHS
SET exe_dir=cmd\SFHS
SET server_dir=server
SET settings_file=settings.txt
SET server_starter_script=start-server.bat
SET cert_script=create-certificates.bat

CD ..
MKDIR "%build_dir%"

:: Build the server.
CD "%exe_dir%\%server_dir%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%server_dir%.exe" ".\..\..\..\%build_dir%\"
CD ".\..\..\..\"

:: Copy some additional files for the server.
COPY "config\SFHS\%settings_file%" "%build_dir%\"
COPY "script\SFHS\%server_starter_script%" "%build_dir%\"
COPY "script\SFHS\%cert_script%" "%build_dir%\"
