# Terminology

## Domains
Domains have a name, resources, zero or more child domains, and 1 or more applications.

### Organization
A top level domain. A business, society, association, etc. given a specific name and probably associated with a registered domain name.  Organizations can contain zero or more child domains. Authentication could be implemented at the organization level and enable single sign on for all the contained domains.

### Micro Service Domain
Child domains of the organization and is a bounded context within the organization. It models data to accomplish it's specific work.  It has its own domain language, cloud resources, roles and permissions, data storage, and it can have one or more applications.  It is a micro service.

[Bounded Context](https://martinfowler.com/bliki/BoundedContext.html)

### Application
Executable code handling domain requirements

## Resource Group Descriptors 

### Global
Resources available to the entire organization. Across all domains, envs and instances.
Examples would be, container repository, table storage with deployment meta data

### Environment
Dedicated environments of work.
(ci, dev, qa, demo, canary, prod)

### Common
Resources available to all instances in the env

### Manual
Resources that have to be created by hand

### Instance
Independant resources




Environment
Instance
    stbl

Types of Names: 
    domainName - global
    domainName - env - manual
    domainName - env - common
    domainName - env - instance
