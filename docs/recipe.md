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
  type: file
  path: ./example/resource/user_account.json
  format: json
  framework_names:
  - user_account_validation
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
            <td><i>./example/resource/user_account.json</i></td>
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
            <td rowspan=2>framework_names</td>
            <td>indicates what frameworks to be executed against a resource. execution of one framework name to another is done <b>sequentially</b> and <b>independently</b>.</td>
            <td rowspan=2>each framework name should point to an existing framework</td>
            <td rowspan=2><i>user_account_validation</i></td>
        </tr>
        <tr>
            <td>sequentially means that every resource will be validated and/or evaluated against every framework (indicated by framework name) from top to buttom. independently means that the result of a framework execution against a resource is not carried to the next framework.</td>
        <tr>
    </tbody>
</table>
