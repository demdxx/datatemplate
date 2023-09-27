# DataTemplate Module

[![Build Status](https://github.com/demdxx/datatemplate/workflows/Tests/badge.svg)](https://github.com/demdxx/datatemplate/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/datatemplate)](https://goreportcard.com/report/github.com/demdxx/datatemplate)
[![GoDoc](https://godoc.org/github.com/demdxx/datatemplate?status.svg)](https://godoc.org/github.com/demdxx/datatemplate)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/datatemplate/badge.svg)](https://coveralls.io/github/demdxx/datatemplate)

The Go Templating Module is a versatile and powerful tool for dynamically generating text and data by processing templates with dynamic content. This module allows you to seamlessly substitute placeholders within templates with real data, making it an ideal choice for creating dynamic reports, generating emails, rendering dynamic content, and even generating dynamic configuration files in JSON or YAML formats for your Go applications.

## Features

- **String and Map Templates**: The module supports both simple string templates and more complex map templates. String templates are useful for straightforward text replacement, while map templates enable you to structure data hierarchically.

- **Conditional Statements**: Incorporate conditional logic into your templates using `$if` and `$else` clauses. Conditionally include or exclude content based on values in the data context.

- **Iteration**: Use `$iterate` clauses to iterate over lists or arrays of data, generating multiple instances of output based on the data provided in the context.

- **Expression Evaluation**: Evaluate expressions enclosed in double curly braces, such as `{{person[0].age > 18 ? 'adult' : 'teenager'}}`, to dynamically compute values during template substitution.

- **Custom Logic**: Implement custom logic within your templates using expressions like `{{s= index + 1}}`, enabling advanced data processing during template rendering.

Custom expression [syntax](https://expr.medv.io/docs/Language-Definition) is supported through the use of the [github.com/antonmedv/expr](https://github.com/antonmedv/expr) library.

## Usage Examples

Explore the provided usage examples and test cases to see how to effectively utilize the Go Templating Module for your specific needs, including the generation of dynamic configuration files in JSON or YAML formats.

## Dynamic Configuration Files

One powerful use case of this module is the generation of dynamic configuration files. You can create templates for JSON or YAML configuration files with placeholders for values that may change depending on your application's environment or user-defined settings. By processing these templates with the Go Templating Module and providing the appropriate data context, you can generate configuration files tailored to specific scenarios, making your application highly adaptable and flexible.

### Example: Generating Dynamic YAML Configuration

Consider the following YAML template:

```yaml
# Template
database:
  $iterate: parseURI(databases)
  host: {{item.host}}
  port: {{item.port}}
  username: {{item.username}}
  password: {{item.passord}}
```

Input example

```json
{
  "databases": [
    "mongodb://user@password:localhost:27017",
    "mysql://user@password:localhost:3306"
  ]
}
```

When processed with the Go Templating Module, the output configuration file will be dynamically generated as follows:

```yaml
# Output
database:
  - host: localhost
    port: 27017
    username: user
    password: password
  - host: localhost
    port: 3306
    username: user
    password: password
```

In this example, the module iterates through the databases array, parses the URI strings, and populates the YAML template with the extracted values, resulting in a customized YAML configuration file that reflects the provided data context.

This illustrates how the Go Templating Module can streamline the process of generating dynamic configuration files, making it a valuable tool for configuring your applications dynamically based on various input data.

> This addition provides a clear and detailed example of how the module can be used to generate dynamic YAML configuration files from templates and input data, showcasing its practical utility.

### Example: in GO

```go
tpl, err := NewTemplateFor(map[string]any{
  "database": map[string]any {
    "$iterate": "parseURI(databases)",
    "host": "{{item.host}}",
    "port": "{{item.port}}",
    "username": "{{item.username}}",
    "password": "{{item.password}}",
  },
})
if err != nil {
    panic(err)
}

result, err := tpl.Process(context.Background(), map[string]any{
  "databases": []string{
    "mongodb://user@password:localhost:27017",
    "mysql://user@password:localhost:3306",
  },
})
if err != nil {
    panic(err)
}

fmt.Println(result) // Output: map[database:[map[host:localhost port:27017 username:user password:password] map[host:localhost port:3306 username:user password:password]]]
```

## Contributing

We welcome contributions from the community to enhance and expand the capabilities of this module. If you have ideas for improvements or encounter issues, please feel free to contribute by opening a pull request or submitting an issue.

## Dependencies

- [github.com/antonmedv/expr](http://github.com/antonmedv/expr)

## License

This Go Templating Module is open-source software released under the [Apache 2.0 License](LICENSE).
