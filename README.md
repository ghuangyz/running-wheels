# running-wheels
Running-wheels is a task scheduler, which process a group of tasks represented as DAG.

### Build
This is how you build `running-wheels`
```
cd bin
./build.sh
```

### Run
This is how you run `running-wheels`
```
running-wheels run --filename=path/to/task/yaml/file
```

Full help
```
usage: running-wheels [<flags>] <command> [<args> ...]

A thread DAG task runner

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  run --filename=FILENAME [<flags>]
    Run a DAG task
```

### Task YAML file

Tasks are specified as YAML file, and a sample DAG specification will look like this:

```
tasks:
 - name: task1
   pipe: true
   commands:
    - cmd: echo
      args:
        - task1 running
    - cmd: tee
      args:
        - -i
        - task1.txt
 - name: task2
   pipe: true
   commands:
     - cmd: echo
       args:
         - task2 running
     - cmd: tee
       args:
         - -i
         - task2.txt
 - name: task3
   pipe: true
   commands:
     - cmd: cat
       args:
         - task1.txt
         - task2.txt
     - cmd: tee
       args:
         - -i
         - task3.txt
   depends:
     - task1
     - task2
 - name: task4
   pipe: true
   commands:
     - cmd: cat
       args:
         - task3.txt
   depends:
     - task3
 - name: task5
   pipe: true
   commands:
     - cmd: rm
       args:
         - task1.txt
         - task2.txt
   depends:
     - task3
 - name: task6
   pipe: true
   commands:
     - cmd: rm
       args:
         - task3.txt
   depends:
     - task4
```
