# Azure Storage Queue SDK for GO sample 

Sample program for sending/receiving queue. 

# How to run the sample

## Prerequisite
Install Go lang (1.13+)

## Set Enviornment Variables

| Key | Description |
| ---- | ------ |
| ConnectionString | The Connection String for the Azure Storage Account |
| QueueName | The name of the queue |

## Run 

### receiver

Receive queue messages.

```bash
$ cd cmd/receive
$ go run receive.go
```

### sender

Send 100 messages to the Azure Storage Queue.

```bash
$ cd cmd/send
$ go run send.go 100
```

# How to debug

To see the behavior, you can debug it. This is the sample of the VSCode `.vscode/launch.json`. Then Start Debugging. 
It requires, VSCode [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go).

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "env": {
                "ConnectionString": "YOUR_STORAGE_ACCOUNT_CONNECTION_STRING_HERE",
                "QueueName":"hello"
            },
            "args": []
        }
    ]
}
```