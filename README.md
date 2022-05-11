# Recipe Count Test

1. Make config from example:
    
        make config

2. Copy your datafile to working directory and change `path` in config.


3. Run application from console:

        go run ./cmd/main.go

    or in Docker (do not forget to copy file into container):

         docker build -t recipe-count-test .
         docker run recipe-count-test 

Run tests:
   
      make test


[Task description](TODO.md)