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
    A(Deploy Common Environment Resources)

    B(Deploy Environment Instance Resources </br></br> teamsrv-dev </br> teamsrv-dev-ci)
    
    C[[Create Manual Environment Resources </br></br> B2C]]

    D(Configure Manual Resources)

    A --> B --> C --> D
```

## Deploy application

After the environment is fully deployed, deploy the desired version of the application to the desired environment instance.   The environment ci instance will automatically be deployed on each successful build.

```mermaid

flowchart LR
    A(Deploy teamsrv)

```

## To create multiple environment instances of the application

If there is a use case you can deploy more environment instances.  For example if a developer need to work on their own instance or maybe a new feature needs a dedicated instance.

```mermaid

flowchart LR
    A(Deploy teamsrv)
    B(Deploy Environment Instance Resources </br></br> teamsrv-dev-feature)

    A --> B
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

# Naming

We are using `Organization` because Azure B2C uses organization.

org-global
    Resource Group
        Container Registry: org-global

org-env-manual
    Resource Group
        B2C Tenant: orgauthenv 

Domain 1
    b2c appRegistration
        ui
        api
            scopes
                Groups
                Roles
                Entitlements

    domain1-env-common
        Resource Group
            KeyVault: domain1-env-keyvault
            CosmosDb: domain1-env-cosmosdb

    domain1-env-instance1
        Resource Group
            Container Apps Environment: domain1-env-instance-cae
            Container App: domain1-env-instance-ui
            Container App: domain1-env-instance-api1
            Container App: domain1-env-instance-api2

    domain1-env-instance2
        Resource Group
            Container Apps Environment: domain1-env-instance-cae
            Container App: domain1-env-instance-ui
            Container App: domain1-env-instance-api1
            Container App: domain1-env-instance-api2

Domain 2
    b2c appRegistration
        ui
        api
            scopes
                Groups
                Roles
                Entitlements

    domain2-env-common
        Resource Group
            KeyVault: domain2-env-keyvault
            CosmosDb: domain2-env-cosmosdb

    domain2-env-instance1
        Resource Group
            Container Apps Environment: domain2-env-instance-cae
            Container App: domain2-env-instance-ui
            Container App: domain2-env-instance-api1
            Container App: domain2-env-instance-api2

    domain2-env-instance2
        Resource Group
            Container Apps Environment: domain-env-instance-cae
            Container App: domain2-env-instance-ui
            Container App: domain2-env-instance-api1
            Container App: domain2-env-instance-api2

## Organization Resources

```mermaid

flowchart

    G(majorleaguemiruken-global</br></br> Container Repository) 
        G --- DevC(       majorleaguemiruken-dev-common  </br></br> Key Vault)
     DevC --- DevM(       majorleaguemiruken-dev-manual  </br></br> B2C Tenant)
     DevM --- DevInst1(   majorleaguemiruken-dev         </br></br> Container Apps Environment </br> Container Apps)
     DevM --- DevInst2(   majorleaguemiruken-dev-ci      </br></br> Container Apps Environment </br> Container Apps)

        G --- QaC(        majorleaguemiruken-qa-common   </br></br> Key Vault)
      QaC --- QaM(        majorleaguemiruken-qa-manual   </br></br> B2C Tenant)
      QaM --- QAInst1(    majorleaguemiruken-qa          </br></br> Container Apps Environment </br> Container Apps)

        G --- ProdC(      majorleaguemiruken-prod-common </br></br> Key Vault)
    ProdC --- ProdM(      majorleaguemiruken-prod-manual </br></br> B2C Tenant)
    ProdM --- ProdInst1(  majorleaguemiruken-qa          </br></br> Container Apps Environment </br> Container Apps)
    
```

# Billing Domain Resources
```mermaid

flowchart

    G(majorleaguemiruken-global</br></br> Container Repository) 

    G --- Dom2(Billing Domain)

    Dom2 --- Dom2DevCommon( billing-dev-common  </br></br> CosmosDB)
    Dom2 --- Dom2QACommon(  billing-qa-common   </br></br> CosmosDB)
    Dom2 --- Dom2ProdCommon(billing-prod-common </br></br> CosmosDB)

    Dom2DevCommon --- Dom2DevInst1(billing-dev           </br></br> Stable Dev            </br> Container Apps Environment </br> Container Apps)
    Dom2DevCommon --- Dom2DevInst2(billing-dev-ci        </br></br> CI/CD                 </br> Container Apps Environment </br> Container Apps)
    Dom2DevCommon --- Dom2DevInst3(billing-dev-developerA</br></br> Isolated Feature Work </br> Container Apps Environment </br> Container Apps)

    Dom2QACommon --- Dom2QAInst1(  billing-qa-1</br></br> Container Apps Environment</br> Container Apps)
    Dom2QACommon --- Dom2QAInst2(  billing-qa-2</br></br> Container Apps Environment</br> Container Apps)

    Dom2ProdCommon --- Dom2ProdInst1(billing-prod    </br></br> Production        </br> Container Apps Environment </br> Container Apps)
    Dom2ProdCommon --- Dom2ProdInst2(billing-prod-dr </br></br> Disaster Recovery </br> Container Apps Environment </br> Container Apps)

```

# League Domain Resources
```mermaid

flowchart

    G(majorleaguemiruken-global</br></br> Container Repository) 
    
    G --- Dom3(League Domain)

    Dom3 --- Dom3DevCommon(  league-dev-common  </br></br> CosmosDB)
    Dom3 --- Dom3QACommon(   league-qa-common   </br></br> CosmosDB)
    Dom3 --- Dom3UATCommon(  league-uat-common  </br></br> CosmosDB)
    Dom3 --- Dom3DemoCommon( league-demo-common </br></br> CosmosDB)
    Dom3 --- Dom3ProdCommon( league-prod-common </br></br> CosmosDB)
    Dom3 --- Dom3DRCommon(   league-dr-common   </br></br> CosmosDB)

     Dom3DevCommon --- Dom3DevInst1(  league-dev            </br></br> Stable Dev                </br> Container Apps Environment </br> Container Apps)
     Dom3DevCommon --- Dom3DevInst2(  league-dev-ci         </br></br> CI/CD                     </br> Container Apps Environment </br> Container Apps)
     Dom3DevCommon --- Dom3DevInst3(  league-dev-developerA </br></br> Isolated Feature Work     </br> Container Apps Environment </br> Container Apps)

      Dom3QACommon --- Dom3QAInst1(   league-qa-1           </br></br> Stable QA Env             </br> Container Apps Environment </br> Container Apps)
      Dom3QACommon --- Dom3QAInst2(   league-qa-2           </br></br> New Feature Work          </br> Container Apps Environment </br> Container Apps)
      
     Dom3UATCommon --- Dom3UATInst2(  league-uat            </br></br> User Acceptance           </br> Container Apps Environment </br> Container Apps)

    Dom3DemoCommon --- Dom3DemoInst2( league-demo           </br></br> New Customer Demos        </br> Container Apps Environment </br> Container Apps)

    Dom3ProdCommon --- Dom3ProdInst1( league-prod           </br></br> Production                </br> Container Apps Environment </br> Container Apps)

    Dom3DRCommon   --- Dom3ProdInst2( league-dr             </br></br> Disaster Recovery         </br> Container Apps Environment </br> Container Apps)

```