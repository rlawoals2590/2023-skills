def lambda_handler(event, context):
    # CloudFront 요청 객체 가져오기
    request = event['Records'][0]['cf']['request']

    # 쿼리 스트링 가져오기
    query_string = request.get('querystring', '')
    
    # 각 파라미터를 분리하여 딕셔너리 형태로 변환
    parameters = {}
    for param in query_string.split('&'):
        key, value = param.split('=')
        parameters[key] = value

    print(parameters)  # 쿼리 스트링 파라미터 출력

    # 필요한 경우 여기에서 추가 로직 수행...

    # 요청을 그대로 반환하거나 수정하여 반환
    return request
