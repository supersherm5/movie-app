# Microservices with Go
---
## Introduction
This repository contains the code for the `Microservices with Go` book. The book is designed to teach you how to use Go to build microservice-based systems. Writtien by Alexander Shuiskov. I did make modifications to the code to make it work with the latest version of Go and the libraries used in the book.
I also added a few more features to the code to make it more interesting and useful. The book is available on [Amazon](https://www.amazon.com/Microservices-Go-Alex-Shuiskov-ebook/dp/B07N6C4N3R).

## What is a microservice?
A microservice is a small, loosely coupled, distributed service. It is designed to be independently deployable and scalable. Microservices are typically built around business capabilities and are independently deployable. They are also designed to be independently scalable. This means that you can scale each microservice independently of the others. This is in contrast to a monolithic application, where the entire application is deployed and scaled as a single unit.

### What is Service Discovery?
Service discovery is the process of finding the location of a service. It is a critical part of building a microservice-based system. In a microservice-based system, services are constantly being deployed and scaled. This means that the location of a service is constantly changing. Service discovery is the process of finding the location of a service, so that you can communicate with it.

### What is a Service Registry?
A service registry is a database of services. It is a central location where services can register themselves, and where clients can look up the location of a service. A service registry is a critical part of building a microservice-based system. It is the central location where services can register themselves, and where clients can look up the location of a service.

#### Service Discovery Models
- Client-side discovery
- Server-side discovery

##### Client-side discovery
In client-side discovery, the client is responsible for finding the location of a service. The client is responsible for finding the location of a service, and then communicating with it. This is typically done using a service registry. The client queries the service registry to find the location of a service, and then communicates with it.

##### Server-side discovery
In server-side discovery, the server is responsible for finding the location of a service. The server is responsible for finding the location of a service, and then communicating with it. This is typically done using a load balancer. The load balancer is responsible for finding the location of a service, and then communicating with it.

### Popular Service Discovery Tools
- Consul
- Kubernetes

Consul is a popular service discovery tool. It is a distributed, highly available, and scalable service discovery tool. It is designed to be used in a microservice-based system. It is a central location where services can register themselves, and where clients can look up the location of a service. It was chosen for this repo.
Kubernetes is a popular container orchestration tool. It is designed to be used in a microservice-based system. It is a central location where services can register themselves, and where clients can look up the location of a service.

## Serialization
Serialization is the process of converting an object into a format that can be transmitted over a network. It is a critical part of building a microservice-based system. In a microservice-based system, services are constantly communicating with each other. This means that objects need to be serialized before they can be transmitted over the network.

### Popular Serialization Formats
- JSON
- Protocol Buffers
- Avro
- Thrift
- XML