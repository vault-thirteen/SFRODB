# SFRODB

## Simple File Read-Only DataBase

A very simplistic database engine which stores information in files and 
provides read-only access to data. 

The architecture of this engine implies simplicity in all possible aspects. 
Data is stored in files which are very easy to use and do not need complicated 
mechanisms or libraries. To reduce the frequency of disk operations and to 
increase the data retrieval speed the database engine uses an internal cache of 
popular items stored in random access memory (RAM). The network protocol is 
also very simple – the database uses its own very simple network protocol based 
on the TCP/IP. Only data retrieval operations are supported making this engine 
suitable for sharing of static content.  

An example of a client can be found in the `cmd\SFRODB\client` folder.

## Caching

Retrieved items automatically get into the cache to avoid future reads from a 
file storage. If for some reason a user wants to update the data file in the 
storage, the cached data must be removed from the cache, the API provides such 
functionality.  

## Dual Port Architecture

To provide additional protection, database uses separate ports for read 
operations and for non-read operations. By non-read operations we mean methods 
for removing a single item from cache and methods for cache cleaning, i.e. 
resetting the cache to an empty state.

## Pool of Clients

The library provides not only a single client for this database. A pool of 
clients is also available. The pool is able to fix broken connections 
automatically and has an adjustable size.

## Building

Use the `build_SFRODB.bat` script included with the source code.

## Installation

`go install github.com/vault-thirteen/SFRODB/cmd/SFRODB/client@latest`  
`go install github.com/vault-thirteen/SFRODB/cmd/SFRODB/server@latest`  

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

### Sample Client

`client.exe <server's host name> <server's main port> <server's aux port>`

Example:  
`client.exe localhost 12345 12346`

## Settings

Format of the settings' file for a server is quite simple. It uses line 
breaks as a separator between parameters. Inside a parameter, sub-parameters 
are separated with a single space symbol (" "). Described below are meanings 
of each line.

1. Hostname.
2. Main port.
3. Auxiliary port.
4. Data folder.
5. Parameters of the cache:
   * File extension for data items;
   * Maximum cache volume for data items (in bytes);
   * Maximum volume of a single data item (in bytes);
   * Item's TTL (in seconds).

**Notes**:
* File extension here is used as a normal extension together with a dot 
(period) prefix, because _Go_ language uses such format for file extensions. 
This is not good, but this is how _Golang_ works.
