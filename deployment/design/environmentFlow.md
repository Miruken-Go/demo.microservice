## Environment Deployment Flow

```mermaid

flowchart LR
    A(Deploy Global Resources)
    
    B(Deploy Common Environment Resources)

    C[[Configure Manual Environment Resources]]

    D(
        Deploy Environment Instance Resources
    
        Configure Manual Resources
    )

    A --> B --> C --> D
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