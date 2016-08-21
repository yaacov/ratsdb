# ratsdb
Rest API for time series data bases

## A Restul API Front end for tsdb's
A Restful API for querying time-value samples. Data is grouped by keys, and labels.
Users can query data using keys and labels. Time-value samples can also grouped in time buckets.

## API Path and varialbes
### Query samples
GET  http://hostname/samples/
#### Query parameters
    key: sample key.
    labels: comma sepreated list of labels.
    start: start time in millisecond since Jan 01 1970.
    end: end time in millisecond since Jan 01 1970.
    bucket: group samples by milliseconds.

### Query one sample
GET  http://hostname/samples/id
### Insert one sample
POST http://hostname/samples/
#### Post body json
```
{ "key": key, "value": value: [, "labels": comma separated list of labels]}
```

## Backends
  Sqlite (memory)

## Examples

### Query all samples
GET: http://localhost:8080/samples
#### Responce

```
[
  {
    "id": 1,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787960527,
    "value": 2
  },
  {
    "id": 2,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787971180,
    "value": 3
  },
  {
    "id": 3,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787976956,
    "value": 4
  }
]
```

### Query one sample by id
GET: http://localhost:8080/samples/2
#### Response
```
{
  "id": 2,
  "key": "Cats",
  "labels": "tabby,ginger",
  "time": 1471787971180,
  "value": 3
}
```

### Query by key and labels
GET: http://localhost:8080/samples/?key=Cats&labels=ginger
```
[
  {
    "id": 1,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787960527,
    "value": 2
  },
  {
    "id": 2,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787971180,
    "value": 3
  },
  {
    "id": 3,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787976956,
    "value": 4
  }
]
```

### Query by time
GET: http://localhost:8080/samples/?key=Cats&start=1471787960527&end=1471787971180
```
[
  {
    "id": 1,
    "key": "Cats",
    "labels": "tabby,ginger",
    "time": 1471787960527,
    "value": 2
  }
]
```

### Query using time buckets
GET: http://localhost:8080/samples/?key=Cats&labels=ginger&bucket=10000
```
[
  {
    "count": 1,
    "key": "Cats",
    "labels": "ginger",
    "start": 1471787960527,
    "end": 1471787960527,
    "min": 2,
    "max": 2,
    "avg": 2
  },
  {
    "count": 2,
    "key": "Cats",
    "labels": "ginger",
    "start": 1471787971180,
    "end": 1471787976956,
    "min": 3,
    "max": 4,
    "avg": 3.5
  }
]
```

### Insert new samples
#### Query
POST: http://localhost:8080/samples

Body:
```
{
  "key": "Cats",
  "labels": "tabby,ginger",
  "value": 2.0
}
```

#### Response
```
{
  "id": 1,
  "key": "Cats",
  "labels": "tabby,ginger",
  "time": 1471787727600,
  "value": 2.0
}
```
