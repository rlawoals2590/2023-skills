import json
import time
import boto3

# athena constant
ATHENA_DATABASE = 'skills_db'
ATHENA_TABLE = 'logs_table'

# S3 constant
ATHENA_OUTPUT_BUCKET = 's3://skills-data-104-abcd/output/'

# number of retries
RETRY_COUNT = 10

def lambda_handler(event, context):
    level = event['queryStringParameters']['level']
    year = event['queryStringParameters']['year']
    month = event['queryStringParameters']['month']
    day = event['queryStringParameters']['day']
    hour = event['queryStringParameters']['hour']

    query = f"SELECT level, host, method, path, code, time FROM {ATHENA_DATABASE}.{ATHENA_TABLE} where level='{level}' and year={year} and month={month} and day={day} and hour={hour}"
    
    athena = boto3.client('athena')

    response = athena.start_query_execution(
        QueryString=query,
        QueryExecutionContext={
            'Database': ATHENA_DATABASE
        },
        ResultConfiguration={
            'OutputLocation': ATHENA_OUTPUT_BUCKET,
        }
    )

    query_execution_id = response['QueryExecutionId']

    # Step 2: Poll for query completion status
    while True:
        response = athena.get_query_execution(QueryExecutionId=query_execution_id)
        status = response['QueryExecution']['Status']['State']

        if status in ['SUCCEEDED', 'FAILED', 'CANCELLED']:
            break

        time.sleep(5)  # Wait for 5 seconds before polling again

    # If query failed, return the reason for failure
    if status == 'FAILED':
        return {
            'statusCode': 400,
            'body': response['QueryExecution']['Status']['StateChangeReason']
        }

    # Step 3: Fetch results
    results = athena.get_query_results(QueryExecutionId=query_execution_id)
    rows = results['ResultSet']['Rows']

    # Only Output data
    # headers = [col['Name'] for col in results['ResultSet']['ResultSetMetadata']['ColumnInfo']]
    # output = [dict(zip(headers, [col['VarCharValue'] for col in row['Data']])) for row in rows[1:]]

    # return {
    #     'statusCode': 200,
    #     'body': output
    # }

    # Return HTML and CSS
    headers = [col['Name'] for col in results['ResultSet']['ResultSetMetadata']['ColumnInfo']]
    html_styles = """
        <style>
            table {
                border-collapse: collapse;
                width: 100%;
                border: 1px solid #ddd;
            }
            th, td {
                border: 1px solid #ddd;
                padding: 8px;
                text-align: left;
            }
            tr:nth-child(even) {
                background-color: #f2f2f2;
            }
            th {
                background-color: #4CAF50;
                color: white;
            }
        </style>
    """
    html_table = "<table>"
    html_table += "<tr>" + "".join([f"<th>{header}</th>" for header in headers]) + "</tr>"
    for row in rows[1:]:
        html_table += "<tr>" + "".join([f"<td>{col['VarCharValue']}</td>" for col in row['Data']]) + "</tr>"
    html_table += "</table>"

    return {
        'statusCode': 200,
        'body': html_styles + html_table,
        'headers': {
            'Content-Type': 'text/html'
        }
    }

