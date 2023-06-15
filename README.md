[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)

# SQLBoiler extensions demo

## Step-by-step

1. Add the repo [sqlboiler-extensions](https://github.com/tiendc/sqlboiler-extensions) to your project as a submodule
    ```shell
    git submodule add --name "db/extensions"  https://github.com/tiendc/sqlboiler-extensions.git "db/extensions"
    git submodule update --init
    ```
   
2. Start the equivalent test env (there are tests for mysql, postgres, and cockroachdb)
    ```shell
    docker-compose -f docker-compose.<<db>>.yaml up -d
    ```
 
3. Generate DB models with SQLboiler
    ```shell
    make prepare-<<db>>
    make gen-models-<<db>>
    ```

4. Run the test code (you can see the detailed usage in [main](./main))
    ```shell
    make run-test-<<db>>
    ```