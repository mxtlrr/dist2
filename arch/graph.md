```
client 1           server           client 2                            notes
  |-----config------->|<--config-------|
  |                   |                |
  |  compute n digits |   compute k    |<-----|   <<<-|
  |                   |     digits     |      |   <<<----- n or k should be sequential. if n=10,
  |<------------------|--------------->|      |             then computation of k should start
  |                   |                |      |             at the 11th digit.
  | validate client 2 |<-----done------|      |    <<<-|    validation is done by checking first
  |<----------------->|                |      |        |---- two digits of n or k. if discrepnacies
  |--------done------>|validate  client|      |    <<<-|   are found, validation is to done by recomputing
  |                   |        1       |      |                n or k digits.
  |                   |<-------------->|      |
  |               combine              |      |
  |               results              |      |
  |-------------------|----------------|      |
                      |                       |
                      |                       |
                      |-----------------------|
```
