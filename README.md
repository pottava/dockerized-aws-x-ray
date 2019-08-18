# Dockerized AWS X-Ray Daemon

[![pottava/xray](http://dockeri.co/image/pottava/xray)](https://hub.docker.com/r/pottava/xray/)

Supported tags and respective `Dockerfile` links:

・latest ([versions/3.1/amd/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/3.1/amd/Dockerfile))  
・3.1 ([versions/3.1/amd/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/3.1/amd/Dockerfile))  
・3.1-arm ([versions/3.1/arm/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/3.1/arm/Dockerfile))  
・3.0 ([versions/3.0/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/3.0/Dockerfile))  
・2.1 ([versions/2.1/Dockerfile](https://github.com/pottava/dockerized-aws-x-ray/blob/master/versions/2.1/Dockerfile))  

## Usage

```sh
$ docker run --rm pottava/xray:3.1 --version
AWS X-Ray daemon version: 3.1.0
$ docker run --rm pottava/xray:3.1 --help
Usage: X-Ray [options]
  -a --resource-arn Amazon Resource Name (ARN) of the AWS resource running the daemon.
  ..
```

### Local

```sh
$ docker run --name xray -d \
    -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY \
    -p 2000:2000/udp -p 2000:2000/tcp \
    pottava/xray:3.1 --region ${AWS_REGION} --local-mode
```

* with Docker-Compose:

```yaml
version: "2.1"
services:
  app:
    image: <your-some-application>
    ports:
      - 80:80
    environment:
      - AWS_XRAY_DAEMON_ADDRESS=xray:2000
    container_name: app
  xray:
    image: pottava/xray:3.1
    command: --region ${AWS_REGION} --local-mode
    environment:
      - AWS_ACCESS_KEY_ID
      - AWS_SECRET_ACCESS_KEY
    container_name: xray
```

### ECS

* with AWS CloudFormation:

```yaml
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
        Image: pottava/xray:3.1
        Cpu: 10
        Memory: 100
        MemoryReservation: 32
    Family: xxxx
    TaskRoleArn: xxxx
```

* with AWS-CLI (JSON format for register-task-definition)

```json
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
    "image": "pottava/xray:3.1",
    "cpu": 10,
    "memory": 100,
    "memoryReservation": 32,
    "essential": false
  }
]
```
