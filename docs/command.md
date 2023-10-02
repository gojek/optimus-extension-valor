# Command

## Description

This documentation explains more about the available commands in Valor.
Assuming Valor is already built (see [here](../README.md#how-to-run)), then when the user run the command:

```zsh
./out/valor
```

The following output will be shown in CLI:

```zsh
Usage:
  valor [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  execute     Execute pipeline based on the specified recipe
  help        Help about any command
  profile     Profile the recipe specified by path

Flags:
  -h, --help   help for valor

Use "valor [command] --help" for more information about a command.
```

## Profile

Profile is a command that will profile the specified Valor recipe. By default, the recipe being profiled is read from `valor.yaml` in the active directory. Like the following example:

```zsh
./out/valor profile
```

If the recipe from the provided example `valor.example.yaml` is copied into `valor.yaml`, then the output will be like the following:

```zsh
RESOURCE:
+--------------+--------+------+--------------------+-------------------------+
|     NAME     | FORMAT | TYPE |        PATH        |        FRAMEWORK        |
+--------------+--------+------+--------------------+-------------------------+
| user_account | json   | dir  | ./example/resource | user_account_evaluation |
+--------------+--------+------+--------------------+-------------------------+

FRAMEWORK:
+-------------------------+------------+---------------------+
|        FRAMEWORK        |    TYPE    |        NAME         |
+-------------------------+------------+---------------------+
| user_account_evaluation | definition | memberships         |
+                         +------------+---------------------+
|                         | schema     | user_account_rule   |
+                         +------------+---------------------+
|                         | procedure  | enrich_user_account |
+                         +------------+---------------------+
|                         | output     | std_output          |
+-------------------------+------------+---------------------+
```

The user can also specify other recipe by using flag `--recipe-path` like the following:

```zsh
./out/valor profile --recipe-path=valor.example.yaml
```

## Execute

Execute is a command that execute pipeline based on the provided recipe. By default, the recipe being executed is read from `valor.yaml` in the active directory. An example of running this command:

```zsh
./out/valor execute
```

Running the above command will execute all frameworks under `valor.yaml` recipe. This command also has several flags.

Flag | Description | Format
--- | --- | ---
--progress-type | specify the progress type during execution | currently available: `progressive` (default) and `iterative`
--recipe-path | customize the recipe that will be executed | it is optional. the value should be a valid recipe path

This command also has sub-command. The currently available sub-commands are explained below.

### Resource

Recipe is required in every execution. However, sometime the user just want a certain resource to executed. The user can define a new recipe to accomplish such requirement. Or, another alternative is by using resource sub-command. Resource sub-command is a command that allows the user to specify which resource to be executed. So, with he same recipe, the user can select the only resource they need to be executed. To use this sub-command, run it like the following:

```zsh
./out/valor execute resource
```

This sub-command inherit the flags from the execute command. In addition to that, it also has its own flags like the following.

Flag | Description | Format
--- | --- | ---
--name | the name of the resource to be executed | it is mandatory and should refer to the resource in the recipe
--format | the format of the input resource | if it's not specified, Valor will use the format in recipe. but if it is, then Valor will use it instead.
--path | the path of the input resource | if it's not specified, Valor will use the format in recipe. but if it is, then Valor will use it instead.
--type | the type of path for the specified resource | if it's not specified, Valor will use the format in recipe. but if it is, then Valor will use it instead.
