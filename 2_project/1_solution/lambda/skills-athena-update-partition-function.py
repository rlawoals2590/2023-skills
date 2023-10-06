import urllib.parse
import boto3

print('Loading function')

s3 = boto3.client('s3')


def lambda_handler(event, context):
    #print("Received event: " + json.dumps(event, indent=2))

    # Get the object from the event and show its content type
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')
    print('Bucket--->' + bucket)
    print('Key------>' + key)
    
    position = key.index('/')
    
    table_name = key[0:int(position)]
    athena_client = boto3.client('athena')
    
    sql = 'MSCK REPAIR TABLE skills_db.' + table_name
    queryContext = {'Database': 'skills_db'}
    resultConfig = {
        'OutputLocation': 's3://' + bucket + '/output/'
    }
    
    try:
        athena_client.start_query_execution(QueryString=sql, QueryExecutionContext=queryContext, ResultConfiguration=resultConfig)
        
        return 'success'
    except Exception as e:
        print(e)
        print('Error getting object {} from bucket {}. Make sure they exist and your bucket is in the same region as this function.'.format(key, bucket))
        raise e
