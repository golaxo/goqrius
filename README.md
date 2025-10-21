# GoQrius

> [!WARNING]
> GoQrius is under heavy develpment.

A powerful `filter` query parameter implementation based on [RFC 8040][rfc8040] for filtering and querying data with expressive syntax.

> [!NOTE]  
> This repository only contains the lexer-parser.
> 
> To check the current data-layer implementations check [Implementation](#-implementation).

## 🚀 Features

GoQrius provides a comprehensive set of logical operators for building complex filter expressions:

### Comparison Operators

- `eq`: to check for equality, e.g. `name eq 'John'`.
- `ne`: to check not equality, e.g. `name ne 'John'`.
- `gt`: to check greater than, e.g. `age gt 18`.
- `ge`: to check greater than or equal, e.g. `age ge 18`.
- `lt`: to check lower than, e.g. `age lt 21`.
- `le`: to check lower than or equal, e.g. `age le 21`.

### Logical Operators

- `not`: to negate the next check, e.g. `not name eq 'John'`.
- `and`: to AND concatenate conditions, e.g. `name eq 'John' and age gt 18`
- `or`: to OR concatenate conditions, e.g. `age le 18 or age ge 65`

## 📚 Examples

Imagine the following users

| Id | Name   | Surname | Age  |
|----|--------|---------|------|
| 1  | John   | Doe     | 20   |
| 2  | Jane   | Doe     | 10   |
| 3  | Alice  | Smith   | 66   |
| 4  | Bob    | Smith   | 30   |

Then it's expected to filter the following conditions:

| Filter                    | John(1) | Jane(2) | Alice(3) | Bob(4) |
|---------------------------|:-------:|:-------:|:--------:|:------:|
| `name eq 'John'`          |    ✅    |    ❌    |    ❌     |   ❌    |
| `name ne 'John'`          |    ❌    |    ✅    |    ✅     |   ✅    |
| `not name eq 'John'`      |    ❌    |    ✅    |    ✅     |   ✅    |
| `age gt 18 and age lt 65` |    ✅    |    ❌    |    ❌     |   ✅    |
| `age le 18 or age gt 65`  |    ❌    |    ✅    |    ✅     |   ❌    |

## 🔧 Implementation

GoQrius is designed to be easily integrated into your REST API endpoints.
The filter parameter parses the expression and converts it into executable conditions for your data layer.

The current data layers implementations for GoQrius are:

- [GormGoQrius](https://github.com/golaxo/gormgoqrius)

[rfc8040]: https://www.rfc-editor.org/rfc/rfc8040.html#section-4.8.4
