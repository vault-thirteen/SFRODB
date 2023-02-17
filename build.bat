:: This script builds the server and an example of a client.
@ECHO OFF

SET build_dir=_build_
SET exe_dir=cmd
SET client_dir=client
SET server_dir=server
SET settings_file=settings.dat
SET server_starter_script=start-server.bat
SET client_starter_script=start-client.bat

MKDIR "%build_dir%"

:: Build the server.
CD "%exe_dir%\%server_dir%"
go build
MOVE "%server_dir%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy some additional files for the server.
COPY "%exe_dir%\%server_dir%\%settings_file%" "%build_dir%\"
COPY "%exe_dir%\%server_dir%\%server_starter_script%" "%build_dir%\"

:: Build an example of a client.
CD "%exe_dir%\%client_dir%"
go build
MOVE "%client_dir%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy some additional files for the client.
COPY "%exe_dir%\%client_dir%\%client_starter_script%" "%build_dir%\"
