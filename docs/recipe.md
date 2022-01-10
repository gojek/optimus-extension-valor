# Recipe

Recipe describes how Valor should execute a pipeline. A pipeline
is one execution flow of Valor. Each recipe consists of
one or more **[resources](#resource)** and also
one or more of its framework **[frameworks](#framework)**.

## Resource

Resource is something to be either validated, evaluated, or both.
Each recipe should contains one or more resources. An example:

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

Each resource is defined by a structure with the following fields:

<table>
    <thead>
        <tr>
            <th>Field</th>
            <th>Description</th>
            <th>Format</th>
            <th>Example</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>name</td>
            <td>a unique name of the resource</td>
            <td>it is suggested to be descriptive and needs to follow regex <i>[a-z_]+</i></td>
            <td><i>user_account</i></td>
        </tr>
        <tr>
            <td>path</td>
            <td>the path where the resource should be read from</td>
            <td>it has to follow the <b>type</b> format. for example, if the <b>type</b> is <i>dir</i> (to indicate local directory), then the format should follow how a directory path looks like</td>
            <td><i>./example/resource</i></td>
        </tr>
        <tr>
            <td rowspan=2>type</td>
            <td rowspan=2>describes the type of path in order to get the resource</td>
            <td>currently available: <i>file</i> and <i>dir</i></td>
            <td rowspan=2><i>dir</i></td>
        </tr>
        <tr>
            <td>
                <ul>
                    <li>file</li>
                    describe that the path is of type <i>file</i>. if the <b>path</b> value is actually a directory but the <b>type</b> is set to be a <i>file</i>, then only the first file in that directory will be read.
                    <li>dir</li>
                    describe that the path is of type <i>dir</i>. if the <b>path</b> value is not a directory, then error might be returned.
                </ul>
            </td>
        </tr>
        <tr>
            <td>format</td>
            <td>indicates what format a resource was stored</td>
            <td>currently available: <i>json</i> and <i>yaml</i></td>
            <td><i>json</i></td>
        </tr>
        <tr>
            <td>framework_names</td>
            <td>indicates what frameworks to be executed against a resource. execution of one framework name to another is done <b>sequentially</b> and <b>independently</b>.</td>
            <td>each framework name should point to an existing framework</td>
            <td><i>user_account_validation</i></td>
        </tr>
    </tbody>
</table>

_Note that every field mentioned above is mandatory unless stated otherwise._

## Framework

Framework describes how to validate and/or evaluate a resource and how to return the result. One framework can be used by multiple resources. An example of framework:

```yaml
...
frameworks:
- name: user_account_evaluation
  schemas:
  - name: user_account_rule
    type: file
    format: json
    path: ./example/schema/user_account_rule.json
  definitions:
  - name: memberships
    format: json
    type: dir
    path: ./example/definition
    function:
      type: file
      path: ./example/procedure/construct_membership_dictionary.jsonnet
  procedures:
  - name: enrich_user_account
    type: file
    format: jsonnet
    path: ./example/procedure/enrich_user_account.jsonnet
    output_is_error: false
  output_targets:
  - name: std_output
    type: std
    format: json
```

The following is the general constructs for a framework:

Field | Required | Description | Format | Output
--- | --- | --- | --- | ---
name | true | defines the name of a particular framework. | it is suggested to be descriptive and needs to follow regex _`[a-z_]+`_ | -
[schemas](#schema) | false | defines how to validate a resource. | it is an array of `schema` that will be executed _sequentially_ and _independently_.| for each schema, the output of validation is either a success or an error message.
[definitions](#definition) | false | definitions are data input that might be required by **procedure**. **definitions** helps evaluation to be more efficient when external data is referenced multiple times. | it is an array of `definition` that defines how a definition should be prepared. | for each definition, the output is expected to be an array of JSON object.
[procedures](#procedure) | false | defines how to evaluate a resource. | it is an array of `procedure` that will be executed sequentially with the ability to pass on information from one procedure to the next. | vary, dependig on how the procedure is constructed.
output_targets | false | defines how an output of validation and/or evaluation should be written out | it is an array of `output_target` that will be executed sequentially and independently. | vary, dependig on the validation and/or evaluation output

### Schema

Schema is mainly used for validation. A schema composes of one or more rules on how a data should look like. Currently, schema only follows the specification by [JSON schema](https://json-schema.org/specification.html). The following is an example of basic construct of a schema in a framework:

```yaml
...
name: user_account_rule
type: file
format: json
path: ./example/schema/user_account_rule.json
...
```

Field | Description | Format
--- | --- | ---
name | the name of schema | it has to be unique within a framework only and should follow _`[a-z_]+`_
type | the type of data to be read from the path specified by **path** | currently available is `file` only
format | the format being used to decode the data | currently available is `json` only, pointing that it's a JSON schema
path | the path where the schema rule to be read from | the valid format based on the **type**

_Note that every field mentioned above is mandatory unless stated otherwise._

And the following is an example of JSON schema, pointed by **path**:

```json
{
    "title": "user_account",
    "description": "Schema to validate user_account.",
    "type": "object",
    "properties": {
        "email": {
            "type": "string"
        },
        "membership_id": {
            "type": "integer"
        },
        "is_active": {
            "type": "boolean"
        }
    },
    "required": [
        "email",
        "membership_id"
    ],
    "additionalProperties": false
}
```

The above example is a validation rule for data `user_account`, where for every value in its `email` field should be a `string`, its `membership` field should an `integer`, and its `is_active` field should be a boolean. If any of the actual resource (or one could say, record) does not comply, then error will be triggered.

## Definition

Definition is external data that could be used by **procedures**. Definition is usually utilized when one or more procedures want to load one or more externals data once and use it multiple times efficiently. Definition is like a static reference data. An example of definition construct in a framework:

```yaml
...
name: memberships
format: json
type: dir
path: ./example/definition
function:
    type: file
    path: ./example/procedure/construct_membership_dictionary.jsonnet
...
```

Field | Description | Format | Output
--- | --- | --- | ---
name | the name of definition | it has to be unique within a framework only and should follow _`[a-z_]+`_ | -
type | the type of data to be read from the path specified by **path** | currently available is `file` and `dir` | -
format | the format being used to decode the data | currently available is `json` and `yaml` | -
path | the path where to read the actual data from | the valid format based on the **type** | -
function | an optional instruction to build a definition, where the instruction follows the [Jsonnet](https://jsonnet.org/) format | - | dictionary where the key is the **name** and the value is up to the actual function defined under **function.path**
function.type | defines the type of path specified by **function.path** | it should be valid for the given **function.path** | -
function.path | defines the path where to read the actual function | should be valid according to the **function.type** | -

_Note that every field mentioned above is mandatory unless stated otherwise._

As mentionend, every definition **function** should follow [Jsonnet](https://jsonnet.org/) format. Apart from that, there are a few additional rules involved when defining a definition:

* the final definition output depends on the **function**:
  * if **function** is not set, then the actual output will be a dictionary where the key is the definition name and the value is an array
  * if **function** is set, then the actual output will be a dictionary where the key is the definition name and the value is up to the actual function to define
* every definition function should define a special [Jsonnet](https://jsonnet.org/) function with the following requirement:
  * it has to be named `construct`
  * it accepts one parameter
  * it outputs one value
* data being passed as the parameter of [Jsonnet](https://jsonnet.org/) function is the raw data, which is an array of definition object
* definition object is the actual data that is stored in the preferred place, such as a file
* if the special function requires custom functions, then they should be initialized above the special function

The following is an example of actual definition function:

```jsonnet
local construct (definitions) = {
    [std.toString(d.id)]: d
    for d in definitions
};
```

As shown above, there's only one function named `construct`. This is a special function that will be called by Valor, much like a "main" function. If needed, then the user can define some custom functions, like:

```jsonnet
local custom_function() {
    // do something
};

local construct (definitions) = {
    custom_function(),

    [std.toString(d.id)]: d
    for d in definitions
};
```

The output when the definition function is not set:

```json
{
    "memberships": [
        {
            "id": 1,
            "name": "premium",
            "description": "Membership which involves payment"
        }
    ]
}
```

When the definition function is defined where it outputs an object, then:

```json
{
    "memberships": {
        "1": {
            "id": 1,
            "name": "premium",
            "description": "Membership which involves payment"
        }
    }
}
```

## Procedure

Procedure is, like the name, one or more instruction to process data. Think of it like a the GO function. Procedure uses [Jsonnet](https://jsonnet.org/) format. Even though it's similar with [definition](#definition) function in term of the format being used, it's acually different. If a [definition](#definition) function's purpose is to accept all external definition data and produces new data, then procedures's purpose is to accept every possible data and may or may not proceds new data. It might be a bit abstract, but let's take a look at its basic construct:

```yaml
...
name: enrich_user_account
type: file
format: jsonnet
path: ./example/procedure/enrich_user_account.jsonnet
output_is_error: false
...
```

Field | Description | Format
--- | --- | ---
name | the name of a procedure | it has to be unique within a framework only and should follow _`[a-z_]+`_
type | the type of data to be read from the path specified by **path** | currently available is `file` only | -
format | the format being used to decode the data | currently available is `jsonnet` only | -
path | the path where to read the actual data from | the valid format based on the **type** | -
output_is_error | indicates whether output from a procedure is considered as error or not | it is optional, with the default value is `false`, where output is not considered as error

_Note that every field mentioned above is mandatory unless stated otherwise._

As mentioned, procedure follows the [Jsonnet](https://jsonnet.org/) format. Though, there are some rules for it to be executed properly by Valor:

* each procedure should have special [Jsonnet](https://jsonnet.org/) function named `evaluate`, which:
  * accepts `resource`, `definition`, and `previous` parameter sequentially, and
  * may or may not return data, depending on how the function is defined
* `resource` parameter in `evaluate` function refers to one resource data defined under **resources**, which is in a JSON format
* `definition` parameter in `evaluate` function refers to the whole definition defined under **definitions**, which is in the form of dictionary (JSON object) where the key is the definition name
* `previous` parameter in `evaluate` function refers to the output of the previous procedure execution within a framework and it will be a null value if the current procedure is the first to be executed
* any additional function which might be required by the special function should be initialized beforehand

The following is an example of a procedure:

```jsonnet
local evaluate(resource, definition, previous) =
    local membership_dict = definition['memberships'];
    local membership_id = std.toString(resource.membership_id);
    local current_membership = membership_dict[membership_id];
    {
        email: resource.email,
        membership: current_membership.name,
        is_active: resource.is_active
    };
```

On the procedure above, there's only one function, which is `evalute`. Behind the scene, Valor will call this function. In the line,

```jsonnet
...
local membership_dict = definition['memberships'];
...
```

this function wants to extract a definition named `memberships`, and use it as a reference to process its business flow. This function may or may not use the provided parameters, and may or may not return any output. It is entirely up to the user. In the above example, this function outputs an object. And since `output_is_error` is set to be `false`, then this value will be sent to the next pipeline, which can be:

* a new procedure, where this output will be sent as parameter under `previous`, or
* an output target, where this output will be written out to output stream, or
* nothing, where the output will not be used.

## Output Target

Output target defines how an output would be written. One framework can zero or more output targets. If output target is not defined, then even if procedure returns output, it won't be written out. If output target is defined, but procedure doesn't return output, then no output will be written out. The basic construct of output target is like the following:

```yaml
...
name: std_output
type: std
format: json
path: ./output
...
```

Field | Description | Format
--- | --- | ---
name | the name of output target | it has to be unique within a framework only and should follow _`[a-z_]+`_
type | the type of output target | currently available: `std` (to write to standard output) and `dir` (to write to a directory)
format | the format for output | currently available: `json` and `yaml`
path | the path where to write the output | only being used when the **type** is set `dir`