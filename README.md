# exporter-idrac

- Using git
  - git branch feature/readme => create new branch
  - git checkout feature/readme => switch to feature/readme branch
  - git branch => check current working branch
  - git branch delete feature/readme => delete branch feature/readme(local PC), you should `git push origin --delete feature/readme` to delete branch on github.com
  - git fetch => to synchronize your local branch with branch on github.com
  - git checkout main && git merge feature/readme => first change current branch to main and then do merge feature/readme branch to main branch
  - git push --set-upstream origin feature/readme => create feature/readme on github.com
- How to Statically compile GOLANG programs - <https://oddcode.daveamit.com/2018/08/16/statically-compile-golang-binary/>
- Follow code with ultimate service which is taught by Jacob Walker
  - <https://github.com/ardanlabs/service/tree/32c75246b11b871ca7aaf07eebb2b1ccef6ee81c>
- <https://developer.aliyun.com/article/683180>
- <https://topic.alibabacloud.com/a/write-prometheus-exporter-using-golang_8_8_10262688.html>
- Folder structure

  ```bash
  ├── cmd
  │   └── idrac
  │       └── controller.go
  ├── config
  ├── migrate
  ├── go.mod
  ├── internal
  ├── main.go
  ├── Makefile
  └── README.md
  ```

  - The migrations directory will contain the SQL migration files for our database.
  - The internal directory will contain various ancillary packages used by our API. It will contain the code for interacting with our database, doing data validation, sending emails and so on. Basically, any code which isn't application-specific and can potentially be reused will live in here. Our Go code under cmd/idrac will import the packages in the internal directory (but never the other way around).
  - The Makefile will contain recipes for automating common administrative tasks — like auditing our Go code, building binaries, and executing database migrations.
  - The config directory will contain the configuration files and setup scripts for our production server

## How to Use

- Content of Makefile

  ```bash
  run:
    go run ./cmd/idrac


  mod:
  # real tab space or get error Makefile:2: *** missing separator.  Stop.
    go mod tidy # remove unused go packages
    go mod vendor # make local copy of third party packages

  ```

- run cli:
  - make run - run main.go in idrac folder
  - make mod
    - run go mod tidy cli
    - run go mod vendor cli
