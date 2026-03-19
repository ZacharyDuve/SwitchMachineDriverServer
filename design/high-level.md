# High Level Design

```mermaid
classDiagram
    Client --> WebsocketBroker
    WebsocketBroker --> WebsocketConnection
    WebsocketConnection <--> Client
```