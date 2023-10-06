CREATE EXTERNAL TABLE IF NOT EXISTS logs_table (
    host string,
    method string,
    path string,
    http string,
    code int,
    mirco float,
    agent string,
    time string,
    epoch_time int
)
PARTITIONED BY (`year` int, `month` int, `day` int, `hour` int, `level` string)
STORED AS PARQUET
LOCATION 's3://skills-data-104-abcd/logs/'
tblproperties ("parquet.compression"="SNAPPY");

MSCK REPAIR TABLE logs_table;

show partitions logs_table;

