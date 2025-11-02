# SFHS

## Static Files HTTP Server

An _HTTP_ server for serving static files. 

Static files are taken from the _SFRODB_ database, so this server may be 
considered as a front-end part for the _SFRODB_ database.

The server serves files of only one file type (extension). So, if you need to 
serve files having two different types (extensions), e.g. `json` and `txt`, you 
will need to start two separate servers, each for a separate file extension. 
Such an approach is used for increased security reasons which may be required
in some specific circumstances. 

The data folder (folder with all your files) can be a single folder when you 
use multiple instances of servers.

### HTTP Status Codes
200 – Successful data retrieval.  
400 – Client has requested wrong data (UID is bad or file does not exist).  
500 – Server error has occurred.  

## Architecture
_HTTP_ protocol is used for serving incoming requests.  

Static files are taken from the _SFRODB_ database.  
The server uses a pool of clients to connect to the _SFRODB_ database.

## Building
Use the `build_SFHS.bat` script included with the source code.

## Installation
`go install github.com/vault-thirteen/SFRODB/cmd/SFHS/server@latest`

## Startup Parameters

### Server
`server.exe <path-to-configuration-file>`  
`server.exe`  

Example:  
`server.exe "settings.txt"`  
`server.exe`  

**Notes**:  
If the path to a configuration file is omitted, the default one is used.  
Default name of the configuration file is `settings.txt`.

## Settings
Format of the settings' file for a server is quite simple. It uses line
breaks as a separator between parameters. Described below are meanings of each 
line.

1. Server's hostname.
2. Server's listen port.
3. Work mode: _HTTP_ or _HTTPS_.
4. Path to the certificate file for the _HTTPS_ work mode.
5. Path to the key file for the _HTTPS_ work mode.
6. Hostname of the _SFRODB_ database.
7. Main port of the _SFRODB_ database.
8. Auxiliary port of the _SFRODB_ database.
9. Size of the client pool for the _SFRODB_ database.
10. File extension of served files.
11. MIME type of served files.
12. TTL of served files, i.e. value of the `max-age` field of the 
`Cache-Control` _HTTP_ header.
13. Allowed origin for _HTTP_ CORS, i.e. value of the 
`Access-Control-Allow-Origin` _HTTP_ header.

**Notes**:
* File extension here may be set without a leading dot symbol. Dot symbol is
    appended to the start of the extension automatically.

## Performance

Performance test of the combination of _SFHS_ together with _SFRODB_ made 
in _Apache JMeter_ may be found in the `test` folder. Quite a decent hardware 
shows about 22 kRPS in _HTTPS_ mode and about 23 kRPS in _HTTP_ mode, while test 
file size was about 1kB.
