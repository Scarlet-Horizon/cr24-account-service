# Description

<div align="center">
    <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
    <img alt="Docker" src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white">
    <img alt="Shell Script" src="https://img.shields.io/badge/Shell_Script-121011?style=for-the-badge&logo=gnu-bash&logoColor=white"/>
    <img alt="AWS" src="https://img.shields.io/badge/Amazon_AWS-232F3E?style=for-the-badge&logo=amazonaws&logoColor=white"/>
    <img alt="Amazon EC2" src="https://img.shields.io/badge/Amazon_EC2-FF9900?style=for-the-badge&logo=amazonec2&logoColor=white"/>
    <img alt="RabbitMQ" src="https://img.shields.io/badge/RabbitMQ-FF6600?style=for-the-badge&logo=RabbitMQ&logoColor=white"/>
    <img alt="Amazon DynamoDB" src="https://img.shields.io/badge/Amazon_DynamoDB-4053D6?style=for-the-badge&logo=Amazon-DynamoDB&logoColor=white"/>
    <img alt="Swagger" src="https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=Swagger&logoColor=white"/>
</div>

API for managing bank accounts. The API allows you to create, read, update and delete bank accounts.

The API was made with **Go** and was hosted on **AWS EC2** as a container with **Docker**. For the database, we chose
**DynamoDB**, which was hosted on **AWS DynamoDB**.

We used **RabbitMQ** for log management and **Swagger** for documentation.

## Table of Contents

- [Description](#description)
    - [Table of Contents](#table-of-contents)
    - [Prerequisites](#prerequisites)
    - [How to run](#how-to-run)
    - [Usage](#usage)
    - [Contributor](#contributor)

## Prerequisites

Go to the `db/` directory and create a `.env` file. Add the following environment variables:

```env
AWS_ACCESS_KEY_ID = <your_aws_access_key_id>
AWS_SECRET_ACCESS_KEY = <your_aws_secret_access_key>
REGION = <your_aws_region>
```

All these variables can be set to any value. They are only used to create the database locally.

Create a `.env` file in the `env/` directory and add the following environment variables:

```env
AWS_ACCESS_KEY_ID = <your_aws_access_key_id>
AWS_SECRET_ACCESS_KEY = <your_aws_secret_access_key>
AWS_REGION = <your_aws_region>
GIN_MODE = release
JWT_SECRET = <your_jwt_secret>
AMQP_URL = <your_amqp_url>
EXCHANGE_QUEUE_NAME = <your_exchange_queue_name>
```

`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_REGION` must the same as in the `db/.env` file. `AMQP_URL`
and `EXCHANGE_QUEUE_NAME` are optional. If you do not specify them, the logs will not be sent to the queue.

## How to run

Firstly, you need to install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).

Then run:

```shell
git clone
```

or download the repository.

Then run

```shell
docker-compose up
```

in the root directory of the project.

## Usage

The API is documented with Swagger. You can access the documentation at `http://localhost:8080/swagger/index.html`.

The API is protected with JWT. To get the token, you need to generate it. You can do this online
at [jwt.io](https://jwt.io/).

The token must contain the following payload:

```json
{
  "sub": "<your_id>",
  "iat": "<current_timestamp>",
  "exp": "<current_timestamp + desired_expiration_time>"
}
```

The JWT token must be sent in the `Authorization` header in the following format:

```text
Authorization <your_jwt_token>
```

## Contributor

<table>
    <tbody>
        <tr>
            <td align="center">
                <a href="https://github.com/david-slatinek">
                    <img src="https://avatars.githubusercontent.com/u/79467409?v=4" width="100px;" alt="David Slatinek Github avatar"/>
                    <br/>
                    <sub><b>David Slatinek</b></sub>
                </a>
            </td>
        </tr>
    </tbody>
</table>
