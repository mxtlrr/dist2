```
client 1           server           client 2                            notes
  |-----config------->|<--config-------|
  |                   |                |
  |  compute n digits |   compute k    |<-----|   <<<-|
  |                   |     digits     |      |   <<<----- n or k should be sequential. if n=10,
  |<------------------|--------------->|      |             then computation of k should start
  |                   |                |      |             at the 11th digit.
  |--------done------>|<----done-------|      |
  |               combine              |      |
  |               results              |      |
  |-------------------|----------------|      |
                      |                       |
                      |                       |
                      |-----------------------|
                      |
                      | all digits done
                      |
|---------------------|----------------|
|  check n digits     | check k digits |<----|
|<--------------------|--------------->|     |
|                     |                |     |
|----incorrect------->|                |     |
|               modify digits          |     |
|                     |<----correct----|     |
|                     |                |     |
|---------------------|----------------|     |
                      |                      |
                      |----------------------|
                      |
                all validation
                    done
                      |---> exit
```
