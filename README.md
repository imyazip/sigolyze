# sigolyze
[![Go](https://github.com/imyazip/sigolyze/actions/workflows/tests.yml/badge.svg)](https://github.com/imyazip/sigolyze/actions/workflows/tests.yml)

Signature based anlysis library for golang

## Signature example
```json
{
    "name": "Example Signature",
    "patterns": [
      {
        "name": "Pattern 1",
        "value": "Value1",
        "is_regex": true
      },
      {
        "name": "Pattern 2",
        "value": "Value2",
        "is_regex": false
      }
    ],
    "tags": ["tag1", "tag2"],
    "meta": [
      {
        "name": "Meta1",
        "info": ["detail1", "detail2"]
      },
      {
        "name": "Meta2",
        "info": ["detail3"]
      }
    ]
  }

```