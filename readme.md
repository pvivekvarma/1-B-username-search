# 1 billion username search using PostgreSQL Database

## Setup
```
docker compose up
```

## Notes
### SimpleSearch 
Uses username as the primary key. Searching 100 mil records only takes 1-2 ms. 
#### Advantages:
1. Search is quick
2. insert time duplicate check not required
#### Drawbacks:
1. Changing a user's username is complex. Update all references of the old primary key in all the tables/databases with the new primary key.
2. What if username need not be unique. 
#### Benchmarks:
1. Inserting 100 mil records - 1h5m
2. Searching for a username - 1.2 ms