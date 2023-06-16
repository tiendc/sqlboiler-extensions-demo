[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)

# SQLBoiler extensions demo

## Step-by-step

1. Add [sqlboiler-extensions](https://github.com/tiendc/sqlboiler-extensions) to your project
    ```shell
    # Method 1: use go get
    go get -u github.com/tiendc/sqlboiler-extensions@latest
    ```
   
    ```shell
    # Method 2: add the extensions to your project as a submodule
    # NOTE: if use this method, you need to use relative path in Makefile (--templates ./db/extensions/...)
    git submodule add --name "db/extensions"  https://github.com/tiendc/sqlboiler-extensions.git "db/extensions"
    git submodule update --init
    ```
   
2. Start the equivalent test env (there are tests for mysql, postgres, and cockroachdb)
    ```shell
    docker-compose -f docker-compose.<<db>>.yaml up -d
    ```
 
3. Generate DB models with SQLboiler (see [Makefile](./Makefile) for details)
    ```shell
    make prepare
    make gen-models-<<db>>
    ```

4. Run the test code (you can see the detailed usage in [main](./main))
    ```shell
    make run-test-<<db>>
    ```