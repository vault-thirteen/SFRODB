# SFRODB

This repository contains two products which are dependent of each other and may 
be considered to be two parts of a single product:

1. SFRODB
2. SFHS

The first one is the _Simple File Read-Only DataBase_.  
The second one is the _Static Files HTTP Server_.  
The second product is an _HTTP_ front-end part for the _SFRODB_ database.

Each of these two products has its own description file.
1. [SFRODB](doc/SFRODB/ReadMe.md)
2. [SFHS](doc/SFHS/ReadMe.md)

## Notes

Please, do note that earlier in the past time, _SFRODB_ and _SFHS_ were stored 
in two separate repositories. At present time, two these interdependent 
products are stored in a single repository. All files of the _SFHS_ repository 
were merged into the _SFRODB_ repository.

## Usage Example

1. First of all, you need to build both of the parts: the back-end part and the 
front-end part. Use the provided build scripts in the root folder of this 
repository: `build_SFRODB.bat` and `build_SFHS.bat`.


2. In the `_build_` folder, which was created, first start the back-end part 
(`_build_\SFRODB\server.exe`), then start the front-end part 
(`_build_\SFHS\server.exe`). If you have any other services already running on 
the default ports of the _SFRODB_ and _SFHS_ servers, an error will be shown.


3. Open your favourite web browser at the following address:
`http://localhost/sample`. A built-in sample _JSON_ file (`sample.json`) will be 
fetched by the browser and shown on the screen.
