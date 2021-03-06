
# 1. Try this application locally

Expose AWS credencials as environment variables:

```
$ export AWS_REGION=
$ export AWS_ACCESS_KEY_ID=
$ export AWS_SECRET_ACCESS_KEY=
```

Run with Docker Compose:

```
$ git clone https://github.com/pottava/dockerized-aws-x-ray.git
$ cd dockerized-aws-x-ray/sample
$ docker run --rm -v "$(pwd)":/go/src/github.com/pottava/dockerized-aws-x-ray/sample \
    -w /go/src/github.com/pottava/dockerized-aws-x-ray/sample/src \
    supinf/go-dep:0.5 ensure
$ docker-compose up
```

Open with your browser:

```
$ open http://localhost:9000
```


# 2. Build as a docker image

```
$ aws_account_id=$( aws sts get-caller-identity --query "Account" --output text )
$ repository=${aws_account_id}.dkr.ecr.${AWS_REGION}.amazonaws.com/xray-sample
$ docker build -t ${repository} .
```

Push it to ECR.

```
$ aws ecr create-repository --repository-name xray-sample
$ aws ecr get-login --no-include-email | sh
$ docker push ${repository}
```


# 3. Provision the sample stack

Choose your stack name & expose your keypair name:

```
$ STACK_NAME=
$ YOUR_KEYPAIR_NAME=
```

Create a stack:

```
$ aws cloudformation create-stack --stack-name ${STACK_NAME} \
  --template-body file://cfn/ecs.yaml \
  --parameters ParameterKey=InstanceType,ParameterValue=t2.small \
               ParameterKey=KeyName,ParameterValue=${YOUR_KEYPAIR_NAME} \
  --capabilities CAPABILITY_IAM
$ aws cloudformation wait stack-create-complete --stack-name ${STACK_NAME}
$ service_name=$( aws cloudformation describe-stacks --stack-name ${STACK_NAME} \
    --query 'Stacks[0].Outputs[?(OutputKey==`Service`)].OutputValue' \
    --output text )
$ aws ecs wait services-stable --cluster ${STACK_NAME} --services ${service_name}
```

Access to the endpoint:

```
$ open "http://$( aws cloudformation describe-stacks --stack-name ${STACK_NAME} \
    --query 'Stacks[0].Outputs[?(OutputKey==`LoadBalancer`)].OutputValue' \
    --output text )"
```


# 4. Update the task definition & service

Update the task definition:

```
$ task_name=$( aws cloudformation describe-stacks --stack-name ${STACK_NAME} \
    --query 'Stacks[0].Outputs[?(OutputKey==`Task`)].OutputValue' \
    --output text )
$ old=$( aws ecs describe-task-definition --task-definition ${task_name} \
    --query 'taskDefinition' )
$ cat << EOF > container-definitions.json
[
  {
    "name": "web",
    "image": "${repository}",
    "portMappings": [{"protocol": "tcp", "containerPort": 80, "hostPort": 0}],
    "environment": [
      {"name": "AWS_REGION", "value": "ap-northeast-1"},
      {"name": "MYSQL_USER", "value": "user"},
      {"name": "MYSQL_PASSWORD", "value": "pass"},
      {"name": "MYSQL_DATABASE", "value": "test"}
    ],
    "links": [
      "xray-daemon:xray",
      "gen-errors:err",
      "mysql:db"
    ],
    "logConfiguration": $(echo ${old} | jq '.containerDefinitions[0].logConfiguration'),
    "memoryReservation": 32,
    "memory": 100,
    "cpu": 10,
    "essential": true
  },
  {
    "name": "xray-daemon",
    "image": "pottava/xray:3.0",
    "logConfiguration": $(echo ${old} | jq '.containerDefinitions[0].logConfiguration'),
    "memoryReservation": 32,
    "memory": 100,
    "cpu": 10,
    "essential": false
  },
  {
    "name": "gen-errors",
    "image": "pottava/http-sw:1.0",
    "logConfiguration": $(echo ${old} | jq '.containerDefinitions[0].logConfiguration'),
    "memoryReservation": 32,
    "memory": 100,
    "cpu": 10,
    "essential": false
  },
  {
    "name": "mysql",
    "image": "mysql:5.7",
    "environment": [
      {"name": "MYSQL_ALLOW_EMPTY_PASSWORD", "value": "true"},
      {"name": "MYSQL_USER", "value": "user"},
      {"name": "MYSQL_PASSWORD", "value": "pass"},
      {"name": "MYSQL_DATABASE", "value": "test"}
    ],
    "logConfiguration": $(echo ${old} | jq '.containerDefinitions[0].logConfiguration'),
    "memoryReservation": 256,
    "memory": 768,
    "cpu": 100,
    "essential": false
  }
]
EOF
$ aws ecs register-task-definition --family ${STACK_NAME} \
    --task-role-arn $(echo ${old} | jq -r '.taskRoleArn') \
    --container-definitions file://container-definitions.json
```

Update ECS service:

```
$ new_task_arn=$( aws ecs list-task-definitions \
    | jq ".taskDefinitionArns | to_entries" \
    | jq "map(select(.value | index(\"${STACK_NAME}\")).value)" \
    | jq -r "sort | .[-1]" )
$ aws ecs update-service --cluster ${STACK_NAME} --service ${service_name} \
    --task-definition ${new_task_arn}
$ aws ecs wait services-stable --cluster ${STACK_NAME} --services ${service_name}
```

Access to the endpoint:

```
$ open "http://$( aws cloudformation describe-stacks --stack-name ${STACK_NAME} \
    --query 'Stacks[0].Outputs[?(OutputKey==`LoadBalancer`)].OutputValue' \
    --output text )"
```
