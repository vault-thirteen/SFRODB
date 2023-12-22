:: This script builds the SFRODB server and an example of a SFRODB client.
@ECHO OFF

SET build_dir=_build_\SFRODB
SET exe_dir=cmd\SFRODB
SET client_dir=client
SET server_dir=server
SET settings_file=settings.txt
SET sample_data_file=sample.json
SET client_starter_script=start-client.bat
SET data_dir=data

MKDIR "%build_dir%"

:: Build the server.
CD "%exe_dir%\%server_dir%"
go build
MOVE "%server_dir%.exe" ".\..\..\..\%build_dir%\"
CD ".\..\..\..\"

:: Copy some additional files for the server.
COPY "%exe_dir%\%server_dir%\%settings_file%" "%build_dir%\"
MKDIR "%build_dir%\%data_dir%"
COPY "%exe_dir%\%server_dir%\%sample_data_file%" "%build_dir%\%data_dir%\"

:: Build an example of a client.
CD "%exe_dir%\%client_dir%"
go build
MOVE "%client_dir%.exe" ".\..\..\..\%build_dir%\"
CD ".\..\..\..\"

:: Copy some additional files for the client.
COPY "%exe_dir%\%client_dir%\%client_starter_script%" "%build_dir%\"
