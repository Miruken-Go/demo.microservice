## Global Resources Deployment Flow

Global resources are shared across all environments.   There will not be many of these types of resources because in general you want strict environment seperation.   The container repository is an example of one such resource.  It is shared across all the environments because we want the exact same bits that were compiled and tested to be deployed to production.  This minimizes the risk of hidden and unintended changes from creaping in after testing.  An example of hidden and unintended changes would be new versions of code dependencies being pulled in during a package restore.

```mermaid

flowchart LR
    A(Deploy Global Resources)
    
    B(Build Default Container Image)

    A --> B
```

## Environment Deployment Flow

Common environment resources are resources that are shared by all the invironment instances.  For example you may decide to share a database to reduce maintainence and mock data creation.  You may also want to share a single key store per environment.  These are long lasting resources that will be recreated less frequently.

Manual resources are created by hand.  In an ideal world there would be no resources created by hand, but some resources do not have the proper support from ARM/Bicep, AZ, or the rest api do be deployed successfully and consistently from a pipeline.  Azure B2C is an example of a resource that needs to be created by hand.  It cannot as of November 2023 be configured with AZ.   The ARM templates have serious limitations.  They can only be deployed once and cannot be modified.  Finally the rest API can only request an autherization token with the "implicit" flow.  It will need to suppor the "client credentials" flow before it is an option.  

```mermaid

flowchart LR
    B(Deploy Common Environment Resources)

    C[[Create Manual Environment Resources </br></br> B2C </br> Env Configuration]]

    D(Configure Manual Resources)

    E(Deploy Environment Instance Resources </br></br> Configure Manual Resources)

    B --> C --> D --> E
```


# Resource Groups
```mermaid

flowchart

    G(
        teamsrv-global

        Container Repository
    ) 
    G  --- DS(
        teamsrv-dev-common

        Key Vault
        Cosmos DB
    )
    DS  --- DM(
        teamsrv-dev-manual

        B2C Tenant
    )
    DM --- D1(
        teamsrv-dev

        Container Apps Environment
        Container App
    )
    DM --- D2(
        teamsrv-dev-ci

        Container Apps Environment
        Container App
    )
    DM --- D3(
        teamsrv-dev-developerA
    
        Container Apps Environment
        Container App
    )

    G  --- QS(
        teamsrv-qa-common

        Key Vault
        Cosmos DB
    )
    QS  --- QM(
        teamsrv-qa-manual

        B2C Tenant
    )
    QM --- Q1(
        teamsrv-qa-1

        Container Apps Environment
        Container App
    )
    QM --- Q2(
        teamsrv-qa-2

        Container Apps Environment
        Container App
    )

    G  --- PS(
        teamsrv-prod-common

        Key Vault
        Cosmos DB
    )

    PS  --- PM(
        teamsrv-prod-manual

        B2C Tenant
    )

    PM --- P1(
        teamsrv-prod

        Container Apps Environment
        Container App
    )

```