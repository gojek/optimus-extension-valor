# valor

## Description

Valor is a tool that can be used to validate and evaluate resource.

## Quick Start

### 1. Clone Repo

Clone this repository into the local by running the following command:

```zsh
git clone github.com/gojek/optimus-extension-valor
```

### 2. Create Recipe

Copy the recipe example `valor.example.yaml` to `valor.yaml`. The resource that will be validated and evaluated from the recipe is like the following:

```yaml
resources:
- name: user_account
  type: dir
  path: ./example/resource
  format: json
  framework_names:
  - user_account_evaluation
...
```

As stated in the recipe, the resource is available under directory `./example/resource`. One example of the resource is like the following:

```json
{
    "email": "valor@github.com",
    "membership_id": 1,
    "is_active": true
}
```

Note that the `membership_id` is in numeric.

### 3. Prepare Valor

To execute the pipeline, one can either:

* download the latest binary from [here](https://github.com/gojek/optimus-extension-valor/releases) and put it in this project directory, or
* try building it by following [this guide](#how-to-build)

### 4. Execute Pipeline

After the preparation is done, try executing Valor from the CLI like the following:

```zsh
./out/valor execute
```

And done. Based on the recipe, the output will be printed to the std out, or in other words in the CLI. The output should contain one of the modified input, like the following:

```zsh
...
---------------------------
example/resource/valor.json
---------------------------
email: valor@github.com
is_active: true
membership: premium
...
```

### 5. Explanation

What Valor does is actually stated in the recipe file `valor.yaml`. Behind the scene, the process is like the following:

1. read the first resource
2. execute framework pointed by field `framework_names`
3. in the framework, run validation based on `schemas`
4. if no error is found, load the required definition under `definitions`
5. if no error is found, execute procedures stated under `procedures`
6. if no error is found, write the result based on `output`

The explanation here is quite brief. For further explanation, try checkout the documentation [here](#documentation).

## Documentation

The complete documentation is available like the following:

* [recipe](./docs/recipe.md), which explains more about recipe and how to set it up for execution
* [command](./docs/command.md), which explains more about the available commands to be run

## Development

Community contribution is very welcome. Though, because of diverse background, one might not be able to easily understand the reasoning behind every decision or approach. So, the following is the general guideline to contribute:

1. create a dedicated branch
2. after all changes are done, create a PR
3. in that PR, please describes the issue or reasining or even approach for such changes
4. if no issue is found, merge will be done shortly

In order to follow the above steps, one might need to setup the local environment. The following is the general information to set it up.

### Dependency

#### GO Language

The GO programming language version `go1.17.1` need to be installed
in the system. Go to [this link](https://golang.org/doc/install) and follow
the instruction to install. To check GO version on the environment, run the following command:

```bash
go version
```

example output:

```bash
go version go1.17.1 darwin/amd64
```

#### How to Test

In this project root directory, run the following command to test:

```bash
make test
```

To show a more complete coverage and uncovered lines:

```bash
make coverage
```

You can check into `coverage.html` file in root project directory.
This command also will open interactive coverage tool in the browser.

#### How to Build

To build the binary executable, in this project root directory, run the following command:

```bash
make bin
```

There will be a new directory named `out` with an executable file `valor` as the result of the built project.

Alternatively, distribution could also be made by running the following command:

```bash
make dist
```

There will be a new directory named `dist` with few executable files.

#### How to Run

In order to run this project, after building the binary executable, run
the following command in this project root directory:

```bash
./out/valor
```
