REST API
========

SESSIONS
--------
####CREATE
Create a new Session.   
```
    PUT http://<host>/sessions
        JSON Payload:  
```

####READ
Get All Sessions in the cluster.  
```
    GET http://<host>/sessions  
```

Get a session information.  
```
    GET http://<host>/sessions/<SessionID> 
```  

Get Sessions based on a filter.  
```    
    GET http://<host>/sessions/filter?state=active
```
   

####UPDATE
Start or Stop an existing Session.
```
    PATCH http://<host>/sessions/<start|stop>/<SessionId>
```
Update an existing Session.
```
    PATCH http://<host>/sessions/update/<SessionId>
          JSON Payload:  
```


####DELETE
Delete an existing Session.
```
    DELETE http://<host>/sessions/<SessionId>
           JSON Payload:  
```



NODES
-------
####CREATE
Add a new node to the cluster  
```
    PUT http://<host>/nodes
        JSON Payload:  
```

####READ
Get all Nodes in the cluster
```   
    GET http://<host>/nodes/
```
Get a Node information

```
    GET http://<host>/nodes/<NodeName>
```
Get a filtered list of Nodes
```
    GET http://<host>/sessions/filter?
```

####UPDATE
Activate or Deactivate an existing Node.
```
    PATCH http://<host>/nodes/<activate|deactivate>/<NodeName>
```
Update an existing Session.
```
    PATCH http://<host>/nodes/update/<NodeName>
          JSON Payload:  
```


####DELETE
Delete an existing Session.
```
    DELETE http://<host>/nodes/<NodeName>
           JSON Payload:  
```

