import base64
import json
import time
from datetime import datetime

print('Loading function')


def lambda_handler(event, context):
    output = []

    for record in event['records']:
        print(record['recordId'])
        payload = base64.b64decode(record['data']).decode('utf-8')
        data_item = json.loads(payload)
        
        timestamp = data_item['time']
        data_item['time'] = str(datetime.strptime(timestamp, '%Y-%m-%dT%H:%M:%SZ'))
        epoch_timestamp = int(datetime.strptime(timestamp, '%Y-%m-%dT%H:%M:%SZ').timestamp())
        data_item['epoch_time'] = epoch_timestamp
        data_item['code'] = int(data_item['code'])
        data_item['mirco'] = float(data_item['mirco'])
        new_payload = json.dumps(data_item)
        print(new_payload)

        # Do custom processing on the payload here

        output_record = {
            'recordId': record['recordId'],
            'result': 'Ok',
            'data': base64.b64encode(new_payload.encode('utf-8')).decode('utf-8')
        }
        output.append(output_record)

    print('Successfully processed {} records.'.format(len(event['records'])))

    return {'records': output}