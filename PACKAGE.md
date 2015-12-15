- cmd
  - selfie
    + main.go
      -> reads config file and creates config.Config object
      -> creates selfie object and passes config.Config into it
      -> call selfie.Start()
- config
- data : can only access libs package inside selfie and other packages outside.
  + db.go
- etc : all configuration files should be here
- logme : uses private global variable which can be only access by public global Funcs
- lib : everything in lib has to be stateless, no global variables
        they must not depends on any other selfie packages.
  - crypto
  - utils
- scripts
- web
  - apis
    - auth
    - users
    - bundles
    - releases
  - middlewares : should be stateless, no global variable
  - security
    - jwt
    + security.go
      -> creates a global variable for each package
  + web.go
    -> setup the security package
    -> returns api's handlers

+ selfie.go
  -> setup logme : internally it creates a private global to its package
  -> setup db : internally it creates a global to its package, e.g. `data.DB`
  -> setup security : internally it creates a global variable for each package

  -> Start()
    -> instantiates web package and assign it to graceful package.
