# TimeSeries Assignment

### The code structure is as follows:
- emitter (for emitter as a server, contains its own main.go)
- listener (for listener as a server, contains its own main.go)
- docker-compose, .env, init.sql

### About Table Schema (I have used PostgreSQL)
- The table schema is specified in the init.sql file.
- To optimise performance indexes are used, index(On the minute timestamp) 
  & composite-index (minute timestamp & listenerId) on the table. 

### Frontend 
- The frontend simply displays all the validated decrypted messages.
- Also, the success rates are displayed in the info logs where docker-compose is fired.

### About Docker & localhost
- I am a very beginner with docker, I have tried my best to get this tool working.
- Once the repo is cloned, follow two steps:
    - docker-compose build
    - docker-compose up
- Also, once the above two steps are followed:
    - Go to (http://localhost:53612/home) to see the messages