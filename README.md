# LEIConvert

Converts GLEIF Golden Copy (LEIs) to csv file

Allows for filtering for specific LEIs

based on this article https://blog.singleton.io/posts/2012-06-19-parsing-huge-xml-files-with-go/
and this sample code https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

GLEIF data available here: https://www.gleif.org/en/lei-data/gleif-golden-copy/download-the-golden-copy#/
look for LEI-CDF v2.1 xml files

## Usage

Usage of `./LEIconvert`:
  - `-debug`<br>
    only do a short run for debugging
  - `-filter <filename>`<br>
    file with all LEI to filter for (default "filter.txt")
  - `-in <filename>`<br>
    input file path (default "gleif-goldencopy-lei2-golden-copy.xml")
  - `-missing <filename>`<br>
    file that gets all missing LEIs (default "missing.txt")
  - `-out <filename>`<br>
    output file path (default "found.csv")



### Example Mac OS:
```
./LEIconvert -in 20190830-0800-gleif-goldencopy-lei2-golden-copy.xml -out found.csv -filter filter.txt

LEI Golden Copy XML converter v0.10
Filter LEIs:  8815
Remaining: 0s   Processed records: 1471988 (100.0%)   Found and written records: 8814/8815 (100.0%)                
Not all LEI could be found:
XYZ_DOES_NOX_EXIST
Total processing time:  2m53s
```
