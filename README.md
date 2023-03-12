# go-tcp-tunnel

_A simple tool that initiates port forwarding based on a socat implementation written in Go_

## Prequiresite

Socat installed for machine.

## Usage

```sh
./go-tcp-tunnel [<json_dataset_file>]
```
Eg:

```sh
./go-tcp-tunnel example/example_dataset.json
```

_**Output**_

```sh
2023/03/12 12:01:25 Start creating forward from 10.51.78.127:4000 -> :4000 [netcat]
2023/03/12 12:01:25 -> Forward [netcat] done.
2023/03/12 12:01:25 Start creating forward from 10.51.78.127:3306 -> :4001 [mysql]
2023/03/12 12:01:25 -> Forward [mysql] done.
2023/03/12 12:01:25 Start creating forward from 10.51.78.127:8000 -> :3333 [service 1]
2023/03/12 12:01:25 -> Forward [service 1] done.
```

**Example dataset forwarding**

_example_dataset.json_

```json
[
    {
        "name": "netcat",
        "source_port": "4000",
        "dest": "10.51.78.127:4000"
    },
    {
        "name": "mysql",
        "source_port": "4001",
        "dest": "10.51.78.127:3306"
    },
    {
        "name": "service 1",
        "source_port": "3333",
        "dest": "10.51.78.127:8000"
    }
]
```