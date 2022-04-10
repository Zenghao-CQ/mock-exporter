# mock-exporter

A simple export to generate mock failures

Args: start, period, failureTypes
* export will generate a mock failure if [start, start+period] seconds
* failure code contains [0,failureTypes], default is 3:
    * 0: no failure
    * 1: CPU resource shortage
    * 2: Men resource shortage
    * 3: Disk resource shortage

## References
- https://github.com/SongLee24/prometheus-exporter
