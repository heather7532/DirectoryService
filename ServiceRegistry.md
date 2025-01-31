# INSERT NAME HERE

# Introduction
The TBD is a simple framework that allows you to publish and discover services in a distributed environment. It is based on the concept of a service registry, which is a database of services that can be queried by clients to find the services they need. The Data Service Finder is designed to be lightweight and easy to use, with a minimal set of features that are sufficient for most use cases.
A service should be easy to deploy.  
+ It should be easy to find.
+ It should offer support for versioned and flexible APIs using multiple protocols.
+ It should be easy to scale.
+ It should be able to be monitized and secured

This project will dust off a few old ideas, modernize them, and put them together in a new way.  
Besides being simple, lightweight, and can be used by a single developer or a large enterprise, it will be open source and free to use.
This will also provide AI and ML services to help with the discovery and selection of services for Retrieval-Augmented Generation (RAG) and other AI/ML applications.

## Architecture

The TBD is composed of the following components:
### Service Registry
+ Registry database(s)
+ Service registration Tool and API
+ Service discovery API
+ RAG API to consume registry data that has been cleaned, normalized, and tokenized
+ Service monitoring and management API

### Service
+ API to authenticate and interface with the registry (service up/down, health check, etc.)
+ The Service Registry shall provide a standard API to start/stop the service. Implementation of the start/stop is up to the service provider.
+ The Service Registry shall provide a standard API to check the health of the service. Implementation of the health check is up to the service provider.


## SDK Support
The TBD will provide SDKs for the following languages:
+ Python
+ Go
+ Java 11+
+ Rust
+ Node.js
+ Others as needed

## Protocol Support
Intially, it will support HTTP/HTTPS. gRPC may be added later.

## Basic Flow
### Service Registration
There will be one or more service registries.  Public registries will require at least two instance (or more).  This will be backed by a PostgreSQL database.  Replication will be specific to the deployment for registries.  For instance, if deployed in AWS, it will use RDS Multi-AZ.  If deployed in Azure, it will use Azure SQL Database with Geo-Replication.  If deployed in Google Cloud, it will use Cloud SQL with High Availability.  
A service can be registered with any registry, and it shall be propagated to all other registries in the Registry Group. Service Registry(s) can be private or public.
Registry Services need not all be owned by the same entity, but they must cooperate with other members of the registry group.  The registry group will come from a bootstrap file that is loaded into the registry.  The registry group will be a list of registry URLs.
Registration of a service does not mean a service running/available, but a service must be registered to run.

#### Primary Registration Attributes
+ Service ID
+ Service Name
+ Service Description
+ Service Cost/Billing (TBD)
+ Service Owner Information (Contact/Support Link/Website)
+ Industry/Category (ISIC or NAICS)
+ Security Information (TBD)
+ Client Rating (5-star system)
#### Instance Registration Attributes (when service starts)
+ Service ID
+ Service Instance ID
+ Service Version (SemVer)
+ Service Host/Port/URL
+ Service API (OpenAPI 3.0 spec)
+ Lat/Lon (optional) - used for network optimization 

#### Registry of Registries
+ DNS information/URLs for the other registries in the group. This will be a bootstrap file that is loaded into the registry.

### Registry Viewer/Editor
There will be a website to view, register new services, update information, and delete services.  Users will also be able to leave reviews of services that will be used to rate the service.  The website will be a React app that will use the registry API to interact with the registry.


### Service Startup
When a service starts, it must update the registry.
The registry will start monitoring the service periodically checking health.  
If a service is unhealthy, it will be marked as down in the registry.
It is up to the service owner to manage the environment the service is executing in.  The registry will only monitor the service using the health check API.


### Client 
The client will query the registry for services based on the service name and version.  The client will receive a list of services that match the query.  The client can then choose a service to use.  
If there are no services available in that registry an option could be provided to search other registries in the registry group (this should only be needed if a registry is undergoing maintenance).  
The client can also specify a lat/lon to optimize the network path to the service.  
The client can also specify a version of the service to use.  If the client specifies the full semantic version (i.e. 1.2.3) the registry will only return services that match that version or newer.  If the client specifies a major version (i.e. 1) the registry will return services that match that major version or newer.  If the client specifies a major and minor version (i.e. 1.2) the registry will return any services that match that major and minor version or newer.

Once the client has chosen a service, it will query for service instances.  The client will receive a list of service instances that match the query.  The client can then choose a service instance to use. If the service instance has latitute and longitude, the list of service instances will be sorted by distance from the client determined by a Great Circle Distance calculation.

If the client uses the Client SDK, the SDK will handle the service and instance selection.  It will also cache the service instance list and periodically refresh the list.  It can also be used to handle retries if the service fails to respond in a timely manner.  It can automatically repeat the call to the next service in the list.
The client will have the ability to give feedback on the service resulting in a star rating.

#### Versioning
The client can specify a version of the service to use.  If the client specifies the full semantic version (i.e. 1.2.3) the registry will only return services that match that version or newer.  If the client specifies a major version (i.e. 1) the registry will return services that match that major version or newer.  If the client specifies a major and minor version (i.e. 1.2) the registry will return services that match that major and minor version or newer.
This will support backward compatibility and allow the client to support multiple versions of a service.  This makes it easier to upgrade services without breaking clients.

### Service Owner Responsibilites
The service owner is responsible for the following:
+ Registering the service with the registry
+ Invoking the registry API to update the service status when it starts and stops
+ Providing a health check API that the registry can use to determine if the service is healthy
+ Providing a support contact for the service
+ Providing a billing mechanism for the service (TBD)
+ Providing a security mechanism for the service (TBD)
+ Providing a service API that conforms to the OpenAPI 3.0 spec
+ Providing a version number for the service
+ Providing a lat/lon for the service (optional)

## Enterprise Level Features
The Data Service Finder will provide the following enterprise level features:
+ Local instance management/monitoring tool(s)
    + Auto start/stop of services
    + Health check monitoring
    + Restart services that are unhealthy
    + Metrics SDK add on
+ Replication of the registry between different cloud providers 
+ Support for monitizing services and billing
+ Support for user managed security and encryption
