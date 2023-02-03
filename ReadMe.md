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

An example of a client can be found in the `example\client` folder.

## Building
Use the `build.bat` script included with the source code.

## Startup Parameters

### Server
`server.exe <path-to-configuration-file>`

Example:  
`server.exe "settings.dat"`


### Sample Client
`client.exe <server's host name> <server's port number>`

Example:  
`client.exe localhost 12345`

## Settings
Format of the settings' file can be learned by studying the source code.
