Supported tags and respective `Dockerfile` links:

・latest ([versions/1.0/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/1.0/Dockerfile))  
・1.0 ([versions/1.0/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/1.0/Dockerfile))  


## Usage

```
$ docker run --rm pottava/xray:1.0 --version
$ docker run --rm pottava/xray:1.0 --help
```

### Local

```
$ docker run -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY \
    --rm -p 2000:2000 pottava/xray:1.0 --local-mode --bind 0.0.0.0:2000
```

* with Docker-Compose:

```
version: "3"
services:

  app:
    image: <your-some-application>
    ports:
      - 80:80
    container_name: app

  xray:
    image: pottava/xray:1.0
    command: --local-mode --bind 0.0.0.0:2000
    environment:
      - AWS_REGION
      - AWS_ACCESS_KEY_ID
      - AWS_SECRET_ACCESS_KEY
    container_name: xray
```

### ECS

* with AWS CloudFormation:

```
TaskDef:
  Type: AWS::ECS::TaskDefinition
  Properties:
    ContainerDefinitions:
      - Name: app
        Image: <your-some-application>
        PortMappings: 
          - ContainerPort: 80
            HostPort: 0
        Environment:
          - Name: AWS_XRAY_DAEMON_ADDRESS
            Value: xray:2000
        Links:
          - xray-daemon:xray
        Cpu: 10
        Memory: 100
        MemoryReservation: 32
        Essential: true
      - Name: xray-daemon
        Image: pottava/xray:1.0
        Cpu: 10
        Memory: 100
        MemoryReservation: 32
    Family: xxxx
    TaskRoleArn: xxxx
```

* with AWS-CLI (JSON format for register-task-definition)

```
[
  {
    "name": "app",
    "image": "<your-some-application>",
    "portMappings": [
      {
        "protocol": "tcp",
        "containerPort": 80,
        "hostPort": 0
      }
    ],
    "environment": [
      {"name": "AWS_XRAY_DAEMON_ADDRESS", "value": "xray:2000"}
    ],
    "links": [
      "xray-daemon:xray"
    ],
    "cpu": 10,
    "memory": 100,
    "memoryReservation": 32,
    "essential": true
  },
  {
    "name": "xray-daemon",
    "image": "pottava/xray:1.0",
    "cpu": 10,
    "memory": 100,
    "memoryReservation": 32,
    "essential": false
  }
]
```
