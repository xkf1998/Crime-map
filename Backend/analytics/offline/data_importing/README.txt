# install maven
  https://maven.apache.org/download.cgi
# run the following commands in this directory
# build 
  mvn package
  export HBASE_CLASSPATH=dataImporting.jar
# run, our table name is group4_rbda_nyu_edu:crimes
  hbase HBaseCrimeImporter ${INPUT_PATH} ${TABLE_NAME}
### test: put data.tsv on HDFS, then run
  hbase HBaseCrimeImporter data.tsv group4:test