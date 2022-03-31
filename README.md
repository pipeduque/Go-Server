# Go-Server

Protocol Specification

Commands:

| ID            | Arguments                                    | Description                      |  
| ------------- | -------------------------------------------- | -------------------------------- | 
| CREATE        | [nameChannel]                                | Create a channel                 |
| JOIN          | [nameChannel]                                | Client enters a channel          |
| LEAVE         | [nameChannel]                                | Client leaves a channel          |
| MSG           | [nameChannel];;[messageContent];;[file]      | Send a message                   |
| LIST_CHN      |                                              | List the channels                |
| LIST_MSG      | [nameChannel]                                | List the messages of a channel   |
| LIST_USR      | [nameChannel]                                | List the users of a channel      |
